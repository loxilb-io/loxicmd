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
package set

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

// NewLogLevelCmd represents the save command
func NewSetBFDCmd(restOptions *api.RESTOptions) *cobra.Command {
	o := api.BFDSessionInfo{}
	SetBFDCmd := &cobra.Command{
		Use:   "bfd remoteIP [--instance=<instance>] [--interval=<interval>] [--retryCount=<count>]",
		Short: "bfd session configuration",
		Long: `bfd session congfigration
--instance   - Cluster Instance name
--interval   - BFD packet Tx interval value in microseconds
--retryCount - Maximum number of retry to detect failure`,

		Aliases: []string{"bfd-session"},
		Run: func(cmd *cobra.Command, args []string) {

			// Make bfdMod
			if err := ReadSetBfdOptions(&o, args); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			resp, err := SetBFDAPICall(restOptions, o)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				PrintSetResult(resp, *restOptions)
				return
			}

		},
	}
	SetBFDCmd.Flags().StringVarP(&o.Instance, "instance", "", "default", "Specify the cluster instance name")
	SetBFDCmd.Flags().Uint64VarP(&o.Interval, "interval", "", 0, "Specify the BFD packet tx interval (in microseconds)")
	SetBFDCmd.Flags().Uint8VarP(&o.RetryCount, "retryCount", "", 0, "Specify the number of retries")

	return SetBFDCmd
}

func ReadSetBfdOptions(o *api.BFDSessionInfo, args []string) error {
	if len(args) > 1 {
		return errors.New("set bfd command get so many args")
	} else if len(args) < 1 {
		return errors.New("set bfd need <RemoteIP> args")
	}

	if val := net.ParseIP(args[0]); val != nil {
		o.RemoteIP = args[0]
	} else {
		return fmt.Errorf("remote IP '%s' is invalid format", args[0])
	}

	return nil
}

func SetBFDAPICall(restOptions *api.RESTOptions, bfdModel api.BFDSessionInfo) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}

	return client.BFDSession().Create(ctx, bfdModel)
}
