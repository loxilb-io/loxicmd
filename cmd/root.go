/*
Copyright © 2022 Netlox Inc <backguyn@netlox.io>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"

	"loxicmd/cmd/create"
	"loxicmd/cmd/delete"
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

	rootCmd.PersistentFlags().Int16VarP(&restOptions.Timeout, "timeout", "t", 5, "Set timeout")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Set output layer")
	rootCmd.PersistentFlags().StringVarP(&restOptions.ServerIP, "apiserver", "s", "127.0.0.1", "Set API server IP address")
	rootCmd.PersistentFlags().Int16VarP(&restOptions.ServerPort, "port", "p", 8181, "Set API server port number")

	rootCmd.AddCommand(create.CreateCmd(restOptions))
	rootCmd.AddCommand(delete.DeleteCmd(restOptions))

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
