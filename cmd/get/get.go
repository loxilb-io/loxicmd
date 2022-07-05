/*
Copyright © 2022 baekgyun jung <backguyn@netlox.io>
Author: Inho gog <inhogog2@netlox.io>
*/
package get

import (
	"fmt"

	"loxicmd/pkg/api"

	"github.com/spf13/cobra"
)

func GetCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetCmd = &cobra.Command{
		Use:   "get",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly Get a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Get called")
		},
	}

	GetCmd.AddCommand(NewGetLoadBalancerCmd(restOptions))
	GetCmd.AddCommand(NewGetConntrackCmd(restOptions))
	return GetCmd
}
