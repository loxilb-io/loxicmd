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
	"os"
	"time"

	"loxicmd/pkg/api"

	"github.com/spf13/cobra"
)

func DeleteFDBValidation(args []string) error {
	if len(args) > 3 {
		fmt.Println("delete FDB command get so many args")
	} else if len(args) <= 1 {
		return errors.New("delete IP Address need <MacAddress> <device> args")
	}
	if _, err := net.ParseMAC(args[0]); err != nil {
		return fmt.Errorf("MacAddress '%s' is invalid format", args[0])
	}

	return nil
}

func NewDeleteFDBCmd(restOptions *api.RESTOptions) *cobra.Command {

	var deleteFDBCmd = &cobra.Command{
		Use:   "fdb <MacAddress> <DeviceName>",
		Short: "Delete a FDB",
		Long:  `Delete a FDB using MacAddress  in the LoxiLB.`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := DeleteFDBValidation(args); err != nil {
				fmt.Println("not valid <MacAddress>")
				return
			}
			MacAddress := args[0]
			Device := args[1]

			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(ctx, time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}
			subResources := []string{
				MacAddress, "dev", Device,
			}
			resp, err := client.FDB().SubResources(subResources).Delete(ctx)
			if err != nil {
				fmt.Printf("Error: Failed to delete FDB : %s", MacAddress)
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

	return deleteFDBCmd
}
