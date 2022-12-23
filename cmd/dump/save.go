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
package dump

import (
	"fmt"
	"os"
	"errors"
	get "loxicmd/cmd/get"
	"loxicmd/pkg/api"

	"github.com/spf13/cobra"
)

type SaveOptions struct {
	SaveIpConfig      bool
	SaveLBConfig      bool
	SaveSessionConfig bool
	SaveUlClConfig    bool
	SaveFWConfig      bool
	SaveAllConfig     bool
}

// saveCmd represents the save command
func SaveCmd(saveOpts *SaveOptions, restOptions *api.RESTOptions) *cobra.Command {
	saveCmd := &cobra.Command{
		Use:   "save",
		Short: "saves current configuration",
		Long:  `saves current configuration in text file`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd
			_ = args
			dpath := "/etc/loxilb/"
			if _, err := os.Stat(dpath); errors.Is(err, os.ErrNotExist) {
				err := os.Mkdir(dpath, os.ModePerm)
				if err != nil {
					fmt.Println("Can't create config dir /etc/loxilb/")
					return
				}
			}
			if saveOpts.SaveIpConfig || saveOpts.SaveAllConfig {
				file := get.Nlpdump(dpath)
				fmt.Println("IP Configuration saved in", file)
			}
			if saveOpts.SaveLBConfig || saveOpts.SaveAllConfig {
				lbfile, err := get.Lbdump(restOptions, dpath)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println("LB Configuration saved in", lbfile)
			}
			if saveOpts.SaveSessionConfig || saveOpts.SaveAllConfig {
				sessionFile, err := get.Sessiondump(restOptions, dpath)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println("Session Configuration saved in", sessionFile)
			}
			if saveOpts.SaveUlClConfig || saveOpts.SaveAllConfig {
				ulclFile, err := get.SessionUlCldump(restOptions, dpath)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println("UlCl Configuration saved in", ulclFile)
			}
			if saveOpts.SaveFWConfig || saveOpts.SaveAllConfig {
				FWFile, err := get.FWdump(restOptions, dpath)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println("Firewall Configuration saved in", FWFile)
			}
		},
	}
	return saveCmd
}
