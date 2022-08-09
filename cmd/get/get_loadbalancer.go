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
	"io/ioutil"
	"loxicmd/pkg/api"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func NewGetLoadBalancerCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetLbCmd = &cobra.Command{
		Use:     "loadbalancer",
		Short:   "Get a LoadBalancer",
		Aliases: []string{"lb", "loadbalancers", "lbs"},
		Long:    `It shows Load balancer Information`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd
			_ = args
			resp, err := LoadbalancerAPICall(restOptions)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				PrintGetLbResult(resp, *restOptions)
				return
			}

		},
	}

	return GetLbCmd
}

func PrintGetLbResult(resp *http.Response, o api.RESTOptions) {
	lbresp := api.LbRuleModGet{}
	var data [][]string
	resultByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}

	if err := json.Unmarshal(resultByte, &lbresp); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	// if json options enable, it print as a json format.
	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(lbresp, "", "    ")
		fmt.Println(string(resultIndent))
		return
	}

	// Table Init
	table := TableInit()

	// Making load balance data
	for _, lbrule := range lbresp.LbRules {
		if o.PrintOption == "wide" {
			table.SetHeader(LOADBALANCER_WIDE_TITLE)
			for i, eps := range lbrule.Endpoints {
				if i == 0 {
					data = append(data, []string{lbrule.Service.ExternalIP, fmt.Sprintf("%d", lbrule.Service.Port), lbrule.Service.Protocol, fmt.Sprintf("%d", lbrule.Service.Sel),
						eps.EndpointIP, fmt.Sprintf("%d", eps.TargetPort), fmt.Sprintf("%d", eps.Weight)})
				} else {
					data = append(data, []string{"", "", "", "", eps.EndpointIP, fmt.Sprintf("%d", eps.TargetPort), fmt.Sprintf("%d", eps.Weight)})
				}
			}
		} else {
			table.SetHeader(LOADBALANCER_TITLE)
			data = append(data, []string{lbrule.Service.ExternalIP, fmt.Sprintf("%d", lbrule.Service.Port), lbrule.Service.Protocol, fmt.Sprintf("%d", lbrule.Service.Sel), fmt.Sprintf("%d", len(lbrule.Endpoints))})
		}
	}

	// Rendering the load balance data to table
	TableShow(data, table)
}

func LoadbalancerAPICall(restOptions *api.RESTOptions) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}
	resp, err := client.LoadBalancer().SetUrl("/config/loadbalancer/all").Get(ctx)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return nil, err
	}
	return resp, nil
}

func Lbdump(restOptions *api.RESTOptions) (string, error) {
	// File Open
	fileP := []string{"lbconfig_", ".txt"}
	t := time.Now()
	file := strings.Join(fileP, t.Local().Format("2006-01-02_15:04:05"))
	f, err := os.Create(file)
	if err != nil {
		fmt.Printf("Can't create dump file\n")
		os.Exit(1)
	}
	defer f.Close()

	// API Call
	resp, err := LoadbalancerAPICall(restOptions)
	if err != nil {
		return "", err
	}
	resultByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
	}
	// Write
	f.Write(resultByte)

	if _, err := os.Stat("/opt/loxilb/lbconfig.txt"); errors.Is(err, os.ErrNotExist) {
		if err != nil {
			fmt.Println("There is no saved config file")
		}
	} else {
		command := "mv /opt/loxilb/lbconfig.txt /opt/loxilb/lbconfig.txt.bk"
		cmd := exec.Command("bash", "-c", command)
		_, err := cmd.Output()
		if err != nil {
			fmt.Println("Can't backup /opt/loxilb/lbconfig.txt")
			return file, err
		}
	}
	command := "cp -R " + file + " /opt/loxilb/lbconfig.txt"
	cmd := exec.Command("bash", "-c", command)
	fmt.Println(cmd)
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	return file, nil
}
