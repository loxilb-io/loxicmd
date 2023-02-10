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

type CreateSessionUlClOptions struct {
	UserID   string
	UlClArgs []string
}

func NewCreateSessionUlClCmd(restOptions *api.RESTOptions) *cobra.Command {
	o := CreateSessionUlClOptions{}

	var createSessionCmd = &cobra.Command{
		Use:   "sessionulcl <userID> --ulclArgs=<QFI>:<ulclIP>,... ",
		Short: "Create a Session UlCl",
		Long: `Create a Session UlCl using LoxiLB

ex) loxicmd create sessionulcl user1 --ulclArgs=16:192.33.125.1
		`,
		Aliases: []string{"ulcl", "sessionulcls", "ulcls"},
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			var SessionMods api.UlclInformationGet
			// Make SessionMod
			if err := ReadCreateSessionUlClOptions(&SessionMods, args, len(o.UlClArgs)); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}

			if err := GetUlClArgsPairList(&SessionMods, o.UlClArgs); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			for _, SessionMod := range SessionMods.UlclInfo {
				resp, err := SessionUlClAPICall(restOptions, SessionMod)
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

			}

		},
	}

	createSessionCmd.Flags().StringSliceVar(&o.UlClArgs, "ulclArgs", o.UlClArgs, "Port pairs can be specified as '<QFI>:<UlClIP>'")

	return createSessionCmd
}

func ReadCreateSessionUlClOptions(o *api.UlclInformationGet, args []string, numUlclArgs int) error {
	if len(args) > 2 {
		return errors.New("create Session ulcl command get so many args")
	} else if len(args) <= 0 {
		return errors.New("create Session ulcl need <userID> arg")
	}

	for i := 0; i < numUlclArgs; i++ {
		t := api.SessionUlClMod{Ident: args[0]}
		o.UlclInfo = append(o.UlclInfo, t)
	}

	return nil
}

func GetUlClArgsPairList(o *api.UlclInformationGet, ulclArgs []string) error {
	for i, ulclArg := range ulclArgs {
		ulclArgsPair := strings.Split(ulclArg, ":")
		if len(ulclArgsPair) != 2 {
			return fmt.Errorf("ulclArgs '%s' is invalid format", ulclArgs)
		}

		// 0 is TeID, 1 is TunnelIP
		qfi, err := strconv.Atoi(ulclArgsPair[0])
		if err != nil {
			return fmt.Errorf("ulclArgs's QFI '%s' is invalid format", ulclArgsPair[0])
		}
		o.UlclInfo[i].Args.Qfi = uint8(qfi)
		if val := net.ParseIP(ulclArgsPair[1]); val != nil {
			o.UlclInfo[i].Args.Addr = val
		} else {
			return fmt.Errorf("ulclArgs's IP '%s' is invalid format", ulclArgsPair[1])
		}

	}
	return nil
}

func SessionUlClAPICall(restOptions *api.RESTOptions, ulclModel api.SessionUlClMod) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}

	return client.SessionUlCL().Create(ctx, ulclModel)
}
