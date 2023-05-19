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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"loxicmd/pkg/api"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func NewGetFirewallCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetfwCmd = &cobra.Command{
		Use:     "firewall",
		Short:   "Get a firewall",
		Aliases: []string{"Firewall", "fw", "firewalls"},
		Long:    `It shows Load balancer Information`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd
			_ = args
			resp, err := FWAPICall(restOptions)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				PrintGetFWResult(resp, *restOptions)
				return
			}

		},
	}

	return GetfwCmd
}

func PrintGetFWResult(resp *http.Response, o api.RESTOptions) {
	fwresp := api.FWInformationGet{}
	var data [][]string
	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}

	if err := json.Unmarshal(resultByte, &fwresp); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	// if json options enable, it print as a json format.
	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(fwresp, "", "    ")
		fmt.Println(string(resultIndent))
		return
	}

	fwresp.Sort()

	// Table Init
	table := TableInit()

	// Making load balance data
	for _, fwrule := range fwresp.FWInfo {
		table.SetHeader(FIREWALL_TITLE)
		data = append(data, []string{fwrule.Rule.SrcIP, fwrule.Rule.DstIP, fmt.Sprintf("%d", fwrule.Rule.SrcPortMin), fmt.Sprintf("%d", fwrule.Rule.SrcPortMax),
			fmt.Sprintf("%d", fwrule.Rule.DstPortMin), fmt.Sprintf("%d", fwrule.Rule.DstPortMax), fmt.Sprintf("%d", fwrule.Rule.Proto),
			fwrule.Rule.InPort, fmt.Sprintf("%d", fwrule.Rule.Pref), MakeFirewallOptionToString(fwrule.Opts)})

	}

	// Rendering the load balance data to table
	TableShow(data, table)
}

func MakeFirewallOptionToString(t api.FwOptArg) (ret string) {
	if t.Allow {
		ret = "Allow"
	} else if t.Drop {
		ret = "Drop"
	} else if t.Trap {
		ret = "Trap"
	} else if t.Rdr {
		ret = fmt.Sprintf("Redirect(%s)", t.RdrPort)
	}
	if t.Record {
		ret += fmt.Sprintf(",Record")
	}
	if t.Mark != 0 {
		ret += fmt.Sprintf(",FwMark(%d)", t.Mark)
	}
	return ret
}

func FWAPICall(restOptions *api.RESTOptions) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}
	resp, err := client.Firewall().SetUrl("/config/firewall/all").Get(ctx)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return nil, err
	}
	return resp, nil
}

func FWdump(restOptions *api.RESTOptions, path string) (string, error) {
	// File Open
	fileP := []string{"FWconfig_", ".txt"}
	t := time.Now()
	file := strings.Join(fileP, t.Local().Format("2006-01-02_15:04:05"))
	f, err := os.Create(file)
	if err != nil {
		fmt.Printf("Can't create dump file\n")
		os.Exit(1)
	}
	defer f.Close()

	// API Call
	resp, err := FWAPICall(restOptions)
	if err != nil {
		return "", err
	}
	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
	}
	// Write
	f.Write(resultByte)

	cfile := path + "FWconfig.txt"
	if _, err := os.Stat(cfile); errors.Is(err, os.ErrNotExist) {
		if err != nil {
			fmt.Println("There is no saved config file")
		}
	} else {
		command := "mv " + cfile + " " + cfile + ".bk"
		cmd := exec.Command("bash", "-c", command)
		_, err := cmd.Output()
		if err != nil {
			fmt.Println("Can't backup ", cfile)
			return file, err
		}
	}
	command := "cp -R " + file + " " + cfile
	cmd := exec.Command("bash", "-c", command)
	fmt.Println(cmd)
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	return file, nil
}
