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

func NewGetVxlanCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetvxlanCmd = &cobra.Command{
		Use:   "vxlan",
		Short: "Get a vxlan",
		Long:  `It shows vxlan Information in the loxiLB`,
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}
			resp, err := client.Vxlan().SetUrl("/config/tunnel/vxlan/all").Get(ctx)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				PrintGetvxlanResult(resp, *restOptions)
				return
			}

		},
	}

	return GetvxlanCmd
}

func PrintGetvxlanResult(resp *http.Response, o api.RESTOptions) {
	vxlanresp := api.VxlanGet{}
	var data [][]string
	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}

	if err := json.Unmarshal(resultByte, &vxlanresp); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	// if json options enable, it print as a json format.
	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(vxlanresp, "", "    ")
		fmt.Println(string(resultIndent))
		return
	}
	// Sort vxlan Data
	vxlanresp.Sort()

	// Table Init
	table := TableInit()

	// Making load balance data
	for _, vxlans := range vxlanresp.VxlanAttr {

		table.SetHeader(VXLAN_TITLE)
		data = append(data, []string{vxlans.VxlanName, fmt.Sprintf("%d", vxlans.VxLanID), vxlans.EndpointDev, MakePeerToSting(vxlans.PeerIP)})

	}
	// Rendering the load balance data to table
	TableShow(data, table)
}
func MakePeerToSting(peerIPs []string) (ret string) {
	for _, peerIP := range peerIPs {
		ret += peerIP + "\n"
	}
	ret = strings.TrimSpace(ret)
	return ret
}
