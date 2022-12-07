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
package cmd

import (
	"os"

	"loxicmd/cmd/create"
	"loxicmd/cmd/delete"
	"loxicmd/cmd/dump"
	"loxicmd/cmd/get"

	"loxicmd/pkg/api"

	"github.com/spf13/cobra"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "loxicmd",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
	}

	restOptions := &api.RESTOptions{}
	saveOptions := &dump.SaveOptions{}
	applyOptions := &dump.ApplyOptions{}

	rootCmd.PersistentFlags().Int16VarP(&restOptions.Timeout, "timeout", "t", 5, "Set timeout")
	rootCmd.PersistentFlags().StringVarP(&restOptions.Protocol, "protocol", "", "http", "Set API server http/https")
	rootCmd.PersistentFlags().StringVarP(&restOptions.PrintOption, "output", "o", "", "Set output layer (ex.) wide, json)")
	rootCmd.PersistentFlags().StringVarP(&restOptions.ServerIP, "apiserver", "s", "127.0.0.1", "Set API server IP address")
	rootCmd.PersistentFlags().Int16VarP(&restOptions.ServerPort, "port", "p", 11111, "Set API server port number")

	rootCmd.AddCommand(get.GetCmd(restOptions))
	rootCmd.AddCommand(create.CreateCmd(restOptions))
	rootCmd.AddCommand(delete.DeleteCmd(restOptions))

	saveCmd := dump.SaveCmd(saveOptions, restOptions)
	applyCmd := dump.ApplyCmd(applyOptions, restOptions)

	saveCmd.Flags().BoolVarP(&saveOptions.SaveAllConfig, "all", "a", false, "Saves all loxilb configuration")
	saveCmd.Flags().BoolVarP(&saveOptions.SaveIpConfig, "ip", "i", false, "Saves IP configuration")
	saveCmd.Flags().BoolVarP(&saveOptions.SaveLBConfig, "lb", "l", false, "Saves Load Balancer rules configuration")
	saveCmd.Flags().BoolVarP(&saveOptions.SaveSessionConfig, "session", "", false, "Saves session configuration")
	saveCmd.Flags().BoolVarP(&saveOptions.SaveUlClConfig, "ulcl", "", false, "Saves ulcl configuration")
	saveCmd.Flags().BoolVarP(&saveOptions.SaveFWConfig, "firewall", "", false, "Saves firewall configuration")

	saveCmd.MarkFlagsMutuallyExclusive("all", "ip", "lb", "session", "ulcl", "firewall")

	applyCmd.Flags().StringVarP(&applyOptions.IpConfigFile, "ip", "i", "", "IP config file to apply")
	applyCmd.Flags().StringVarP(&applyOptions.Intf, "per-intf", "", "", "Apply configuration only for specific interface")
	applyCmd.Flags().BoolVarP(&applyOptions.Route, "ipv4route", "r", false, "Apply route configuration only for specific interface")
	applyCmd.Flags().StringVarP(&applyOptions.ConfigPath, "config-path", "c", "/opt/loxilb/ipconfig/", "Configuration path only for applying per interface config")
	applyCmd.Flags().StringVarP(&applyOptions.LBConfigFile, "lb", "l", "", "Load Balancer config file to apply")
	applyCmd.Flags().StringVarP(&applyOptions.SessionConfigFile, "session", "", "", "Session config file to apply")
	applyCmd.Flags().StringVarP(&applyOptions.SessionUlClConfigFile, "ulcl", "", "", "Ulcl config file to apply")

	rootCmd.AddCommand(saveCmd)
	rootCmd.AddCommand(applyCmd)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
