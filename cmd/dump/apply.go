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
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"loxicmd/cmd/create"
	"loxicmd/pkg/api"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

type ConfigFiles struct {
	IpConfigFile string
	LBConfigFile string
}

// applyCmd represents the save command
func ApplyCmd(cfgFiles *ConfigFiles, restOptions *api.RESTOptions) *cobra.Command {
	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply configuration",
		Long:  `Reads and apply configuration from the text file`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd
			_ = args
			if len(cfgFiles.IpConfigFile) == 0 && len(cfgFiles.LBConfigFile) == 0 {
				fmt.Println("Provide valid filename")
				return
			}
			if len(cfgFiles.IpConfigFile) > 0 {
				ApplyIpConfig(cfgFiles.IpConfigFile)
				fmt.Printf("Configuration applied - %s\n", cfgFiles.IpConfigFile)
			}
			if len(cfgFiles.LBConfigFile) > 0 {
				ApplyLbConfig(cfgFiles.LBConfigFile, restOptions)
				fmt.Printf("Configuration applied - %s\n", cfgFiles.LBConfigFile)
			}
		},
	}
	return applyCmd
}

func ApplyIpConfig(file string) {
	// open file
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// do something with a line
		fmt.Printf("%s\n", scanner.Text())
		cmd := exec.Command("bash", "-c", scanner.Text())
		output, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("%v\n", string(output))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func ApplyLbConfig(file string, restOptions *api.RESTOptions) {
	// open file
	var lbresp api.LbRuleModGet
	byteBuf, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Unmashal to Json
	if err := json.Unmarshal(byteBuf, &lbresp); err != nil {
		fmt.Printf("Error: Failed to unmarshal File: (%s)\n", err.Error())
		return
	}

	// POST the dump
	for _, lb := range lbresp.LbRules {
		lbModel := api.LoadBalancerModel{}
		lbService := api.LoadBalancerService{
			ExternalIP: lb.Service.ExternalIP,
			Protocol:   lb.Service.Protocol,
			Port:       lb.Service.Port,
		}

		lbModel.Service = lbService
		for _, ep := range lb.Endpoints {
			endPoint := api.LoadBalancerEndpoint{
				EndpointIP: ep.EndpointIP,
				TargetPort: ep.TargetPort,
				Weight:     ep.Weight,
			}
			lbModel.Endpoints = append(lbModel.Endpoints, endPoint)
		}

		resp, err := create.LoadbalancerAPICall(restOptions, lbModel)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
		defer resp.Body.Close()
	}
}
