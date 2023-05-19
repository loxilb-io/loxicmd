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
package get

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"loxicmd/pkg/api"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

func NewGetMirrorCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetMirrorCmd = &cobra.Command{
		Use:     "mirror",
		Short:   "Get a Mirror",
		Long:    `It shows Mirror Information in the LoxiLB`,
		Aliases: []string{"mirror", "mirr", "mirrors"},

		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}
			resp, err := client.Mirror().SetUrl("/config/mirror/all").Get(ctx)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				PrintGetMirrorResult(resp, *restOptions)
				return
			}

		},
	}

	return GetMirrorCmd
}

func PrintGetMirrorResult(resp *http.Response, o api.RESTOptions) {
	Mirrorresp := api.MirrorGet{}
	var data [][]string
	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}

	if err := json.Unmarshal(resultByte, &Mirrorresp); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	// if json options enable, it print as a json format.
	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(Mirrorresp, "", "    ")
		fmt.Println(string(resultIndent))
		return
	}

	Mirrorresp.Sort()

	// Table Init
	table := TableInit()

	// Making load balance data
	for _, Mirrorrule := range Mirrorresp.Mirrors {
		if o.PrintOption == "wide" {
			table.SetHeader(MIRROR_WIDE_TITLE)
			data = append(data, []string{Mirrorrule.Ident, MakeMirrInfoString(Mirrorrule.Info), fmt.Sprintf("%d", Mirrorrule.Target.AttachMent), Mirrorrule.Target.MirrObjName, fmt.Sprintf("%d", Mirrorrule.Sync)})

		} else {
			table.SetHeader(MIRROR_TITLE)
			data = append(data, []string{Mirrorrule.Ident, MakeMirrInfoString(Mirrorrule.Info), MakeAttachmentToString(Mirrorrule.Target.AttachMent), Mirrorrule.Target.MirrObjName})
		}
	}
	// Rendering the load balance data to table
	TableShow(data, table)
}

func MakeMirrInfoString(infos api.MirrInfo) (ret string) {
	// Print Mirror Type
	if infos.MirrType == 0 {
		ret = "Type : SPAN\n"
	} else if infos.MirrType == 1 {
		ret = "Type : RSPAN\n"
	} else if infos.MirrType == 2 {
		ret = "Type : ERSPAN\n"
	}

	// Add additional information
	if infos.MirrPort != "" {
		ret += fmt.Sprintf("Port : %s\n", infos.MirrPort)
	}
	if infos.MirrRip != "" {
		ret += fmt.Sprintf("Remote IP : %s\n", infos.MirrRip)
	}
	if infos.MirrSip != "" {
		ret += fmt.Sprintf("Source IP : %s\n", infos.MirrSip)
	}
	if infos.MirrTid != 0 {
		ret += fmt.Sprintf("Tunnel ID : %d\n", infos.MirrTid)
	}
	if infos.MirrVlan != 0 {
		ret += fmt.Sprintf("Vlan ID : %d\n", infos.MirrVlan)
	}
	return ret
}

func MakeAttachmentToString(objs api.MirrObjType) (ret string) {
	if objs == 1 {
		ret = "Port"
	} else if objs == 2 {
		ret = "LB Rule"
	}
	return ret
}
