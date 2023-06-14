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
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

type LogOptions struct {
	IpConfigFile string
}

// NewLogLevelCmd represents the save command
func NewSetLogLevelCmd(restOptions *api.RESTOptions) *cobra.Command {
	LogLevelCmd := &cobra.Command{
		Use:     "log-level",
		Short:   "log-level configuration",
		Long:    `log-level congfigration`,
		Aliases: []string{"loglevel"},
		Run: func(cmd *cobra.Command, args []string) {
			var parmaMod api.ParamDump
			// Make paramMod
			if err := ReadSetLogLevelOptions(&parmaMod, args); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			resp, err := ParamAPICall(restOptions, parmaMod)
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
	return LogLevelCmd
}

func ReadSetLogLevelOptions(o *api.ParamDump, args []string) error {
	if len(args) > 1 {
		return errors.New("Set log-level command get so many args")
	} else if len(args) < 1 {
		return errors.New("Set log-level need <log-level> args")
	}

	// Validate Log level option
	if IsValidLogLevel(args[0]) {
		o.LogLevel = args[0]
	} else {
		return errors.New("Set log-level in the debug,info,error,warning,notice,critical,emergency,alert.")
	}

	return nil
}

func IsValidLogLevel(loglevel string) bool {
	switch loglevel {
	case
		"debug",
		"info",
		"error",
		"warning",
		"notice",
		"critical",
		"emergency",
		"alert":
		return true
	}
	return false
}

func ParamAPICall(restOptions *api.RESTOptions, parmaModel api.ParamDump) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}

	return client.Param().Create(ctx, parmaModel)
}
