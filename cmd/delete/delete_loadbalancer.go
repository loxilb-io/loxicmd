/*
Copyright © 2022 Baekgyun Jung <backguyn@netlox.io>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package delete

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"loxicmd/pkg/api"

	"github.com/spf13/cobra"
)

func validation(args []string) error {
	if len(args) > 1 {
		fmt.Println("create lb command get so many args")
		fmt.Println(args)
	} else if len(args) <= 0 {
		return errors.New("delete lb need EXTERNAL-IP args")
	}

	// TODO: need validation check
	return nil
}

func NewDeleteLoadBalancerCmd(restOptions *api.RESTOptions) *cobra.Command {
	var tcpPortNumberList []int
	var externalIP string
	//var endpointList []string

	var deleteLbCmd = &cobra.Command{
		Use:   "lb EXTERNAL-IP [--tcp=<port>:<targetPort>] [--endpoints=<ip>:<weight>]",
		Short: "Delete a LoadBalancer",
		Long: `Delete a LoadBalancer. what the hell!!
	what the hell!!`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := validation(args); err != nil {
				fmt.Println("not valid EXTERNAL-IP")
				return
			}
			externalIP = args[0]

			fmt.Println("delete lb called")
			fmt.Printf("ExternalIP: %s, tcp: %v", externalIP, tcpPortNumberList)

			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(ctx, time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}

			for _, portNum := range tcpPortNumberList {
				// TODO: need validation check
				subResources := []string{
					"externalipaddress", externalIP,
					"port", strconv.Itoa(portNum),
					"protocol", "tcp",
				}
				err := client.LoadBalancer().SubResources(subResources).Delete(ctx)
				if err != nil {
					fmt.Printf("Error: Failed to delete LoadBalancer(ExternalIP: %s, Protocol:%s, Port:%d)", externalIP, "tcp", portNum)
					return
				}
			}
			return
		},
	}

	deleteLbCmd.Flags().IntSliceVar(&tcpPortNumberList, "tcp", tcpPortNumberList, "TCP port list can be specified as '<port>,<port>...'")
	return nil
}
