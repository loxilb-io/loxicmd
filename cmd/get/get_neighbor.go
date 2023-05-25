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

func NewGetNeighborCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetNeighborCmd = &cobra.Command{
		Use:     "neighbor",
		Short:   "Get a neighbors",
		Long:    `It shows neighbors Information in the LoxiLB`,
		Aliases: []string{"nei", "neigh"},

		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}
			resp, err := client.Neighbor().SetUrl("/config/neighbor/all").Get(ctx)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				PrintGetNeighborResult(resp, *restOptions)
				return
			}

		},
	}

	return GetNeighborCmd
}

func PrintGetNeighborResult(resp *http.Response, o api.RESTOptions) {
	Neighborsresp := api.NeighborModGet{}
	var data [][]string
	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}

	if err := json.Unmarshal(resultByte, &Neighborsresp); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	// if json options enable, it print as a json format.
	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(Neighborsresp, "", "    ")
		fmt.Println(string(resultIndent))
		return
	}

	Neighborsresp.Sort()

	// Table Init
	table := TableInit()

	// Making load balance data
	for _, neighbor := range Neighborsresp.NeighborAttr {
		table.SetHeader(NEIGHBOR_TITLE)
		data = append(data, []string{neighbor.IP, neighbor.Dev, neighbor.MacAddress})

	}
	// Rendering the load balance data to table
	TableShow(data, table)
}
