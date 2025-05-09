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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"loxicmd/pkg/api"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type DeleteLoadBalancerResult struct {
	Result string `json:"result"`
}

func validation(args []string) error {
	if len(args) > 1 {
		fmt.Println("delete lb command too many args")
		fmt.Println(args)
	} else if len(args) <= 0 {
		return errors.New("delete lb needs <EXTERNAL-IP> arg")
	}

	return nil
}

func NewDeleteLoadBalancerCmd(restOptions *api.RESTOptions) *cobra.Command {
	var tcpPortNumberList string
	var udpPortNumberList string
	var sctpPortNumberList string
	var icmpPortNumberList bool
	var BGP bool
	var Mark uint16
	var Name string
	var Host string

	var externalIP string
	//var endpointList []string

	var deleteLbCmd = &cobra.Command{
		Use:   "lb <EXTERNAL-IP> [--tcp portNumber] [--udp portNumber] [--sctp portNumber] [--icmp portNumber] [--bgp] [--mark=<val>] [--name=<service-name>] [--host=<url>]",
		Short: "Delete a LoadBalancer",
		Long:  `Delete a LoadBalancer.`,
		PreRun: func(cmd *cobra.Command, args []string) {
			//if len(args) == 0 {
			//	cmd.Help()
			//	os.Exit(0)
			//}
		},
		Run: func(cmd *cobra.Command, args []string) {

			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(ctx, time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}

			if Name != "" {
				subResources := []string{
					"name", Name,
				}
				resp, err := client.LoadBalancer().SubResources(subResources).Delete(ctx)
				if err != nil {
					fmt.Printf("Error: Failed to delete LoadBalancer(Name: %s)\n", Name)
					return
				}
				defer resp.Body.Close()
				fmt.Printf("Debug: response.StatusCode: %d\n", resp.StatusCode)
				if resp.StatusCode == http.StatusOK {
					PrintDeleteResult(resp, *restOptions)
				}
				return
			}

			if err := validation(args); err != nil {
				fmt.Println("not valid EXTERNAL-IP")
				return
			}
			externalIP = args[0]
			PortNumberList := make(map[string][]string)

			if len(tcpPortNumberList) > 0 {
				PortNumberList["tcp"] = strings.Split(tcpPortNumberList, "-")
				if len(PortNumberList["tcp"]) > 2 {
					fmt.Printf("Error: Too many ports in list to delete LoadBalancer(ExternalIP: %s, tcp, Port:%s)\n", externalIP, tcpPortNumberList)
				}
			}
			if len(udpPortNumberList) > 0 {
				PortNumberList["udp"] = strings.Split(udpPortNumberList, "-")
				if len(PortNumberList["udp"]) > 2 {
					fmt.Printf("Error: Too many ports in list to delete LoadBalancer(ExternalIP: %s, udp, Port:%s)\n", externalIP, udpPortNumberList)
				}
			}
			if len(sctpPortNumberList) > 0 {
				PortNumberList["sctp"] = strings.Split(sctpPortNumberList, "-")
				if len(PortNumberList["sctp"]) > 2 {
					fmt.Printf("Error: Too many ports in list to delete LoadBalancer(ExternalIP: %s, sctp, Port:%s)\n", externalIP, sctpPortNumberList)
					return
				}
			}
			if icmpPortNumberList {
				PortNumberList["icmp"] = []string{"0", "0"}
			}
			fmt.Printf("PortNumberList: %v\n", PortNumberList)
			if Host == "" {
				Host = "any"
			}
			for proto, portNum := range PortNumberList {
				sPortMin := portNum[0]
				sPortMax := "0"
				if len(portNum) > 1 {
					sPortMax = portNum[1]
				}
				_, err := strconv.Atoi(sPortMin)
				if err != nil {
					fmt.Printf("Error: Invalid port in list to delete LoadBalancer(ExternalIP: %s,Port:%s)\n", externalIP, sPortMin)
					return
				}
				_, err = strconv.Atoi(sPortMax)
				if err != nil {
					fmt.Printf("Error: Invalid port in list to delete LoadBalancer(ExternalIP: %s,Port:%s)\n", externalIP, sPortMax)
					return
				}
				subResources := []string{
					"hosturl", Host,
					"externalipaddress", externalIP,
					"port", sPortMin,
					"portmax", sPortMax,
					"protocol", proto,
				}
				qmap := map[string]string{}
				qmap["bgp"] = fmt.Sprintf("%v", BGP)
				qmap["block"] = fmt.Sprintf("%v", Mark)
				fmt.Printf("subResources: %v\n", subResources)
				resp, err := client.LoadBalancer().SubResources(subResources).Query(qmap).Delete(ctx)
				if err != nil {
					fmt.Printf("Error: Failed to delete LoadBalancer(ExternalIP: %s, Protocol:%s, Port:%v)\n", externalIP, proto, portNum)
					return
				}
				defer resp.Body.Close()
				fmt.Printf("Debug: response.StatusCode: %d\n", resp.StatusCode)
				if resp.StatusCode == http.StatusOK {
					PrintDeleteResult(resp, *restOptions)
					return
				}
			}
		},
	}

	deleteLbCmd.Flags().StringVar(&tcpPortNumberList, "tcp", tcpPortNumberList, "TCP port list can be specified as '<portMin>-<portMax>'")
	deleteLbCmd.Flags().StringVar(&udpPortNumberList, "udp", udpPortNumberList, "UDP port list can be specified as '<portMin>-<portMax>...'")
	deleteLbCmd.Flags().StringVar(&sctpPortNumberList, "sctp", sctpPortNumberList, "SCTP port list can be specified as '<portMin>-<portMax>...'")
	deleteLbCmd.Flags().BoolVarP(&icmpPortNumberList, "icmp", "", false, "ICMP port list can't be specified'")
	deleteLbCmd.Flags().BoolVarP(&BGP, "bgp", "", false, "BGP enable information'")
	deleteLbCmd.Flags().Uint16VarP(&Mark, "mark", "", 0, "Specify the mark num to segregate a load-balancer VIP service")
	deleteLbCmd.Flags().StringVarP(&Name, "name", "", Name, "Name for load balancer rule")
	deleteLbCmd.Flags().StringVarP(&Host, "host", "", Host, "Ingress Host URL Path")

	return deleteLbCmd
}

func PrintDeleteResult(resp *http.Response, o api.RESTOptions) {
	result := DeleteLoadBalancerResult{}
	resultByte, err := io.ReadAll(resp.Body)
	//fmt.Printf("Debug: response.Body: %s\n", string(resultByte))

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
