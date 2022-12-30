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
	"fmt"
	"loxicmd/pkg/api"
	"net/http"
	"time"
	"net"

	"github.com/spf13/cobra"
)

type CreateEndPointOptions struct {
	Host		  string
	Description   string
	ProbeType     string
	ProbeReq      string
	ProbeResp     string
	ProbePort     int
	ProbeDuration int
	ProbeReTries  int
}

func NewCreateEndPointCmd(restOptions *api.RESTOptions) *cobra.Command {
	o := CreateEndPointOptions{}

	var createEndPointCmd = &cobra.Command{
		Use:   "endpoint IP [--desc=<desc>] [--probetype=<probetype>] [--req=<probereq>] [--resp=<proberesp>] [--l4port=<port>] [--period=<period>] [--retries=<retries>]",
		Short: "Create a LB EndPoint for monitoring",
		Long: `Create a LB EndPoint for monitoring using LoxiLB

ex) loxicmd create endpoint 32.32.32.1 --desc=zone1host --probetype=http --l4port=8080 --period=60 --retries=2
`,
		Aliases: []string{"Endpoint", "ep", "endpoints"},

		Run: func(cmd *cobra.Command, args []string) {
			var EPMod api.EndPointMod
			// Make EndPointMod
			if len(args) <= 0 {
				fmt.Printf("create ep need HOST-IP args\n")
				return
			}

			if val := net.ParseIP(args[0]); val != nil {
				o.Host = args[0]
			} else {
				fmt.Printf("HOSTIP '%s' is invalid format\n", args[0])
				return
			}

			if o.ProbeType != "http" && o.ProbeType != "ping" && o.ProbeType != "connect-tcp" &&
				o.ProbeType != "connect-udp" && o.ProbeType != "connect-sctp" {
				fmt.Printf("probetype '%s' is invalid\n", o.ProbeType)
				return
			}

			if o.ProbeType == "http"  || o.ProbeType == "connect-tcp" ||
				o.ProbeType == "connect-udp" || o.ProbeType == "connect-sctp" {
				if o.ProbePort == 0 {
					fmt.Printf("probeport cant be 0 for '%s' probes\n", o.ProbeType)
					return
				}
			}

			if o.ProbeDuration > 24*60*60 {
				fmt.Printf("probe period is out of bounds\n")
				return
			}

			EPMod.Name = o.Host
			EPMod.Desc = o.Description
			EPMod.ProbeDuration = uint32(o.ProbeDuration)
			EPMod.ProbePort = uint16(o.ProbePort)
			EPMod.ProbeType = o.ProbeType
			EPMod.InActTries = o.ProbeReTries
			EPMod.ProbeReq = o.ProbeReq
			EPMod.ProbeResp = o.ProbeResp

			resp, err := EndPointAPICall(restOptions, EPMod)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			defer resp.Body.Close()

			fmt.Printf("Debug: response.StatusCode: %d\n", resp.StatusCode)
			if resp.StatusCode != http.StatusOK {
				PrintCreateLbResult(resp, *restOptions)
				return
			}
		},
	}

	createEndPointCmd.Flags().StringVar(&o.Description, "desc", "", "Description of and end-point")
	createEndPointCmd.Flags().StringVar(&o.ProbeType, "probetype", "ping", "Probe-type:ping,http,connect-udp,connect-tcp, connect-sctp")
	createEndPointCmd.Flags().StringVar(&o.ProbeReq, "probereq", "", "If probe is http, one can specify additional uri path")
	createEndPointCmd.Flags().StringVar(&o.ProbeResp, "proberesp", "", "If probe is http, one can specify custom response string")
	createEndPointCmd.Flags().IntVar(&o.ProbePort, "l4port", 0, "If probe is http,tcp,udp,sctp one can specify custom l4port to use")
	createEndPointCmd.Flags().IntVar(&o.ProbeDuration, "period", 60, "Period of probing")
	createEndPointCmd.Flags().IntVar(&o.ProbeReTries, "retries", 2, "Number of retries before marking endPoint inactive")

	return createEndPointCmd
}

func EndPointAPICall(restOptions *api.RESTOptions, EndPointModel api.EndPointMod) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}

	return client.EndPoint().Create(ctx, EndPointModel)
}
