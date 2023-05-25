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

func NewGetEndPointCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetfwCmd = &cobra.Command{
		Use:     "endpoint",
		Short:   "Get endpoints",
		Aliases: []string{"endpoint", "ep", "endpoints"},
		Long:    `It shows End Point Information`,
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
			resp, err := client.Firewall().SetUrl("/config/endpoint/all").Get(ctx)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				PrintGetEPResult(resp, *restOptions)
				return
			}
		},
	}

	return GetfwCmd
}

func PrintGetEPResult(resp *http.Response, o api.RESTOptions) {
	epResp := api.EPInformationGet{}
	var data [][]string
	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}

	if err := json.Unmarshal(resultByte, &epResp); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	// if json options enable, it print as a json format.
	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(epResp, "", "    ")
		fmt.Println(string(resultIndent))
		return
	}

	epResp.Sort()

	// Table Init
	table := TableInit()

	// Making load balance data
	for _, ep := range epResp.EPInfo {
		table.SetHeader(ENDPOINT_TITLE)
		data = append(data, []string{ep.HostName, ep.Name, fmt.Sprintf("%s:%s", ep.ProbeType, ep.ProbeReq), fmt.Sprintf("%d", ep.ProbePort),
			fmt.Sprintf("%d", ep.ProbeDuration), fmt.Sprintf("%d", ep.InActTries),
			ep.MinDelay, ep.AvgDelay, ep.MaxDelay, ep.CurrState})

	}

	// Rendering the load balance data to table
	TableShow(data, table)
}

func EPAPICall(restOptions *api.RESTOptions) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}
	resp, err := client.Firewall().SetUrl("/config/endpoint/all").Get(ctx)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return nil, err
	}
	return resp, nil
}

func EPdump(restOptions *api.RESTOptions, path string) (string, error) {
	epResp := api.EPInformationGet{}
	// File Open
	fileP := []string{"EPconfig_", ".txt"}
	t := time.Now()
	file := strings.Join(fileP, t.Local().Format("2006-01-02_15:04:05"))
	f, err := os.Create(file)
	if err != nil {
		fmt.Printf("Can't create dump file\n")
		os.Exit(1)
	}
	defer f.Close()

	// API Call
	resp, err := EPAPICall(restOptions)
	if err != nil {
		return "", err
	}
	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
	}

	if err := json.Unmarshal(resultByte, &epResp); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return "", err
	}

	epMs := api.EPConfig{}
	for _, ep := range epResp.EPInfo {
		var epm api.EndPointMod
		epm.HostName = ep.HostName
		epm.Name = ep.Name
		epm.InActTries = ep.InActTries
		epm.ProbeType = ep.ProbeType
		epm.ProbePort = ep.ProbePort
		epm.ProbeReq = ep.ProbeReq
		epm.ProbeResp = ep.ProbeResp
		epm.ProbeDuration = ep.ProbeDuration
		epMs.EPInfo = append(epMs.EPInfo, epm)
	}

	cfgResultByte, err := json.Marshal(epMs)
	if err != nil {
		fmt.Printf("Error: Failed to marshal EP Cfg: (%s)\n", err.Error())
		return "", err
	}

	// Write
	f.Write(cfgResultByte)

	cfile := path + "EPconfig.txt"
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
