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
	"net/http"
	"os"
	"time"

	"loxicmd/pkg/api"

	"github.com/spf13/cobra"
)

type DeleteSessionUlClOptions struct {
	UserID   string
	UlClArgs []string
}

func DeleteSessionUlClValidation(args []string) error {
	if len(args) > 1 {
		fmt.Println("delete Session command get so many args")
		fmt.Println(args)
	} else if len(args) <= 0 {
		return errors.New("delete Session need <UserID> args")
	}

	return nil
}

func NewDeleteSessionUlClCmd(restOptions *api.RESTOptions) *cobra.Command {
	o := DeleteSessionUlClOptions{}

	var deleteLbCmd = &cobra.Command{
		Use:     "sessionulcl <UserID> --ulclArgs=<UlCLIP>,...",
		Short:   "Delete a Ulcl configuration in the LoxiLB.",
		Long:    `Delete a Ulcl configuration in the LoxiLB.`,
		Aliases: []string{"ulcl", "sessionulcls", "ulcls"},
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := DeleteSessionUlClValidation(args); err != nil {
				fmt.Println("not valid <UserID>")
				return
			}

			o.UserID = args[0]
			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(ctx, time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}
			for _, ulclIP := range o.UlClArgs {
				subResources := []string{
					"ident", o.UserID,
					"ulclAddress", ulclIP,
				}
				resp, err := client.SessionUlCL().SubResources(subResources).Delete(ctx)
				if err != nil {
					fmt.Printf("Error: Failed to delete Session(UserID: %s)", o.UserID)
					return
				}
				defer resp.Body.Close()
				fmt.Printf("Debug: response.StatusCode: %d\n", resp.StatusCode)
				if resp.StatusCode == http.StatusOK {
					PrintDeleteResult(resp, *restOptions)
					return
				}

			}

		},
	}
	deleteLbCmd.Flags().StringSliceVar(&o.UlClArgs, "ulclArgs", o.UlClArgs, "UlCl IP address can be specified as '<UlClIP>'. It don't need qfi.")

	return deleteLbCmd
}
