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
	"errors"
	"fmt"
	"io"
	"loxicmd/pkg/api"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func NewGetSessionULCLCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetulclCmd = &cobra.Command{
		Use:     "sessionulcl",
		Short:   "Get a sessionUlcl",
		Aliases: []string{"ulcl", "sessionulcls", "ulcls"},
		Long:    `It shows Session UlCl Information`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd
			_ = args
			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}
			resp, err := client.SessionUlCL().SetUrl("/config/sessionulcl/all").Get(ctx)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				PrintGetSessionULCLResult(resp, *restOptions)
				return
			}

		},
	}

	return GetulclCmd
}

func PrintGetSessionULCLResult(resp *http.Response, o api.RESTOptions) {
	ulclresp := api.UlclInformationGet{}
	var data [][]string
	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}

	if err := json.Unmarshal(resultByte, &ulclresp); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	// if json options enable, it print as a json format.
	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(ulclresp, "", "    ")
		fmt.Println(string(resultIndent))
		return
	}

	ulclresp.Sort()

	// Table Init
	table := TableInit()

	// Making load balance data
	for _, ulcl := range ulclresp.UlclInfo {
		table.SetHeader(ULCL_TITLE)
		if len(data) == 0 {
			data = append(data, []string{ulcl.Ident, ulcl.Args.Addr.String(), fmt.Sprintf("%d", ulcl.Args.Qfi)})
		} else {
			// Check duplicatited User Ident.
			dataFlag := true
			for _, tdata := range data {
				if ulcl.Ident == tdata[0] {
					data = append(data, []string{"", ulcl.Args.Addr.String(), fmt.Sprintf("%d", ulcl.Args.Qfi)})
					dataFlag = false
				}
			}
			if dataFlag {
				data = append(data, []string{ulcl.Ident, ulcl.Args.Addr.String(), fmt.Sprintf("%d", ulcl.Args.Qfi)})
			}
		}

	}
	// Rendering the load balance data to table
	TableShow(data, table)
}

func SessionUlClAPICall(restOptions *api.RESTOptions) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}
	resp, err := client.SessionUlCL().SetUrl("/config/sessionulcl/all").Get(ctx)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return nil, err
	}
	return resp, nil
}

func SessionUlCldump(restOptions *api.RESTOptions, path string) (string, error) {
	// File Open
	fileP := []string{"sessionulclconfig_", ".txt"}
	t := time.Now()
	file := strings.Join(fileP, t.Local().Format("2006-01-02_15:04:05"))
	f, err := os.Create(file)
	if err != nil {
		fmt.Printf("Can't create dump file\n")
		os.Exit(1)
	}
	defer f.Close()

	// API Call
	resp, err := SessionUlClAPICall(restOptions)
	if err != nil {
		return "", err
	}
	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
	}
	// Write
	f.Write(resultByte)

	cfile := path + "sessionulclconfig.txt"
	if _, err := os.Stat(cfile); errors.Is(err, os.ErrNotExist) {
		if err != nil {
			fmt.Println("There is no saved config file")
		}
	} else {
		command := "mv " + cfile + " " + cfile + ".bk"
		cmd := exec.Command("bash", "-c", command)
		_, err := cmd.Output()
		if err != nil {
			fmt.Println("Can't backup ", cfile)
			return file, err
		}
	}
	command := "cp -R " + file + " " + cfile
	cmd := exec.Command("bash", "-c", command)
	fmt.Println(cmd)
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	return file, nil
}
