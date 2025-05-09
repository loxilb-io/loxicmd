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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"loxicmd/pkg/api"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type CreateLoadBalancerOptions struct {
	ExternalIP     string
	TCP            []string
	UDP            []string
	ICMP           bool
	Mode           string
	BGP            bool
	Security       string
	Monitor        bool
	Attach         bool
	Detach         bool
	Timeout        uint32
	Mark           uint32
	SCTP           []string
	Endpoints      []string
	SecIPs         []string
	Select         string
	Name           string
	Host           string
	AllowedSources []string
	PPv2En         bool
	Egress         bool
}

type CreateLoadBalancerResult struct {
	Result string `json:"result"`
}

const CreateLoadBalancerSuccess = "success"

func ReadCreateLoadBalancerOptions(o *CreateLoadBalancerOptions, args []string) error {
	if len(args) > 1 {
		fmt.Println("create lb command get so many args")
		fmt.Println(args)
	} else if len(args) <= 0 {
		return errors.New("create lb need EXTERNAL-IP args")
	}

	if val := net.ParseIP(args[0]); val != nil {
		o.ExternalIP = args[0]
	} else {
		return fmt.Errorf("externel IP '%s' is invalid format", args[0])
	}
	return nil
}

func SelectToNum(sel string) int {
	var ret int
	switch sel {
	case "rr":
		ret = 0
	case "hash":
		ret = 1
	case "priority":
		ret = 2
	case "persist":
		ret = 3
	case "lc":
		ret = 4
	case "n2":
		ret = 5
	case "n3":
		ret = 6
	default:
		ret = 0
	}
	return ret
}

func ModeToNum(sel string) int {
	var ret int
	switch sel {
	case "onearm":
		ret = 1
	case "fullnat":
		ret = 2
	case "dsr":
		ret = 3
	case "fullproxy":
		ret = 4
	case "hostonearm":
		ret = 5
	default:
		ret = 0
	}
	return ret
}

func SecStringToNum(sec string) int {
	var ret int
	switch sec {
	case "https":
		ret = 1
	case "tls":
		ret = 1
	case "e2ehttps":
		ret = 2
	case "e2etls":
		ret = 2
	default:
		ret = 0
	}
	return ret
}

func NewCreateLoadBalancerCmd(restOptions *api.RESTOptions) *cobra.Command {
	o := CreateLoadBalancerOptions{}

	var createLbCmd = &cobra.Command{
		Use:   "lb IP [--select=<rr|hash|priority|persist>] [--tcp=<ports>:<targetPorts>] [--udp=<ports>:<targetPorts>] [--sctp=<ports>:<targetPorts>] [--icmp] [--mark=<val>] [--secips=<ip>,] [--sources=<ip>,] [--endpoints=<ip>:<weight>,] [--mode=<onearm|fullnat>] [--bgp] [--monitor] [--inatimeout=<to>] [--name=<service-name>] [--attachEP] [--detachEP] [--security=<https|e2ehttps|none>] [--host=<url>] [--ppv2en] [--egress]",
		Short: "Create a LoadBalancer",
		Long: `Create a LoadBalancer

--select value options
  	rr - select the lb end-points based on round-robin
	hash - select the lb end-points based on hashing
	priority - select the lb based on weighted round-robin
	persist - select the lb end-point based on sender
	n2 - select the lb end-point base on N2 interface params (only available with fullproxy mode)
	n3 - select the lb end-point base on N3 interface params
	lc - select the lb end-point base on least connection

--mode value options
	onearm - LB put LB-IP as srcIP
	fullnat - LB put Service IP as scrIP
	dsr - LB in DSR mode allows return traffic to bypass the load balancer (only available with hash select)
	fullproxy - LB operating as a L7 proxy
	hostonearm - LB operating in host one-arm

ex)
	loxicmd create lb 192.168.0.200 --tcp=80:32015 --endpoints=10.212.0.1:1,10.212.0.2:1,10.212.0.3:1
	loxicmd create lb 192.168.0.200 --tcp=8080-8081:32015 --endpoints=10.212.0.1:1,10.212.0.2:1,10.212.0.3:1
	loxicmd create lb 192.168.0.200 --tcp=5000:5201-5300 --endpoints=10.212.0.1:1,10.212.0.2:1,10.212.0.3:1
	loxicmd create lb 192.168.0.200 --tcp=80:32015 --endpoints=10.212.0.1:1,10.212.0.2:1,10.212.0.3:1 --security=https
	loxicmd create lb 192.168.0.200 --tcp=80:32015 --endpoints=10.212.0.1:1,10.212.0.2:1,10.212.0.3:1 --host=loxilb.io
	loxicmd create lb 192.168.0.200 --tcp=80:32015 --name="http-service" --endpoints=10.212.0.1:1,10.212.0.2:1,10.212.0.3:1
	loxicmd create lb 192.168.0.200 --udp=80:32015 --endpoints=10.212.0.1:1,10.212.0.2:1,10.212.0.3:1 --mark=10
	loxicmd create lb 192.168.0.200 --tcp=80:32015 --udp=80:32015 --endpoints=10.212.0.1:1,10.212.0.2:1,10.212.0.3:1
	loxicmd create lb 192.168.0.200 --select=hash --tcp=80:32015 --endpoints=10.212.0.1:1,10.212.0.2:1,10.212.0.3:1
	loxicmd create lb 192.168.0.200 --tcp=80:80 --endpoints=10.212.0.1:1,10.212.0.2:1,10.212.0.3:1 --mode=dsr --select=hash
	loxicmd create lb 192.168.0.200 --sctp=37412:38412 --secips=192.168.0.201,192.168.0.202 --endpoints=10.212.0.1:1,10.212.0.2:1,10.212.0.3:1
	loxicmd create lb 192.168.0.200 --tcp=80:32015 --endpoints=10.212.0.1:1,10.212.0.2:1,10.212.0.3:1 --sources=10.10.10.1/32

	loxicmd create lb  2001::1 --tcp=2020:8080 --endpoints=4ffe::1:1,5ffe::1:1,6ffe::1:1
	loxicmd create lb  2001::1 --tcp=2020:8080 --endpoints=31.31.31.1:1,32.32.32.1:1,33.33.33.1:1
	loxicmd create lb 10.10.10.254 --sctp=2020:8080 --endpoints=33.33.33.1:1 --attachEP
	loxicmd create lb 100.100.100.1 --tcp=8080:80 --endpoints=10.10.10.1:1 --ppv2en
	`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			var sctp bool
			if err := ReadCreateLoadBalancerOptions(&o, args); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}

			ProtoPortpair := make(map[string][]string)
			// TCP LoadBalancer
			if len(o.TCP) > 0 {
				ProtoPortpair["tcp"] = o.TCP
			}
			if len(o.UDP) > 0 {
				ProtoPortpair["udp"] = o.UDP
			}
			if len(o.SCTP) > 0 {
				ProtoPortpair["sctp"] = o.SCTP
				sctp = true
			}
			if o.ICMP {
				//icmpProtoPortpair := make(map[string][]string)
				ProtoPortpair["icmp"] = []string{"0:0"}
			}
			if !sctp && len(o.SecIPs) > 0 {
				fmt.Printf("Secondary IPs allowed in SCTP only\n")
				return
			}

			fmt.Printf("ProtoPortpair: %v\n", ProtoPortpair)
			// Commom Part of the load balancer.
			endpointPair, err := GetEndpointWeightPairList(o.Endpoints)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			for proto, portPairList := range ProtoPortpair {
				portTargetPorts, err := GetPortPairList(portPairList)
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
					return
				}
				if len(portTargetPorts) <= 0 || len(portTargetPorts) > 2 {
					fmt.Printf("portPair: None specified\n")
					return
				}

				startSPort := uint16(0)
				endSPort := uint16(0)
				first := false

				for sPort := range portTargetPorts {
					if !first {
						startSPort = sPort
						first = true
					} else {
						if sPort > startSPort {
							endSPort = sPort
						} else {
							endSPort = startSPort
							startSPort = sPort
						}
					}
				}

				lbModel := api.LoadBalancerModel{}
				oper := 0
				if o.Attach {
					oper = 1
				} else if o.Detach {
					oper = 2
				}
				lbService := api.LoadBalancerService{
					ExternalIP: o.ExternalIP,
					Protocol:   proto,
					Port:       startSPort,
					PortMax:    endSPort,
					Sel:        api.EpSelect(SelectToNum(o.Select)),
					BGP:        o.BGP,
					Monitor:    o.Monitor,
					Mode:       api.LbMode(ModeToNum(o.Mode)),
					Timeout:    o.Timeout,
					Block:      o.Mark,
					Name:       o.Name,
					Oper:       api.LbOP(oper),
					Security:   api.LbSec(SecStringToNum(o.Security)),
					Host:       o.Host,
					PpV2:       o.PPv2En,
					Egress:     o.Egress,
				}

				lbModel.Service = lbService
				for endpoint, weight := range endpointPair {
					targetPorts := portTargetPorts[startSPort]
					for _, targetPort := range targetPorts {
						if o.Mode == "dsr" && targetPort != startSPort {
							fmt.Printf("Error: No port-translation in dsr mode\n")
							return
						}
						ep := api.LoadBalancerEndpoint{
							EndpointIP: endpoint,
							TargetPort: targetPort,
							Weight:     weight,
						}
						lbModel.Endpoints = append(lbModel.Endpoints, ep)
					}
				}

				for _, sip := range o.SecIPs {
					sp := api.LoadBalancerSecIp{
						SecondaryIP: sip,
					}
					lbModel.SecondaryIPs = append(lbModel.SecondaryIPs, sp)
				}

				for _, sip := range o.AllowedSources {
					sp := api.LbAllowedSrcIPArg{
						Prefix: sip,
					}
					lbModel.SrcIPs = append(lbModel.SrcIPs, sp)
				}

				resp, err := LoadbalancerAPICall(restOptions, lbModel)
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
					return
				}

				defer resp.Body.Close()

				//fmt.Printf("Debug: request: %v\n", lbModel)

				fmt.Printf("Debug: response.StatusCode: %d\n", resp.StatusCode)
				if resp.StatusCode == http.StatusOK {
					PrintCreateResult(resp, *restOptions)
					return
				}
			}
		},
	}

	createLbCmd.Flags().StringSliceVar(&o.TCP, "tcp", o.TCP, "Port pairs can be specified as '<port>:<targetPort>'")
	createLbCmd.Flags().StringSliceVar(&o.UDP, "udp", o.UDP, "Port pairs can be specified as '<port>:<targetPort>'")
	createLbCmd.Flags().StringSliceVar(&o.SCTP, "sctp", o.SCTP, "Port pairs can be specified as '<port>:<targetPort>'")
	createLbCmd.Flags().BoolVarP(&o.ICMP, "icmp", "", false, "ICMP Ping packet Load balancer")
	createLbCmd.Flags().StringVarP(&o.Mode, "mode", "", o.Mode, "NAT mode for load balancer rule")
	createLbCmd.Flags().BoolVarP(&o.BGP, "bgp", "", false, "Enable BGP in the load balancer")
	createLbCmd.Flags().BoolVarP(&o.Monitor, "monitor", "", false, "Enable monitoring end-points of this rule")
	createLbCmd.Flags().StringSliceVar(&o.SecIPs, "secips", o.SecIPs, "Secondary IPs for SCTP multihoming rule specified as '<secondaryIP>'")
	createLbCmd.Flags().StringVarP(&o.Select, "select", "", "rr", "Select the hash algorithm for the load balance.(ex) rr, hash, priority, persist, lc")
	createLbCmd.Flags().Uint32VarP(&o.Timeout, "inatimeout", "", 0, "Specify the timeout (in seconds) after which a LB session will be reset for inactivity")
	createLbCmd.Flags().Uint32VarP(&o.Mark, "mark", "", 0, "Specify the mark num to segregate a load-balancer VIP service")
	createLbCmd.Flags().StringSliceVar(&o.Endpoints, "endpoints", o.Endpoints, "Endpoints is pairs that can be specified as '<endpointIP>:<Weight>'")
	createLbCmd.Flags().StringVarP(&o.Name, "name", "", o.Name, "Name for load balancer rule")
	createLbCmd.Flags().BoolVarP(&o.Attach, "attachEP", "", false, "Attach endpoints to the load balancer rule")
	createLbCmd.Flags().BoolVarP(&o.Detach, "detachEP", "", false, "Detach endpoints from the load balancer rule")
	createLbCmd.Flags().StringVarP(&o.Security, "security", "", o.Security, "Security mode for load balancer rule")
	createLbCmd.Flags().StringVarP(&o.Host, "host", "", o.Host, "Ingress Host URL Path")
	createLbCmd.Flags().StringSliceVar(&o.AllowedSources, "sources", o.AllowedSources, "Allowed sources for this rule as '<allowedSources>'")
	createLbCmd.Flags().BoolVarP(&o.PPv2En, "ppv2en", "", false, "Enable proxy procotol v2")
	createLbCmd.Flags().BoolVarP(&o.Egress, "egress", "", false, "Specify egress rule")

	return createLbCmd
}

func PrintCreateResult(resp *http.Response, o api.RESTOptions) {
	result := CreateLoadBalancerResult{}
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
		// TODO: need to test MarshalIndent
		resultIndent, _ := json.MarshalIndent(resp.Body, "", "\t")
		fmt.Println(string(resultIndent))
		return
	}

	fmt.Printf("%s\n", result.Result)
}

func GetPortPairList(portPairStrList []string) (map[uint16][]uint16, error) {
	result := make(map[uint16][]uint16)
	for _, portPairStr := range portPairStrList {
		portPair := strings.Split(portPairStr, ":")
		if len(portPair) != 2 {
			continue
		}

		servicePorts := strings.Split(portPair[0], "-")

		// 0 is port, 1 is targetPort
		var portList []int
		for _, servicePort := range servicePorts {
			port, err := strconv.Atoi(servicePort)
			if err != nil {
				return nil, fmt.Errorf("port '%s' is not integer", servicePort)
			}
			portList = append(portList, port)
		}

		var err error
		startTP := 0
		endTP := 0

		targetPortRange := strings.Split(portPair[1], "-")
		if len(targetPortRange) > 2 {
			continue
		} else if len(targetPortRange) == 2 {
			startTP, err = strconv.Atoi(targetPortRange[0])
			if err != nil {
				return nil, fmt.Errorf("targetPort0 '%s' is not integer", targetPortRange[0])
			}
			endTP, err = strconv.Atoi(targetPortRange[1])
			if err != nil {
				return nil, fmt.Errorf("targetPort1 '%s' is not integer", targetPortRange[1])
			}
			if endTP < startTP {
				return nil, fmt.Errorf("targetPort2 '%s' <  targetPort2 '%s'", targetPortRange[1], targetPortRange[0])
			}
		} else {
			startTP, err = strconv.Atoi(targetPortRange[0])
			if err != nil {
				return nil, fmt.Errorf("targetPort0 '%s' is not integer", targetPortRange[0])
			}
			endTP = startTP
		}

		for _, port := range portList {
			result[uint16(port)] = make([]uint16, 0)

			for targetPort := startTP; targetPort <= endTP; targetPort++ {
				result[uint16(port)] = append(result[uint16(port)], uint16(targetPort))
			}
		}
	}
	return result, nil
}

func GetEndpointWeightPairList(endpointsList []string) (map[string]uint8, error) {
	result := make(map[string]uint8)
	for _, endpointStr := range endpointsList {
		weightIdx := strings.LastIndex(endpointStr, ":")
		if weightIdx < 0 {
			return nil, fmt.Errorf("endpoint '%s' is invalid format", endpointStr)
		}
		// 0 is endpoint IP, 1 is weight
		weight, err := strconv.Atoi(endpointStr[weightIdx+1:])
		if err != nil {
			return nil, fmt.Errorf("endpoint's weight '%s' is invalid format", endpointStr[weightIdx+1:])
		}
		result[endpointStr[:weightIdx]] = uint8(weight)
	}

	return result, nil
}

func LoadbalancerAPICall(restOptions *api.RESTOptions, lbModel api.LoadBalancerModel) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}

	return client.LoadBalancer().Create(ctx, lbModel)
}
