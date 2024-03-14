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
	"errors"
	"fmt"
	"loxicmd/cmd/create"
	"loxicmd/cmd/set"
	"loxicmd/pkg/api"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type ApplyOptions struct {
	IpConfigFile          string
	LBConfigFile          string
	SessionConfigFile     string
	SessionUlClConfigFile string
	FWConfigFile          string
	NormalConfigFile      string
	BFDConfigFile         string
	Intf                  string
	ConfigPath            string
	Route                 bool
}

// applyCmd represents the save command
func ApplyCmd(options *ApplyOptions, restOptions *api.RESTOptions) *cobra.Command {
	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply configuration",
		Long:  `Reads and apply configuration from the text file`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd
			_ = args
			if len(options.IpConfigFile) == 0 &&
				len(options.LBConfigFile) == 0 &&
				len(options.SessionConfigFile) == 0 &&
				len(options.SessionUlClConfigFile) == 0 &&
				len(options.FWConfigFile) == 0 &&
				len(options.Intf) == 0 &&
				len(options.NormalConfigFile) == 0 &&
				len(options.BFDConfigFile) == 0 {
				fmt.Println("Provide valid options")
				cmd.Help()
				return
			}
			if len(options.IpConfigFile) > 0 {
				ApplyIpConfig(options.IpConfigFile)
				fmt.Printf("Configuration applied - %s\n", options.IpConfigFile)
			}

			if options.Route && len(options.Intf) > 0 {
				addRoute(options.ConfigPath, options.Intf)
				fmt.Printf("Route Configuration applied for - %s\n", options.Intf)
				return
			}
			if len(options.Intf) > 0 {
				ApplyIpConfigPerInterface(options.ConfigPath, options.Intf)
				fmt.Printf("Configuration applied for - %s\n", options.Intf)
			}

			if len(options.LBConfigFile) > 0 {
				ApplyLbConfig(options.LBConfigFile, restOptions)
				fmt.Printf("Configuration applied - %s\n", options.LBConfigFile)
			}
			if len(options.SessionConfigFile) > 0 {
				ApplySessionConfig(options.SessionConfigFile, restOptions)
				fmt.Printf("Configuration applied - %s\n", options.SessionConfigFile)
			}
			if len(options.SessionUlClConfigFile) > 0 {
				ApplySessionUlClConfig(options.SessionUlClConfigFile, restOptions)
				fmt.Printf("Configuration applied - %s\n", options.SessionUlClConfigFile)
			}
			if len(options.FWConfigFile) > 0 {
				ApplyFWConfig(options.FWConfigFile, restOptions)
				fmt.Printf("Configuration applied - %s\n", options.FWConfigFile)
			}
			if len(options.BFDConfigFile) > 0 {
				ApplyBFDConfig(options.BFDConfigFile, restOptions)
				fmt.Printf("Configuration applied - %s\n", options.BFDConfigFile)
			}
			if len(options.NormalConfigFile) > 0 {
				if err := ApplyFileConfig(options.NormalConfigFile, restOptions); err != nil {
					fmt.Printf("Configuration failed - %s\n", options.NormalConfigFile)
				} else {
					fmt.Printf("Configuration applied - %s\n", options.NormalConfigFile)
				}
			}

		},
	}
	// -f filename option
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

func bashCommand(command string) {
	fmt.Println(command)
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%v\n", string(output))
}

func getType(path string, intf string) (string, error) {
	file := path + "/" + intf + "/type"
	var text string
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
		text = scanner.Text()
		break
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return text, nil
}

func getBondMode(path string, intf string) (string, error) {
	file := path + "/" + intf + "/mode"
	var text string
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
		text = scanner.Text()
		break
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return text, nil
}

func addMTU(path string, intf string) {
	// open file
	file := path + "/" + intf + "/mtu"
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		return
	}
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
		//fmt.Printf("%s\n", scanner.Text())
		command := "ip link set dev " + intf + " mtu " + scanner.Text()
		bashCommand(command)
		break
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func addIpAddr(path string, intf string) {
	// open file
	file := path + "/" + intf + "/ipv4addr"
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		return
	}

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
		//fmt.Printf("%s\n", scanner.Text())
		command := "ip addr add " + scanner.Text() + " dev " + intf
		bashCommand(command)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func addL2FDBs(path string, intf string) {
	// open file
	file := path + "/" + intf + "/l2fdbs"
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		return
	}

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
		//fmt.Printf("%s\n", scanner.Text())
		command := "bridge fdb add " + scanner.Text() + " dev " + intf + " permanent"
		bashCommand(command)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func addVxFDBs(path string, intf string) {
	// open file
	file := path + "/" + intf + "/vxfdbs"
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		return
	}
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
		//fmt.Printf("%s\n", scanner.Text())
		command := "bridge fdb append " + scanner.Text() + " dev " + intf + " permanent"
		bashCommand(command)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func addNeigh(path string, intf string) {
	// open file
	file := path + "/" + intf + "/ipv4neigh"
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		return
	}
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
		//fmt.Printf("%s\n", scanner.Text())
		command := "ip neigh add " + scanner.Text() + " dev " + intf + " permanent"
		bashCommand(command)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func addRoute(path string, intf string) {
	// open file
	file := path + "/" + intf + "/ipv4route"
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		return
	}
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
		//fmt.Printf("%s\n", scanner.Text())
		command := "ip route replace " + scanner.Text() + " proto static"
		bashCommand(command)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func addSubIntf(path string, intf string) {
	// open file
	file := path + "/" + intf + "/subintf"
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		return
	}
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
		//fmt.Printf("%s\n", scanner.Text())
		token := strings.Split(scanner.Text(), "|")

		command := "ip link add " + token[0] + " link " + token[1] + " type vlan id " + token[2]
		bashCommand(command)

		command = "ip link set " + token[0] + " up"
		bashCommand(command)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func addBridge(intf string) {

	command := "ip link add " + intf + " type bridge "
	bashCommand(command)
	command = "ip link set " + intf + " up"
	bashCommand(command)
}

func addBond(path string, intf string) {

	command := "ip link add " + intf + " type bond "
	bashCommand(command)
	mode, err := getBondMode(path, intf)
	if err == nil {
		command = "ip link set " + intf + " type bond mode " + mode
		bashCommand(command)
	}
	command = "ip link set " + intf + " up"
	bashCommand(command)
}

func addVxlan(path string, intf string) {
	// open file
	file := path + "/" + intf + "/info"
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		return
	}
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
		//fmt.Printf("%s\n", scanner.Text())
		token := strings.Split(scanner.Text(), "|")

		command := "ip link add " + intf + " type vxlan id " + token[0] +
			" local " + token[1] + " dev " + token[2] + " dst 4789"
		bashCommand(command)

		command = "ip link set " + intf + " up"
		bashCommand(command)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func setMaster(intf string, master string) {

	command := "ip link set " + intf + " down"
	bashCommand(command)
	command = "ip link set " + intf + " master " + master
	bashCommand(command)
	command = "ip link set " + intf + " up"
	bashCommand(command)
}

func setDev(intf string) {

	command := "ip link set " + intf + " up"
	bashCommand(command)
}

func addMaster(path string, intf string) {
	// open file
	file := path + "/" + intf + "/master"
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		return
	}
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
		//fmt.Printf("%s\n", scanner.Text())
		token := strings.Split(scanner.Text(), "|")

		if token[1] == "bridge" {
			addBridge(token[0])
			setMaster(intf, token[0])
		} else if token[1] == "vxlan" {
			addVxlan(path, token[0])
		} else if token[1] == "bond" {
			addBond(path, token[0])
			setMaster(intf, token[0])
		}
		break
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func ApplyIpConfigPerInterface(path string, intf string) {

	itype, err := getType(path, intf)
	if err == nil {
		if itype == "phy" || itype == "bond" {
			setDev(intf)
			addMTU(path, intf)
			addIpAddr(path, intf)
			addL2FDBs(path, intf)
			addNeigh(path, intf)
			addRoute(path, intf)
			addSubIntf(path, intf)
			addMaster(path, intf)
		} else if itype == "subintf" {
			setDev(intf)
			addIpAddr(path, intf)
			addNeigh(path, intf)
			addRoute(path, intf)
			addSubIntf(path, intf)
			addMaster(path, intf)
		} else if itype == "bridge" {
			setDev(intf)
			addIpAddr(path, intf)
			addNeigh(path, intf)
			addRoute(path, intf)
			addMaster(path, intf)
		} else if itype == "vxlan" {
			setDev(intf)
			addIpAddr(path, intf)
			addVxFDBs(path, intf)
			addNeigh(path, intf)
			addRoute(path, intf)
			addMaster(path, intf)
		}
	} else {
		fmt.Printf("Unable to get type of intf(%v)", err)
		return
	}
}

func ApplyLbConfig(file string, restOptions *api.RESTOptions) {
	// open file
	var lbresp api.LbRuleModGet
	byteBuf, err := os.ReadFile(file)
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
			Sel:        lb.Service.Sel,
			Mode:       lb.Service.Mode,
			BGP:        lb.Service.BGP,
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

func ApplySessionConfig(file string, restOptions *api.RESTOptions) {
	// open file
	var resp api.SessionInformationGet
	byteBuf, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Unmashal to Json
	if err := json.Unmarshal(byteBuf, &resp); err != nil {
		fmt.Printf("Error: Failed to unmarshal File: (%s)\n", err.Error())
		return
	}

	// POST the dump
	for _, sess := range resp.SessionInfo {
		resp, err := create.SessionAPICall(restOptions, sess)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
		defer resp.Body.Close()
	}
}

func ApplySessionUlClConfig(file string, restOptions *api.RESTOptions) {
	// open file
	var resp api.UlclInformationGet
	byteBuf, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Unmashal to Json
	if err := json.Unmarshal(byteBuf, &resp); err != nil {
		fmt.Printf("Error: Failed to unmarshal File: (%s)\n", err.Error())
		return
	}

	// POST the dump
	for _, ulcl := range resp.UlclInfo {
		fmt.Printf("ulcl: %v\n", ulcl)
		resp, err := create.SessionUlClAPICall(restOptions, ulcl)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
		defer resp.Body.Close()
	}
}

func ApplyFWConfig(file string, restOptions *api.RESTOptions) {
	// open file
	var resp api.FWInformationGet
	byteBuf, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Unmashal to Json
	if err := json.Unmarshal(byteBuf, &resp); err != nil {
		fmt.Printf("Error: Failed to unmarshal File: (%s)\n", err.Error())
		return
	}

	// POST the dump
	for _, fw := range resp.FWInfo {
		fmt.Printf("fw: %v\n", fw)
		resp, err := create.FirewallAPICall(restOptions, fw)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
		defer resp.Body.Close()
	}
}

func ApplyBFDConfig(file string, restOptions *api.RESTOptions) {
	// open file
	var resp api.BFDSessionGet
	byteBuf, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Unmashal to Json
	if err := json.Unmarshal(byteBuf, &resp); err != nil {
		fmt.Printf("Error: Failed to unmarshal File: (%s)\n", err.Error())
		return
	}

	// POST the dump
	for _, b := range resp.BFDSessionAttr {
		fmt.Printf("bfd: %v\n", b)
		resp, err := set.SetBFDAPICall(restOptions, b)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
		defer resp.Body.Close()
	}
}
