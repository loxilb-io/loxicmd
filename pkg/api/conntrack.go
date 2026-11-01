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
package api

import (
	"fmt"
	"sort"
)

type Conntrack struct {
	CommonAPI
}

type CtInformationGet struct {
	CtInfo []ConntrackInformation `json:"ctAttr"`
}

type ConntrackInformation struct {
	Dip      string `json:"destinationIP"`
	Sip      string `json:"sourceIP"`
	Dport    uint16 `json:"destinationPort"`
	Sport    uint16 `json:"sourcePort"`
	Proto    string `json:"protocol"`
	Ident    string `json:"ident"`
	CState   string `json:"conntrackState"`
	CAct     string `json:"conntrackAct"`
	Pkts     uint64 `json:"packets"`
	Bytes    uint64 `json:"bytes"`
	ServName string `json:"servName"`
}

func (ct ConntrackInformation) Key() string {
	return fmt.Sprintf("%s|%s|%05d|%05d|%s", ct.Dip, ct.Sip, ct.Dport, ct.Sport, ct.Proto)
}

// Sort orders by Key. Each key is made once: at 200 K entries, making it per
// comparison was seconds of Sprintf. sort.Slice stays, and the rows move with
// their keys, so equal keys keep the order they always had.
func (ctresp CtInformationGet) Sort() {
	type keyedRow struct {
		key string
		idx int
	}
	rows := make([]keyedRow, len(ctresp.CtInfo))
	for i := range ctresp.CtInfo {
		rows[i] = keyedRow{ctresp.CtInfo[i].Key(), i}
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].key < rows[j].key
	})
	sorted := make([]ConntrackInformation, len(rows))
	for i := range rows {
		sorted[i] = ctresp.CtInfo[rows[i].idx]
	}
	copy(ctresp.CtInfo, sorted)
}
