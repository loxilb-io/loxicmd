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
package delete

import (
	"fmt"

	"loxicmd/pkg/api"

	"github.com/spf13/cobra"
)

func DeleteCmd(restOptions *api.RESTOptions) *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a Load balance features in the LoxiLB.",
		Long:  `Delete a Load balance features in the LoxiLB. `,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("delete called")
		},
	}

	deleteCmd.AddCommand(NewDeleteLoadBalancerCmd(restOptions))
	deleteCmd.AddCommand(NewDeleteSessionCmd(restOptions))
	deleteCmd.AddCommand(NewDeleteSessionUlClCmd(restOptions))
	deleteCmd.AddCommand(NewDeletePolicyCmd(restOptions))
	deleteCmd.AddCommand(NewDeleteRouteCmd(restOptions))
	deleteCmd.AddCommand(NewDeleteIPv4AddressCmd(restOptions))
	deleteCmd.AddCommand(NewDeleteNeighborsCmd(restOptions))
	deleteCmd.AddCommand(NewDeleteFDBCmd(restOptions))
	deleteCmd.AddCommand(NewDeleteVlanBridgeCmd(restOptions))
	deleteCmd.AddCommand(NewDeleteVlanMemberCmd(restOptions))
	deleteCmd.AddCommand(NewDeleteVxlanBridgeCmd(restOptions))
	deleteCmd.AddCommand(NewDeleteVxlanPeerCmd(restOptions))

	return deleteCmd
}
