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
package delete

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"loxicmd/pkg/api"

	"github.com/spf13/cobra"
)

type DeleteFirewallOptions struct {
	FirewallRule []string
}

func NewDeleteFirewallCmd(restOptions *api.RESTOptions) *cobra.Command {
	o := DeleteFirewallOptions{}
	var deleteFirewallCmd = &cobra.Command{
		Use:   "firewall --firewallRule=<ruleKey>:<ruleValue>",
		Short: "Delete a Firewall",
		Long: `Delete a Firewall using ruleKey in the LoxiLB.

--<ruleKey>s of firewallRule
sourceIP(string) - Source IP in CIDR notation
destinationIP(string) - Destination IP in CIDR notation	
minSourcePort(int) - Minimum source port range	
maxSourcePort(int) - Maximum source port range	
minDestinationPort(int) - Minimum destination port range	
maxDestinationPort(int) - Maximum source port range	
protocol(int) - the protocol	
portName(string) - the incoming port	
preference(int) - User preference for ordering	

ex) loxicmd delete firewall --firewallRule="sourceIP:1.2.3.2/32,destinationIP:2.3.1.2/32,preference:200"
		`,
		Aliases: []string{"Firewall", "fw", "firewalls"},
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(o.FirewallRule) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			client := api.NewLoxiClient(restOptions)
			ctx := context.TODO()
			var cancel context.CancelFunc
			if restOptions.Timeout > 0 {
				ctx, cancel = context.WithTimeout(ctx, time.Duration(restOptions.Timeout)*time.Second)
				defer cancel()
			}

			qeury, err := MakefirewallDeleteQurey(o.FirewallRule)
			if err != nil {
				fmt.Printf("Error: Failed to delete Firewall")
				return
			}
			resp, err := client.Firewall().Query(qeury).Delete(ctx)
			if err != nil {
				fmt.Printf("Error: Failed to delete Firewall")
				return
			}
			defer resp.Body.Close()
			fmt.Printf("Debug: response.StatusCode: %d\n", resp.StatusCode)
			if resp.StatusCode == http.StatusOK {
				PrintDeleteResult(resp, *restOptions)
				return
			}

		},
	}

	deleteFirewallCmd.Flags().StringSliceVar(&o.FirewallRule, "firewallRule", o.FirewallRule, "Information related to firewall rule")
	return deleteFirewallCmd
}

func MakefirewallDeleteQurey(FirewallRule []string) (map[string]string, error) {
	query := map[string]string{}
	for _, v := range FirewallRule {
		firewallArgsPair := strings.Split(v, ":")
		if len(firewallArgsPair) != 2 {
			return nil, fmt.Errorf("Error: Failed to delete Firewall")
		}
		query[firewallArgsPair[0]] = firewallArgsPair[1]
	}
	return query, nil
}
