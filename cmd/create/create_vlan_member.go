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
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

type CreateVlanMemberOptions struct {
	Tagged bool
}

func NewCreateVlanMemberCmd(restOptions *api.RESTOptions) *cobra.Command {
	o := CreateVlanMemberOptions{}
	var createvlanCmd = &cobra.Command{
		Use:   "vlanmember <Vid> <DeviceName> --tagged=<Tagged>",
		Short: "Create a vlanmember",
		Long: `Create a vlanmember using LoxiLB. It is working as "brctl addif vlan<Vid> <DeviceName>.<tagged>"

		
ex) loxicmd create vlanmember 100 eno7 --tagged=true
	loxicmd create vlanmember 100 eno7 
`,
		Aliases: []string{"vlanMember", "vlan-member", "vlan_member"},
		Run: func(cmd *cobra.Command, args []string) {
			var vlanMod api.VlanMemberMod
			// Make vlanMod
			if err := ReadCreateVlanMemberOptions(&vlanMod, args); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			// Args Setting
			url := fmt.Sprintf("/config/vlan/%s/member", args[0])
			vlanMod.Dev = args[1]
			vlanMod.Tagged = o.Tagged

			resp, err := VlanMemberAPICall(restOptions, vlanMod, url)
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
	createvlanCmd.Flags().BoolVarP(&o.Tagged, "tagged", "", false, "Tagged mode Vlan")
	return createvlanCmd
}

func ReadCreateVlanMemberOptions(o *api.VlanMemberMod, args []string) error {
	if len(args) > 3 {
		return errors.New("create vlan member command get so many args")
	} else if len(args) <= 1 {
		return errors.New("create vlan member need <Vid> , <Device> args")
	}
	_, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	return nil
}

func VlanMemberAPICall(restOptions *api.RESTOptions, vlanModel api.VlanMemberMod, url string) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}

	return client.Vlan().SetUrl(url).Create(ctx, vlanModel)
}
