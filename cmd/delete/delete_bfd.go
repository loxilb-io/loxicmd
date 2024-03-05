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
	"loxicmd/pkg/api"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func NewDeleteBFDCmd(restOptions *api.RESTOptions) *cobra.Command {
	o := api.BFDSessionInfo{}

	var deleteBFDCmd = &cobra.Command{
		Use:   "bfd remoteIP [--instance=<instance>]",
		Short: "Delete a BFD session",
		Long: `Delete a BFD session for HA failover.

ex) loxicmd delete bfd 32.32.32.2 --instance=default"
		`,
		Aliases: []string{"bfd-session"},
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(ctx, time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}

			if val := net.ParseIP(args[0]); val != nil {
				o.RemoteIP = args[0]
			} else {
				fmt.Printf("remoteIP '%s' is invalid format\n", args[0])
				return
			}
			subResources := []string{
				"remoteIP", o.RemoteIP,
			}

			qmap := map[string]string{}
			qmap["instance"] = o.Instance

			resp, err := client.BFDSession().SubResources(subResources).Query(qmap).Delete(ctx)
			if err != nil {
				fmt.Printf("Error: Failed to delete Firewall")
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
	deleteBFDCmd.Flags().StringVarP(&o.Instance, "instance", "", "default", "Specify the cluster instance name")
	
	return deleteBFDCmd
}
