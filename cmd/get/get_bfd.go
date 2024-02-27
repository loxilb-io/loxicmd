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
	"strings"
	"os"
	"os/exec"
	"errors"

	"github.com/spf13/cobra"
)

func NewGetBFDCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetBFDCmd = &cobra.Command{
		Use:     "bfd",
		Short:   "Get all BFD sessions",
		Long:    `It shows BFD Sessions in the LoxiLB`,

		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}
			resp, err := client.Status().SetUrl("config/bfd/all").Get(ctx)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				PrintGetBFDResult(resp, *restOptions)
				return
			}

		},
	}

	return GetBFDCmd
}

func PrintGetBFDResult(resp *http.Response, o api.RESTOptions) {
	BFDresp := api.BFDSessionGet{}
	var data [][]string
	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}

	if err := json.Unmarshal(resultByte, &BFDresp); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	// if json options enable, it print as a json format.
	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(BFDresp, "", "    ")
		fmt.Println(string(resultIndent))
		return
	}

	// Table Init
	table := TableInit()
	
	// Making data
	for _, bfd := range BFDresp.BFDSessionAttr {
		if (o.PrintOption == "wide") {
			table.SetHeader(BFD_WIDE_TITLE)
			data = append(data, []string{bfd.Instance, bfd.RemoteIP, bfd.SourceIP,
			 fmt.Sprintf("%d",bfd.Port), fmt.Sprintf("%d us",bfd.Interval), fmt.Sprintf("%d",bfd.RetryCount), bfd.State})
		} else {
			table.SetHeader(BFD_TITLE)
			data = append(data, []string{bfd.Instance, bfd.RemoteIP, bfd.State})
		}
	}
	// Rendering the data to table
	TableShow(data, table)
}

func BFDdump(restOptions *api.RESTOptions, path string) (string, error) {
	BFDresp := api.BFDSessionGet{}

	// File Open
	fileP := []string{"BFDconfig_", ".txt"}
	t := time.Now()
	file := strings.Join(fileP, t.Local().Format("2006-01-02_15:04:05"))
	f, err := os.Create(file)
	if err != nil {
		fmt.Printf("Can't create dump file\n")
		os.Exit(1)
	}
	defer f.Close()

	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}
	resp, err := client.Status().SetUrl("config/bfd/all").Get(ctx)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return "", err
	}
	if resp.StatusCode == http.StatusOK {

		resultByte, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
			return "", err
		}

		if err := json.Unmarshal(resultByte, &BFDresp); err != nil {
			fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
			return "", err
		}
		
		bfds := api.BFDSessionGet{}
		bfds.BFDSessionAttr = BFDresp.BFDSessionAttr
		/*
		for _, b := range BFDresp.BFDSessionAttr {
			bfds.BFDSessionAttr = append(bfds.BFDSessionAttr, b)

			data = append(data, []string{bfd.Instance, bfd.RemoteIP, bfd.SourceIP,
				fmt.Sprintf("%d",bfd.Port), fmt.Sprintf("%d us",bfd.Interval), fmt.Sprintf("%d",bfd.RetryCount), bfd.State})
		} */
		cfgResultByte, err := json.Marshal(bfds)
		if err != nil {
			fmt.Printf("Error: Failed to marshal BFD Cfg: (%s)\n", err.Error())
			return "", err
		}

		// Write
		f.Write(cfgResultByte)
		cfile := path + "BFDconfig.txt"
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
			fmt.Println("Failed copy file to", cfile)
			return file, err
		}
		return file, nil
	}
	return "", err
}