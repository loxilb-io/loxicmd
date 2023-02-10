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
	"os"
	"strconv"
	"time"

	"loxicmd/pkg/api"

	"github.com/spf13/cobra"
)

func DeleteVxlanPeerValidation(args []string) error {
	if len(args) > 3 {
		fmt.Println("delete vxlanPeer command get so many args")
	}
	_, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	return nil
}

func NewDeleteVxlanPeerCmd(restOptions *api.RESTOptions) *cobra.Command {
	var deleteVxlanPeerCmd = &cobra.Command{
		Use:   "vxlanpeer <Vnid> <PeerIP> ",
		Short: "Delete a vxlanPeer",
		Long: `Delete a vxlanPeer using Vnid in the LoxiLB.
ex) loxicmd delete vxlan-peer 100 30.1.3.1
		
		`,
		Aliases: []string{"vxlanPeer", "vxlan-peer", "vxlan_peer"},
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := DeleteVxlanPeerValidation(args); err != nil {
				fmt.Println("not valid <Vnid>")
				return
			}
			Vnid := args[0]
			PeerIP := args[1]
			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(ctx, time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}
			subResources := []string{
				Vnid, "peer", PeerIP,
			}
			resp, err := client.Vxlan().SubResources(subResources).Delete(ctx)
			if err != nil {
				fmt.Printf("Error: Failed to delete vxlanPeer : %s", Vnid)
				return
			}
			defer resp.Body.Close()
			fmt.Printf("Debug: response.StatusCode: %d\n", resp.StatusCode)
			if resp.StatusCode == http.StatusOK {
				PrintDeleteResult(resp, *restOptions)
				return
			}

		},
	}

	return deleteVxlanPeerCmd
}
