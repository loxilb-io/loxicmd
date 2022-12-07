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

func NewGetStatusProcessCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetStatussCmd = &cobra.Command{
		Use:     "process",
		Short:   "Get a process status",
		Long:    `It shows process status in the LoxiLB`,
		Aliases: []string{"Process", "processes"},

		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}
			resp, err := client.Status().SetUrl("status/process").Get(ctx)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				PrintGetStatusProcessResult(resp, *restOptions)
				return
			}

		},
	}

	return GetStatussCmd
}

func PrintGetStatusProcessResult(resp *http.Response, o api.RESTOptions) {
	Processresp := api.ProcessGet{}
	var data [][]string
	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}

	if err := json.Unmarshal(resultByte, &Processresp); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	// if json options enable, it print as a json format.
	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(Processresp, "", "    ")
		fmt.Println(string(resultIndent))
		return
	}

	// Table Init
	table := TableInit()

	// Making load balance data
	for _, Process := range Processresp.ProcessAttr {

		table.SetHeader(PROCESS_TITLE)
		data = append(data, []string{Process.Pid, Process.User, Process.Priority, Process.Nice, Process.VirtMemory,
			Process.ResidentSize, Process.SharedMemory, Process.Status,
			Process.CPUUsage, Process.MemoryUsage, Process.Command})

	}
	// Rendering the load balance data to table
	TableShow(data, table)
}

func NewGetStatusDeviceCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetStatussCmd = &cobra.Command{
		Use:     "device",
		Short:   "Get a device status",
		Long:    `It shows device status in the LoxiLB`,
		Aliases: []string{"devices"},

		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}
			resp, err := client.Status().SetUrl("status/device").Get(ctx)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				PrintGetStatusDeviceResult(resp, *restOptions)
				return
			}

		},
	}

	return GetStatussCmd
}

func PrintGetStatusDeviceResult(resp *http.Response, o api.RESTOptions) {
	Deviceresp := api.DeviceGet{}
	var data [][]string
	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}

	if err := json.Unmarshal(resultByte, &Deviceresp); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	// if json options enable, it print as a json format.
	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(Deviceresp, "", "    ")
		fmt.Println(string(resultIndent))
		return
	}

	// Table Init
	table := TableInit()

	// Making load balance data

	table.SetHeader(PROCESS_TITLE)
	data = append(data, []string{Deviceresp.HostName, Deviceresp.MachineID, Deviceresp.BootID, Deviceresp.OS, Deviceresp.Kernel, Deviceresp.Architecture, Deviceresp.Uptime})

	// Rendering the load balance data to table
	TableShow(data, table)
}

func NewGetStatusFileSystemCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetStatussCmd = &cobra.Command{
		Use:     "filesystem",
		Short:   "Get a filesystem status",
		Long:    `It shows filesystem status in the LoxiLB`,
		Aliases: []string{"fs"},

		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}
			resp, err := client.Status().SetUrl("status/filesystem").Get(ctx)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				PrintGetStatusFilesystemResult(resp, *restOptions)
				return
			}

		},
	}

	return GetStatussCmd
}

func PrintGetStatusFilesystemResult(resp *http.Response, o api.RESTOptions) {
	Processresp := api.ProcessGet{}
	var data [][]string
	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}

	if err := json.Unmarshal(resultByte, &Processresp); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	// if json options enable, it print as a json format.
	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(Processresp, "", "    ")
		fmt.Println(string(resultIndent))
		return
	}

	// Table Init
	table := TableInit()

	// Making load balance data
	for _, Process := range Processresp.ProcessAttr {

		table.SetHeader(PROCESS_TITLE)
		data = append(data, []string{Process.Pid, Process.User, Process.Priority, Process.Nice, Process.VirtMemory,
			Process.ResidentSize, Process.SharedMemory, Process.Status,
			Process.CPUUsage, Process.MemoryUsage, Process.Command})

	}
	// Rendering the load balance data to table
	TableShow(data, table)
}
