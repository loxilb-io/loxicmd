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

func NewGetConntrackCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetctCmd = &cobra.Command{
		Use:     "conntrack",
		Aliases: []string{"ct", "conntracks", "cts"},
		Short:   "Get a Conntrack",
		Long:    `It shows connection track Information`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd
			_ = args
			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}
			resp, err := client.Conntrack().Get(ctx)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				PrintGetCTResult(resp, *restOptions)
				return
			}

		},
	}
	GetctCmd.Flags().StringVarP(&restOptions.ServiceName, "servName", "", restOptions.ServiceName, "Name for load balancer rule")
	return GetctCmd
}

func PrintGetCTResult(resp *http.Response, o api.RESTOptions) {
	ctresp := api.CtInformationGet{}
	var data [][]string
	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}

	if err := json.Unmarshal(resultByte, &ctresp); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	// if json options enable, it print as a json format.
	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(ctresp, "", "    ")
		fmt.Println(string(resultIndent))
		return
	}

	ctresp.Sort()

	// Table Init
	table := TableInit()
	table.SetHeader([]string{"Service Name","destIP", "srcIP", "dport", "sport", "proto", "state", "act", "packets", "bytes"})
	// Making load balance data
	data = makeConntrackData(o, ctresp)

	// Rendering the load balance data to table
	TableShow(data, table)
}

func makeConntrackData(o api.RESTOptions,ctresp api.CtInformationGet) (data [][]string) {
	for _, conntrack := range ctresp.CtInfo {
		if o.ServiceName != "" && o.ServiceName != conntrack.ServName {
			continue
		}
		data = append(data, []string{
			conntrack.ServName,
			conntrack.Dip,
			conntrack.Sip,
			fmt.Sprintf("%d", conntrack.Dport),
			fmt.Sprintf("%d", conntrack.Sport),
			conntrack.Proto,
			conntrack.CState,
			conntrack.CAct,
			fmt.Sprintf("%v", conntrack.Pkts),
			fmt.Sprintf("%v", conntrack.Bytes),
		})
	}
	return data
}
