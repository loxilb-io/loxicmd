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
	"time"

	"github.com/spf13/cobra"
)

func NewCreateVxlanPeerCmd(restOptions *api.RESTOptions) *cobra.Command {
	var createvxlanCmd = &cobra.Command{
		Use:   "vxlanpeer <Vnid> <PeerIP>",
		Short: "Create a vxlan",
		Long: `Create a vxlan using LoxiLB.

ex) loxicmd create vxlan-peer 100 30.1.3.1
`,
		Aliases: []string{"vxlanPeer", "vxlan-peer", "vxlan_peer"},
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			var vxlanMod api.VxlanPeerMod
			// Make vxlanMod
			if err := ReadCreateVxlanPeerOptions(&vxlanMod, args); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			url := fmt.Sprintf("/config/tunnel/vxlan/%s/peer", args[0])
			resp, err := VxlanPeerAPICall(restOptions, vxlanMod, url)
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

	return createvxlanCmd
}

func ReadCreateVxlanPeerOptions(o *api.VxlanPeerMod, args []string) error {
	if len(args) > 2 {
		return errors.New("create vxlan command get so many args")
	} else if len(args) < 1 {
		return errors.New("create vxlan need <MacAddress>  args")
	}

	if val := net.ParseIP(args[1]); val == nil {
		return fmt.Errorf("PeerIP '%s' is invalid format", args[1])
	}
	o.PeerIP = args[1]

	return nil
}

func VxlanPeerAPICall(restOptions *api.RESTOptions, vxlanModel api.VxlanPeerMod, url string) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}

	return client.Vxlan().SetUrl(url).Create(ctx, vxlanModel)
}
