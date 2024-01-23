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
package create

import (
	"fmt"
	"loxicmd/pkg/api"

	"github.com/spf13/cobra"
)

func CreateCmd(restOptions *api.RESTOptions) *cobra.Command {
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a Load balance features in the LoxiLB.",
		Long: `Create a Load balance features in the LoxiLB.
Create - Service type external load-balancer, Vlan, Vxlan, Qos Policies, 
	 Endpoint client,FDB, IPaddress, Neighbor, Route,Firewall, Mirror, Session, UlCl
`,
		Run: func(cmd *cobra.Command, args []string) {
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

	createCmd.AddCommand(NewCreateLoadBalancerCmd(restOptions))
	createCmd.AddCommand(NewCreateSessionCmd(restOptions))
	createCmd.AddCommand(NewCreateSessionUlClCmd(restOptions))
	createCmd.AddCommand(NewCreatePolicyCmd(restOptions))
	createCmd.AddCommand(NewCreateRouteCmd(restOptions))
	createCmd.AddCommand(NewCreateIPv4AddressCmd(restOptions))
	createCmd.AddCommand(NewCreateNeighborsCmd(restOptions))
	createCmd.AddCommand(NewCreateFDBCmd(restOptions))
	createCmd.AddCommand(NewCreateVlanBridgeCmd(restOptions))
	createCmd.AddCommand(NewCreateVlanMemberCmd(restOptions))
	createCmd.AddCommand(NewCreateVxlanBridgeCmd(restOptions))
	createCmd.AddCommand(NewCreateVxlanPeerCmd(restOptions))
	createCmd.AddCommand(NewCreateMirrorCmd(restOptions))
	createCmd.AddCommand(NewCreateFirewallCmd(restOptions))
	createCmd.AddCommand(NewCreateEndPointCmd(restOptions))
	createCmd.AddCommand(NewCreateBGPNeighborCmd(restOptions))

	return createCmd
}
