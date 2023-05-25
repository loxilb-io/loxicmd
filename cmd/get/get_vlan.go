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
	"time"

	"github.com/spf13/cobra"
)

func NewGetVlanCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetVlanCmd = &cobra.Command{
		Use:   "vlan",
		Short: "Get a Vlan",
		Long:  `It shows Vlan Information in the loxiLB`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}
			resp, err := client.Vlan().SetUrl("/config/vlan/all").Get(ctx)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				PrintGetVlanResult(resp, *restOptions)
				return
			}

		},
	}

	return GetVlanCmd
}

func PrintGetVlanResult(resp *http.Response, o api.RESTOptions) {
	Vlanresp := api.VlanGet{}
	var data [][]string
	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}

	if err := json.Unmarshal(resultByte, &Vlanresp); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	// if json options enable, it print as a json format.
	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(Vlanresp, "", "    ")
		fmt.Println(string(resultIndent))
		return
	}

	// Sort port Data
	Vlanresp.Sort()

	// Table Init
	table := TableInit()

	// Making load balance data
	for _, vlans := range Vlanresp.Vlans {
		if o.PrintOption == "wide" {
			table.SetHeader(VLAN_WIDE_TITLE)
			data = append(data, []string{vlans.Dev, fmt.Sprintf("%d", vlans.Vid), MemberToString(vlans.Member), VlanStatToString(vlans.Statistic)})
		} else {
			table.SetHeader(VLAN_TITLE)
			data = append(data, []string{vlans.Dev, fmt.Sprintf("%d", vlans.Vid), MemberToString(vlans.Member)})
		}
	}
	// Rendering the load balance data to table
	TableShow(data, table)
}

func MemberToString(members []api.VlanMemberMod) (ret string) {
	for _, member := range members {
		ret += fmt.Sprintf("Device: %s\ntagged: %v\n", member.Dev, member.Tagged)
	}
	return ret
}

func VlanStatToString(stat api.VlanStat) (ret string) {
	ret = fmt.Sprintf("In/Out byte : %d/%d \nIn/Out packets : %d/%d", stat.InBytes, stat.OutBytes, stat.InPackets, stat.OutPackets)
	return ret
}
