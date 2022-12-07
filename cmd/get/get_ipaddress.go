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
	"fmt"
	"io"
	"loxicmd/pkg/api"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func NewGetIPAddressCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetIPAddressCmd = &cobra.Command{
		Use:     "ip",
		Short:   "Get a IP Address",
		Long:    `It shows IP Address Information in the LoxiLB`,
		Aliases: []string{"ipv4address", "ipv4"},

		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}
			resp, err := client.IPv4Address().SetUrl("/config/ipv4address/all").Get(ctx)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				PrintGetIPAddressResult(resp, *restOptions)
				return
			}

		},
	}

	return GetIPAddressCmd
}

func PrintGetIPAddressResult(resp *http.Response, o api.RESTOptions) {
	IPv4Addressresp := api.Ipv4AddrModGet{}
	var data [][]string
	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}

	if err := json.Unmarshal(resultByte, &IPv4Addressresp); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	// if json options enable, it print as a json format.
	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(IPv4Addressresp, "", "    ")
		fmt.Println(string(resultIndent))
		return
	}

	// Table Init
	table := TableInit()

	// Making load balance data
	for _, IPv4Addressrule := range IPv4Addressresp.IPv4Attr {
		if o.PrintOption == "wide" {
			table.SetHeader(IP_WIDE_TITLE)
			data = append(data, []string{IPv4Addressrule.Dev, MakeIPv4String(IPv4Addressrule.IP), fmt.Sprintf("%d", IPv4Addressrule.Sync)})

		} else {
			table.SetHeader(IP_TITLE)
			data = append(data, []string{IPv4Addressrule.Dev, MakeIPv4String(IPv4Addressrule.IP)})
		}
	}
	// Rendering the load balance data to table
	TableShow(data, table)
}

func MakeIPv4String(ips []string) (ret string) {
	for _, ip := range ips {
		ret += fmt.Sprintf("%s\n", ip)
	}
	ret = strings.TrimSpace(ret)
	return ret
}
