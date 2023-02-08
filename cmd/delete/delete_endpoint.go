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
	"net"
	"net/http"
	"time"

	"loxicmd/pkg/api"

	"github.com/spf13/cobra"
)

type DeleteEndPoint struct {
	Host string
}

func NewDeleteEndPointCmd(restOptions *api.RESTOptions) *cobra.Command {
	o := DeleteEndPoint{}
	var deleteEndPointCmd = &cobra.Command{
		Use:   "endpoint IP",
		Short: "Delete a LB EndPoint from monitoring",
		Long: `Delete a LB EndPoint from monitoring in the LoxiLB.

ex) loxicmd delete endpoint 31.31.31.31"
		`,
		Aliases: []string{"EndPoint", "ep", "endpoints"},

		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(ctx, time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}

			if val := net.ParseIP(args[0]); val != nil {
				o.Host = args[0]
			} else {
				fmt.Printf("HOSTIP '%s' is invalid format\n", args[0])
				return
			}

			qeury, err := MakeEndPointDeleteQurey(o.Host)
			if err != nil {
				fmt.Printf("Error: Failed to create ep query for delete\n")
				return
			}
			resp, err := client.EndPoint().Query(qeury).Delete(ctx)
			if err != nil {
				fmt.Printf("Error: Failed to delete EndPoint\n")
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

	return deleteEndPointCmd
}

func MakeEndPointDeleteQurey(host string) (map[string]string, error) {
	query := map[string]string{}
	query["hostName"] = host
	return query, nil
}
