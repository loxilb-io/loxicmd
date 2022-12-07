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
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"loxicmd/pkg/api"

	"github.com/spf13/cobra"
)

func DeleteIPv4AddressValidation(args []string) error {
	if len(args) > 3 {
		fmt.Println("delete IPv4Address command get so many args")
	} else if len(args) <= 1 {
		return errors.New("delete IP Address need <DeviceIPNet> <device> args")
	}
	if _, _, err := net.ParseCIDR(args[0]); err != nil {
		return fmt.Errorf("DeviceIPNet '%s' is invalid format", args[0])
	}

	return nil
}

func NewDeleteIPv4AddressCmd(restOptions *api.RESTOptions) *cobra.Command {

	var deleteIPv4AddressCmd = &cobra.Command{
		Use:     "ip <DeviceIPNet> <device>",
		Short:   "Delete a IPv4Address",
		Long:    `Delete a IPv4Address using DeviceIPNet  in the LoxiLB.`,
		Aliases: []string{"ipv4address", "ipv4", "ipaddress"},
		Run: func(cmd *cobra.Command, args []string) {
			if err := DeleteIPv4AddressValidation(args); err != nil {
				fmt.Println("not valid <DeviceIPNet>")
				return
			}
			DeviceIPNet := args[0]
			Device := args[1]

			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(ctx, time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}
			subResources := []string{
				DeviceIPNet, "dev", Device,
			}
			resp, err := client.IPv4Address().SubResources(subResources).Delete(ctx)
			if err != nil {
				fmt.Printf("Error: Failed to delete IPv4Address : %s", DeviceIPNet)
				return
			}
			defer resp.Body.Close()
			fmt.Printf("Debug: response.StatusCode: %d\n", resp.StatusCode)
			if resp.StatusCode != http.StatusOK {
				PrintDeleteLbResult(resp, *restOptions)
				return
			}

		},
	}

	return deleteIPv4AddressCmd
}
