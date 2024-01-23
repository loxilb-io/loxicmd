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
		Short: "Get a Load balance features from LoxiLB.",
		Long: `Get a Load balance features from LoxiLB.
	Get - Service type external load-balancer, Vlan, Vxlan, Qos Policies, 
	 Endpoint client,FDB, IPaddress, Neighbor, Route,Firewall, Mirror, Session, UlCl
	Get Port(interface) dump used by loxilb or its docker
	Get Connection track (TCP/UDP/ICMP/SCTP) information	
`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd
			_ = args
			if len(args) == 0 {
				cmd.Help()
			}
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			fmt.Printf("Error: unknown command \"%v\"for \"loxicmd\" \nRun \"loxicmd --help\" for usage.\n", args)
			cmd.Help()
			return err
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
	GetCmd.AddCommand(NewGetMirrorCmd(restOptions))
	GetCmd.AddCommand(NewGetFirewallCmd(restOptions))
	GetCmd.AddCommand(NewGetFDBCmd(restOptions))
	GetCmd.AddCommand(NewGetVxlanCmd(restOptions))
	GetCmd.AddCommand(NewGetEndPointCmd(restOptions))
	GetCmd.AddCommand(NewGetLogLevelCmd(restOptions))
	GetCmd.AddCommand(NewGetBGPNeighborCmd(restOptions))

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
