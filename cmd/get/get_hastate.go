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

func NewGetHaStateCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetHAStateCmd = &cobra.Command{
		Use:     "hastate",
		Short:   "Get a HA state",
		Long:    `It shows HA status in the LoxiLB`,
		Aliases: []string{"ha", "HAstate"},

		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}
			resp, err := client.Status().SetUrl("config/cistate/all").Get(ctx)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				PrintGetHAStateResult(resp, *restOptions)
				return
			}

		},
	}

	return GetHAStateCmd
}

func PrintGetHAStateResult(resp *http.Response, o api.RESTOptions) {
	HAStateresp := api.HAStateGet{}
	var data [][]string
	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}

	if err := json.Unmarshal(resultByte, &HAStateresp); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	// if json options enable, it print as a json format.
	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(HAStateresp, "", "    ")
		fmt.Println(string(resultIndent))
		return
	}

	// Table Init
	table := TableInit()
	table.SetHeader(HASTATE_TITLE)
	// Making load balance data
	for _, HAState := range HAStateresp.HAStateAttr {
		data = append(data, []string{HAState.Instance, HAState.State})
	}
	// Rendering the load balance data to table
	TableShow(data, table)
}