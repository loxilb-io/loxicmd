// Author: Inho gog <inhogog2@netlox.io>
package get

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"loxicmd/pkg/api"
	"net/http"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func NewGetLoadBalancerCmd(restOptions *api.RESTOptions) *cobra.Command {
	var GetLbCmd = &cobra.Command{
		Use:     "loadbalancer",
		Short:   "Get a LoadBalancer",
		Aliases: []string{"lb", "loadbalancers", "lbs"},
		Long:    `It shows Load balancer Information`,
		Run: func(cmd *cobra.Command, args []string) {
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
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	// Making load balance data
	for _, lbrule := range lbresp.LbRules {
		if o.PrintOption == "wide" {
			table.SetHeader([]string{"ExternalIP", "Port", "Protocol", "Select", "EndpointIP", "TargetPort", "Weight"})
			for i, eps := range lbrule.Endpoints {
				if i == 0 {
					data = append(data, []string{lbrule.Service.ExternalIP, fmt.Sprintf("%d", lbrule.Service.Port), lbrule.Service.Protocol, fmt.Sprintf("%d", lbrule.Service.Sel),
						eps.EndpointIP, fmt.Sprintf("%d", eps.TargetPort), fmt.Sprintf("%d", eps.Weight)})
				} else {
					data = append(data, []string{"", "", "", "", eps.EndpointIP, fmt.Sprintf("%d", eps.TargetPort), fmt.Sprintf("%d", eps.Weight)})
				}
			}
		} else {
			table.SetHeader([]string{"ExternalIP", "Port", "Protocol", "Select", "# of Endpoints"})
			data = append(data, []string{lbrule.Service.ExternalIP, fmt.Sprintf("%d", lbrule.Service.Port), lbrule.Service.Protocol, fmt.Sprintf("%d", lbrule.Service.Sel), fmt.Sprintf("%d", len(lbrule.Endpoints))})
		}
	}

	// Rendering the load balance data to table
	table.AppendBulk(data)
	table.Render()
}
