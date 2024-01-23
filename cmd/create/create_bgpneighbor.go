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
	"time"

	"github.com/spf13/cobra"
)

type CreateBGPNeighborOptions struct {
	SetMultiHtop bool
	RemotePort   uint8
}

func NewCreateBGPNeighborCmd(restOptions *api.RESTOptions) *cobra.Command {
	o := CreateBGPNeighborOptions{}

	var createBGPNeighborCmd = &cobra.Command{
		Use:   "bgpneighbor <PeerIP> <ASN> [--remotePort=<remoteBGPPort>] [--setMultiHtop]",
		Short: "Create a BGP Neighbor",
		Long: `Create a BGP Neighbor using LoxiLB.
ex) loxicmd create bgpneighbor 10.10.10.1 64512
`,
		Aliases: []string{"bgpnei", "bgpneigh"},
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			var BGPNeighborMod api.BGPNeighborMod
			// Make BGPNeighborMod
			if err := ReadCreateBGPNeighborOptions(&BGPNeighborMod, args); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			// option
			if o.RemotePort != 179 {
				BGPNeighborMod.RemotePort = int(o.RemotePort)
			}
			if o.SetMultiHtop {
				BGPNeighborMod.SetMultiHop = o.SetMultiHtop
			}
			resp, err := BGPNeighborAPICall(restOptions, BGPNeighborMod)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				PrintCreateResult(resp, *restOptions)
				return
			}

		},
	}
	createBGPNeighborCmd.Flags().BoolVarP(&o.SetMultiHtop, "setMultiHtop", "", false, "Enable Multihop BGP in the load balancer")
	createBGPNeighborCmd.Flags().Uint8VarP(&o.RemotePort, "remotePort", "", 179, "BGP Port number of the remote site")

	return createBGPNeighborCmd
}

func ReadCreateBGPNeighborOptions(o *api.BGPNeighborMod, args []string) error {
	if len(args) > 2 {
		return errors.New("create BGPNeighbor command get so many args")
	} else if len(args) <= 1 {
		return errors.New("create BGPNeighbor need <RemoteAS> args")
	}

	if val := net.ParseIP(args[0]); val != nil {
		o.IPaddress = args[0]
	} else {
		return fmt.Errorf("Peer IP '%s' is invalid format", args[0])
	}

	RemoteAs, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}
	o.RemoteAs = RemoteAs
	return nil
}

func BGPNeighborAPICall(restOptions *api.RESTOptions, BGPNeighborModel api.BGPNeighborMod) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}

	return client.BGPNeighbor().Create(ctx, BGPNeighborModel)
}
