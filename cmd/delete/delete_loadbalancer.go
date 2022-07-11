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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"loxicmd/pkg/api"

	"github.com/spf13/cobra"
)

type DeleteLoadBalancerResult struct {
	Result string `json:"result"`
}

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
		Long:  `Delete a LoadBalancer.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := validation(args); err != nil {
				fmt.Println("not valid EXTERNAL-IP")
				return
			}
			externalIP = args[0]

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
				resp, err := client.LoadBalancer().SubResources(subResources).Delete(ctx)
				if err != nil {
					fmt.Printf("Error: Failed to delete LoadBalancer(ExternalIP: %s, Protocol:%s, Port:%d)", externalIP, "tcp", portNum)
					return
				}
				defer resp.Body.Close()
				fmt.Printf("Debug: response.StatusCode: %d\n", resp.StatusCode)
				if resp.StatusCode != http.StatusOK {
					PrintDeleteLbResult(resp, *restOptions)
					return
				}
			}
			return
		},
	}

	deleteLbCmd.Flags().IntSliceVar(&tcpPortNumberList, "tcp", tcpPortNumberList, "TCP port list can be specified as '<port>,<port>...'")
	return deleteLbCmd
}

func PrintDeleteLbResult(resp *http.Response, o api.RESTOptions) {
	result := DeleteLoadBalancerResult{}
	resultByte, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Debug: response.Body: %s\n", string(resultByte))

	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}
	if err := json.Unmarshal(resultByte, &result); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(resp.Body, "", "\t")
		fmt.Println(string(resultIndent))
		return
	}

	fmt.Printf("%s\n", result.Result)
}
