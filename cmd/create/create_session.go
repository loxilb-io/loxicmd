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
	"errors"
	"fmt"
	"loxicmd/pkg/api"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type CreateSessionOptions struct {
	UserID    string
	SessionIP string
	ANTunnel  string
	CNTunnel  string
}

func NewCreateSessionCmd(restOptions *api.RESTOptions) *cobra.Command {
	o := CreateSessionOptions{}

	var createSessionCmd = &cobra.Command{
		Use:   "session <userID> <sessionIP> --accessNetworkTunnel=<TeID>:<TunnelIP> --coreNetworkTunnel=<TeID>:<TunnelIP>",
		Short: "Create a Session",
		Long: `Create a Session using LoxiLB
		
ex) loxicmd create session user1 192.168.20.1 --accessNetworkTunnel=1:1.232.16.1 coreNetworkTunnel=1:1.233.16.1

		`,
		Aliases: []string{"session", "sessions"},
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			var SessionMod api.SessionMod
			// Make SessionMod
			if err := ReadCreateSessionOptions(&SessionMod, args); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if err := GetNetworkTunnelPairList(&SessionMod, o.ANTunnel, true); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if err := GetNetworkTunnelPairList(&SessionMod, o.CNTunnel, false); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}

			resp, err := SessionAPICall(restOptions, SessionMod)
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

	createSessionCmd.Flags().StringVarP(&o.ANTunnel, "accessNetworkTunnel", "", "", "accessNetworkTunnel has pairs that can be specified as '<TeID>:<IP>'")
	createSessionCmd.Flags().StringVarP(&o.CNTunnel, "coreNetworkTunnel", "", "", "coreNetworkTunnel has pairs that can be specified as '<TeID>:<IP>'")

	return createSessionCmd
}

func ReadCreateSessionOptions(o *api.SessionMod, args []string) error {
	if len(args) > 3 {
		return errors.New("create Session command get so many args")
	} else if len(args) <= 1 {
		return errors.New("create Session need <userID> and <sessionIP> args")
	}

	o.Ident = args[0]
	if val := net.ParseIP(args[1]); val != nil {
		o.Ip = val
	} else {
		return fmt.Errorf("Session IP '%s' is invalid format", args[1])
	}
	return nil
}

func GetNetworkTunnelPairList(o *api.SessionMod, networkTunnel string, An bool) error {
	networkTunnelPair := strings.Split(networkTunnel, ":")
	if len(networkTunnelPair) != 2 {
		return fmt.Errorf("NetworkTunnel '%s' is invalid format", networkTunnel)
	}

	// 0 is TeID, 1 is TunnelIP
	TeID, err := strconv.Atoi(networkTunnelPair[0])
	if err != nil {
		return fmt.Errorf("NetworkTunnel's TeID '%s' is invalid format", networkTunnelPair[0])
	}
	if An {
		o.AnTun.TeID = uint32(TeID)
		if val := net.ParseIP(networkTunnelPair[1]); val != nil {
			o.AnTun.Addr = val
		} else {
			return fmt.Errorf("Tunnel IP '%s' is invalid format", networkTunnelPair[1])
		}
	} else {
		o.CnTun.TeID = uint32(TeID)
		if val := net.ParseIP(networkTunnelPair[1]); val != nil {
			o.CnTun.Addr = val
		} else {
			return fmt.Errorf("Tunnel IP '%s' is invalid format", networkTunnelPair[1])
		}
	}

	return nil
}

func SessionAPICall(restOptions *api.RESTOptions, sessionModel api.SessionMod) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}

	return client.Session().Create(ctx, sessionModel)
}
