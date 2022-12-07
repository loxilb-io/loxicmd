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
	"fmt"
	"os"

	"loxicmd/pkg/api"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func GetCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetCmd = &cobra.Command{
		Use:   "get",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly Get a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd
			_ = args
			fmt.Println("Get called")
		},
	}

	GetCmd.AddCommand(NewGetLoadBalancerCmd(restOptions))
	GetCmd.AddCommand(NewGetConntrackCmd(restOptions))
	GetCmd.AddCommand(NewGetPortCmd(restOptions))
	GetCmd.AddCommand(NewGetSessionCmd(restOptions))
	GetCmd.AddCommand(NewGetSessionULCLCmd(restOptions))
	GetCmd.AddCommand(NewGetPolicyCmd(restOptions))
	GetCmd.AddCommand(NewGetRouteCmd(restOptions))
	GetCmd.AddCommand(NewGetIPAddressCmd(restOptions))
	GetCmd.AddCommand(NewGetNeighborCmd(restOptions))
	GetCmd.AddCommand(NewGetStatusProcessCmd(restOptions))
	GetCmd.AddCommand(NewGetVlanCmd(restOptions))

	return GetCmd
}

func TableInit() *tablewriter.Table {
	// Table Init
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	return table
}

func TableShow(data [][]string, table *tablewriter.Table) {
	table.AppendBulk(data)
	table.Render()
}
