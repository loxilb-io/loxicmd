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

func NewGetBFDCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetBFDCmd = &cobra.Command{
		Use:     "bfd",
		Short:   "Get all BFD sessions",
		Long:    `It shows BFD Sessions in the LoxiLB`,

		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}
			resp, err := client.Status().SetUrl("config/bfd/all").Get(ctx)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				PrintGetBFDResult(resp, *restOptions)
				return
			}

		},
	}

	return GetBFDCmd
}

func PrintGetBFDResult(resp *http.Response, o api.RESTOptions) {
	BFDresp := api.BFDSessionGet{}
	var data [][]string
	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}

	if err := json.Unmarshal(resultByte, &BFDresp); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	// if json options enable, it print as a json format.
	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(BFDresp, "", "    ")
		fmt.Println(string(resultIndent))
		return
	}

	// Table Init
	table := TableInit()
	
	// Making load balance data
	for _, bfd := range BFDresp.BFDSessionAttr {
		if (o.PrintOption == "wide") {
			table.SetHeader(BFD_WIDE_TITLE)
			data = append(data, []string{bfd.Instance, bfd.RemoteIP, bfd.SourceIP,
			 fmt.Sprintf("%d",bfd.Port), fmt.Sprintf("%d us",bfd.Interval), fmt.Sprintf("%d",bfd.RetryCount), bfd.State})
		} else {
			table.SetHeader(BFD_TITLE)
			data = append(data, []string{bfd.Instance, bfd.RemoteIP, bfd.State})
		}
	}
	// Rendering the load balance data to table
	TableShow(data, table)
}