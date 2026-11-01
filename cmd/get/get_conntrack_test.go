/*
 * Copyright (c) 2022 NetLOX Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package get

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"loxicmd/pkg/api"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

// The golden files hold what this command printed before its table output was
// buffered and its sort changed. Regenerate them only from a build whose output
// is known to be right.
var updateGolden = flag.Bool("update", false, "rewrite testdata/ctget from the current code")

func ctRow(dip, sip string, dport, sport uint16, proto string, pkts uint64) api.ConntrackInformation {
	return api.ConntrackInformation{
		Dip: dip, Sip: sip, Dport: dport, Sport: sport, Proto: proto, Ident: "0",
		CState: "est", CAct: "fdnat-10.0.0.1:80:w1", Pkts: pkts, Bytes: pkts * 100, ServName: "",
	}
}

func ctCases() map[string][]api.ConntrackInformation {
	named := ctRow("20.20.20.1", "10.0.0.1", 80, 1000, "tcp", 1)
	named.ServName = "web"
	spaced := ctRow("20.20.20.1", "10.0.0.2", 80, 1000, "tcp", 2)
	spaced.ServName = "a service name long enough for the table to wrap it"
	spaced.CAct = "fdnat-10.0.0.1:80:w1 and another action past the width of a column"
	icmp := ctRow("20.20.20.1", "10.0.0.1", 0, 0, "icmp", 3)
	icmp.Ident = "4321"
	icmp.CState = "bi-est"

	return map[string][]api.ConntrackInformation{
		"empty": {},
		"one":   {ctRow("20.20.20.1", "10.0.0.1", 80, 1000, "tcp", 1)},
		"two": {
			ctRow("20.20.20.2", "10.0.0.1", 80, 1000, "tcp", 1),
			ctRow("20.20.20.1", "10.0.0.1", 80, 1000, "tcp", 2),
		},
		// The same key, told apart by the packet count: their order must not change
		"duplicates": {
			ctRow("20.20.20.1", "10.0.0.1", 80, 1000, "tcp", 3),
			ctRow("20.20.20.1", "10.0.0.1", 80, 1000, "tcp", 1),
			ctRow("20.20.20.1", "10.0.0.1", 80, 1000, "tcp", 2),
			ctRow("20.20.20.0", "10.0.0.1", 80, 1000, "tcp", 4),
			ctRow("20.20.20.1", "10.0.0.1", 80, 1000, "tcp", 5),
		},
		// String order, not numeric: 10.0.0.10 sorts before 10.0.0.9
		"mixedfamily": {
			ctRow("10.0.0.9", "10.0.0.1", 80, 1000, "tcp", 1),
			ctRow("10.0.0.10", "10.0.0.1", 80, 1000, "tcp", 2),
			ctRow("2001:db8:1::102", "2001:db8:1::10", 80, 1000, "tcp", 3),
			ctRow("2001:db8:0:0:0:0:0:1", "2001:db8:ffff:ffff:ffff:ffff:ffff:ffff", 65535, 1, "udp", 4),
			ctRow("20.20.20.1", "2001:db8:1::10", 9, 65535, "tcp", 5),
		},
		"protocols": {
			ctRow("20.20.20.1", "10.0.0.1", 80, 1000, "tcp", 1),
			ctRow("20.20.20.1", "10.0.0.1", 80, 1000, "udp", 2),
			ctRow("20.20.20.1", "10.0.0.1", 80, 1000, "sctp", 3),
			icmp,
			ctRow("2001:db8:1::102", "2001:db8:1::10", 0, 0, "icmp6", 4),
			ctRow("20.20.20.1", "10.0.0.1", 80, 1000, "unknown", 5),
		},
		"cellwrap": {named, spaced},
		"bigcounters": {
			ctRow("20.20.20.1", "10.0.0.1", 80, 1000, "tcp", 18446744073709551615/100),
		},
	}
}

// Past the sort's insertion-sort size, with many equal keys: the order of equal
// rows then depends on every comparison the sort makes
func ctLarge(n int) []api.ConntrackInformation {
	rows := make([]api.ConntrackInformation, n)
	x := uint32(2463534242)
	for i := range rows {
		x ^= x << 13
		x ^= x >> 17
		x ^= x << 5
		rows[i] = ctRow(fmt.Sprintf("20.20.%d.%d", x%3, x%7), fmt.Sprintf("10.0.%d.%d", (x>>8)%2, (x>>16)%9),
			80, uint16(1000+(x>>4)%5), []string{"tcp", "udp", "sctp"}[(x>>12)%3], uint64(i))
	}
	return rows
}

// captureStdout runs f with os.Stdout on a pipe, as the command's output goes to a pipe
func captureStdout(t *testing.T, f func()) []byte {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = stdout }()

	got := make(chan []byte)
	go func() {
		b, _ := io.ReadAll(r)
		got <- b
	}()
	f()
	w.Close()
	return <-got
}

func runCtGet(t *testing.T, body []byte, o api.RESTOptions) []byte {
	t.Helper()
	return captureStdout(t, func() {
		PrintGetCTResult(&http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(body))}, o)
	})
}

func checkGolden(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join("testdata", "ctget", name)
	if *updateGolden {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("%s differs from the golden output\n--- got\n%s\n--- want\n%s", name, got, want)
	}
}

func TestPrintGetCTResultGolden(t *testing.T) {
	for name, rows := range ctCases() {
		body, err := json.Marshal(api.CtInformationGet{CtInfo: rows})
		if err != nil {
			t.Fatal(err)
		}
		for label, format := range map[string]string{"table": "", "wide": "wide", "json": "json"} {
			t.Run(name+"-"+label, func(t *testing.T) {
				checkGolden(t, name+"."+label+".golden", runCtGet(t, body, api.RESTOptions{PrintOption: format}))
			})
		}
	}
}

func TestPrintGetCTResultServiceFilter(t *testing.T) {
	body, err := json.Marshal(api.CtInformationGet{CtInfo: ctCases()["cellwrap"]})
	if err != nil {
		t.Fatal(err)
	}
	checkGolden(t, "filter.match.golden", runCtGet(t, body, api.RESTOptions{ServiceName: "web"}))
	checkGolden(t, "filter.none.golden", runCtGet(t, body, api.RESTOptions{ServiceName: "absent"}))
}

// Errors print the same, and print once
func TestPrintGetCTResultBadBody(t *testing.T) {
	checkGolden(t, "badbody.golden", runCtGet(t, []byte("{not json"), api.RESTOptions{}))
}

// The large outputs are kept as a hash: 5000 rows of table are not worth a file
func TestPrintGetCTResultLarge(t *testing.T) {
	for _, n := range []int{13, 50, 5000} {
		body, err := json.Marshal(api.CtInformationGet{CtInfo: ctLarge(n)})
		if err != nil {
			t.Fatal(err)
		}
		out := runCtGet(t, body, api.RESTOptions{})
		checkGolden(t, fmt.Sprintf("large%d.sha256", n), []byte(fmt.Sprintf("%x\n", sha256.Sum256(out))))
	}
}

// LOXICMD_CT_RESPONSE names a captured REST response; the 200 K one is too large
// for the repository. Its output must hash to LOXICMD_CT_RESPONSE_SHA256.
func TestPrintGetCTResultCaptured(t *testing.T) {
	path := os.Getenv("LOXICMD_CT_RESPONSE")
	if path == "" {
		t.Skip("LOXICMD_CT_RESPONSE not set")
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got := fmt.Sprintf("%x", sha256.Sum256(runCtGet(t, body, api.RESTOptions{})))
	if want := os.Getenv("LOXICMD_CT_RESPONSE_SHA256"); want != "" && got != want {
		t.Errorf("table output hash %s, want %s", got, want)
	}
	t.Logf("table output sha256 %s", got)
}
