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
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"loxicmd/pkg/api"

	"github.com/spf13/cobra"
)

func DeleteVlanMemberValidation(args []string) error {
	if len(args) > 3 {
		fmt.Println("delete VlanMember command get so many args")
	}
	_, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	return nil
}

type DeleteVlanMemberOptions struct {
	Tagged bool
}

func NewDeleteVlanMemberCmd(restOptions *api.RESTOptions) *cobra.Command {
	o := DeleteVlanMemberOptions{}
	var deleteVlanMemberCmd = &cobra.Command{
		Use:     "vlanmember <Vid> <DeviceName> --tagged=<Tagged>",
		Short:   "Delete a VlanMember",
		Long:    `Delete a VlanMember using Vid in the LoxiLB.`,
		Aliases: []string{"vlanMember", "vlan-member", "vlan_member"},

		Run: func(cmd *cobra.Command, args []string) {
			if err := DeleteVlanMemberValidation(args); err != nil {
				fmt.Println("not valid <Vid>")
				return
			}
			Vid := args[0]
			Dev := args[1]
			Tagged := fmt.Sprintf("%v", o.Tagged)

			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(ctx, time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}
			subResources := []string{
				Vid, "member", Dev, "tagged", Tagged,
			}
			resp, err := client.Vlan().SubResources(subResources).Delete(ctx)
			if err != nil {
				fmt.Printf("Error: Failed to delete VlanMember : %s", Vid)
				return
			}
			defer resp.Body.Close()
			fmt.Printf("Debug: response.StatusCode: %d\n", resp.StatusCode)
			if resp.StatusCode != http.StatusOK {
				PrintDeleteLbResult(resp, *restOptions)
				return
			}

		},
	}
	deleteVlanMemberCmd.Flags().BoolVarP(&o.Tagged, "tagged", "", false, "Tagged mode Vlan")

	return deleteVlanMemberCmd
}
