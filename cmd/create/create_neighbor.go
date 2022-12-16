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
	"time"

	"github.com/spf13/cobra"
)

type CreateNeighborOptions struct {
	macAddress string
}

func NewCreateNeighborsCmd(restOptions *api.RESTOptions) *cobra.Command {
	o := CreateNeighborOptions{}
	var createNeighborsCmd = &cobra.Command{
		Use:   "neighbor <DeviceIP> <DeviceName> [--macAddress=aa:aa:aa:aa:aa:aa]",
		Short: "Create a Neighbors",
		Long: `Create a Neighbors using LoxiLB. It is working as "ip neigh add <DeviceIP> dev <device> lladdr <--macAddress>"

ex) loxicmd create neighbor 192.168.0.1 eno7 --macAddress=aa:aa:aa:aa:aa:aa
`,
		Aliases: []string{"nei", "neigh"},
		Run: func(cmd *cobra.Command, args []string) {
			var NeighborsMod api.NeighborMod
			// Make NeighborsMod
			if err := ReadCreateNeighborsOptions(&NeighborsMod, args); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			NeighborsMod.MacAddress = o.macAddress
			resp, err := NeighborsAPICall(restOptions, NeighborsMod)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				PrintCreateLbResult(resp, *restOptions)
				return
			}

		},
	}
	createNeighborsCmd.Flags().StringVarP(&o.macAddress, "macAddress", "", "", "Hardware MAC address")
	return createNeighborsCmd
}

func ReadCreateNeighborsOptions(o *api.NeighborMod, args []string) error {
	if len(args) > 4 {
		return errors.New("create Neighbors command get so many args")
	} else if len(args) < 2 {
		return errors.New("create Neighbors need <DeviceIPNet>  args")
	}

	if val := net.ParseIP(args[0]); val == nil {
		return fmt.Errorf("DeviceIP '%s' is invalid format", args[0])
	}
	o.IP = args[0]
	o.Dev = args[1]

	return nil
}

func NeighborsAPICall(restOptions *api.RESTOptions, NeighborsModel api.NeighborMod) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}

	return client.Neighbor().Create(ctx, NeighborsModel)
}
