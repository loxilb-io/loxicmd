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
package create

import (
	"context"
	"fmt"
	"loxicmd/pkg/api"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type CreateFirewallOptions struct {
	FirewallRule []string
	Redirect     []string
	Allow        bool
	Drop         bool
	Trap         bool
	Record       bool
	Mark         int
}

func NewCreateFirewallCmd(restOptions *api.RESTOptions) *cobra.Command {
	o := CreateFirewallOptions{}

	var createFirewallCmd = &cobra.Command{
		Use:   "firewall --firewallRule=<ruleKey>:<ruleValue>, [--allow] [--drop] [--trap] [--record] [--redirect=<PortName>] [--setmark=<FwMark>]",
		Short: "Create a Firewall",
		Long: `Create a Firewall using LoxiLB

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


ex) loxicmd create firewall --firewallRule="sourceIP:1.2.3.2/32,destinationIP:2.3.1.2/32,preference:200" --allow
    loxicmd create firewall --firewallRule="sourceIP:1.2.3.2/32,destinationIP:2.3.1.2/32,preference:200" --allow --record
	loxicmd create firewall --firewallRule="sourceIP:1.2.3.2/32,destinationIP:2.3.1.2/32,preference:200" --allow --setmark=10
    loxicmd create firewall --firewallRule="sourceIP:1.2.3.2/32,destinationIP:2.3.1.2/32,preference:200" --drop
	loxicmd create firewall --firewallRule="sourceIP:1.2.3.2/32,destinationIP:2.3.1.2/32,preference:200" --trap
	loxicmd create firewall --firewallRule="sourceIP:1.2.3.2/32,destinationIP:2.3.1.2/32,preference:200" --redirect=hs1
`,
		Aliases: []string{"Firewall", "fw", "firewalls"},
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(o.FirewallRule) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			var FirewallMods api.FwRuleMod
			// Make FirewallMod
			if err := GetFirewallRulePairList(&FirewallMods, o.FirewallRule); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}

			if err := GetFWOptionPairList(&FirewallMods, o); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			resp, err := FirewallAPICall(restOptions, FirewallMods)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
			defer resp.Body.Close()

			fmt.Printf("Debug: response.StatusCode: %d\n", resp.StatusCode)
			if resp.StatusCode != http.StatusOK {
				PrintCreateLbResult(resp, *restOptions)
				return
			}

		},
	}

	createFirewallCmd.Flags().StringSliceVar(&o.FirewallRule, "firewallRule", o.FirewallRule, "Information related to firewall rule")
	createFirewallCmd.Flags().StringSliceVar(&o.Redirect, "redirect", o.Redirect, "Redirect any matching rule")
	createFirewallCmd.Flags().BoolVarP(&o.Allow, "allow", "", false, "Allow any matching rule")
	createFirewallCmd.Flags().BoolVarP(&o.Drop, "drop", "", false, "Drop any matching rule")
	createFirewallCmd.Flags().BoolVarP(&o.Record, "record", "", false, "Record/Dump any matching rule")
	createFirewallCmd.Flags().BoolVarP(&o.Trap, "trap", "", false, " Trap anything matching rule")
	createFirewallCmd.Flags().IntVarP(&o.Mark, "setmark", "", 0, " Add a fw mark")

	return createFirewallCmd
}

func GetFirewallRulePairList(o *api.FwRuleMod, FWrule []string) error {
	for _, FirewallArg := range FWrule {
		FirewallArgsPair := strings.Split(FirewallArg, ":")
		if len(FirewallArgsPair) != 2 {
			return fmt.Errorf("FirewallArgs '%s' is invalid format", FWrule)
		} else if FirewallArgsPair[0] == "protocol" {
			protocol, err := strconv.Atoi(FirewallArgsPair[1])
			if err != nil {
				return err
			}
			o.Rule.Proto = uint8(protocol)
		} else if FirewallArgsPair[0] == "sourceIP" {
			o.Rule.SrcIP = FirewallArgsPair[1]
		} else if FirewallArgsPair[0] == "destinationIP" {
			o.Rule.DstIP = FirewallArgsPair[1]
		} else if FirewallArgsPair[0] == "portName" {
			o.Rule.InPort = FirewallArgsPair[1]
		} else if FirewallArgsPair[0] == "minSourcePort" {
			minSourcePort, err := strconv.Atoi(FirewallArgsPair[1])
			if err != nil {
				return err
			}
			o.Rule.SrcPortMin = uint16(minSourcePort)
		} else if FirewallArgsPair[0] == "maxSourcePort" {
			maxSourcePort, err := strconv.Atoi(FirewallArgsPair[1])
			if err != nil {
				return err
			}
			o.Rule.SrcPortMax = uint16(maxSourcePort)
		} else if FirewallArgsPair[0] == "minDestinationPort" {
			minDestinationPort, err := strconv.Atoi(FirewallArgsPair[1])
			if err != nil {
				return err
			}
			o.Rule.DstPortMin = uint16(minDestinationPort)
		} else if FirewallArgsPair[0] == "maxDestinationPort" {
			maxDestinationPort, err := strconv.Atoi(FirewallArgsPair[1])
			if err != nil {
				return err
			}
			o.Rule.DstPortMax = uint16(maxDestinationPort)
		} else if FirewallArgsPair[0] == "preference" {
			preference, err := strconv.Atoi(FirewallArgsPair[1])
			if err != nil {
				return err
			}
			o.Rule.Pref = uint16(preference)
		}

	}
	return nil
}

func GetFWOptionPairList(FirewallMods *api.FwRuleMod, o CreateFirewallOptions) error {
	// Option boolean check
	if o.Allow {
		FirewallMods.Opts.Allow = true
	} else if o.Drop {
		FirewallMods.Opts.Drop = true
	} else if o.Trap {
		FirewallMods.Opts.Trap = true
	} else if len(o.Redirect) != 0 {
		FirewallMods.Opts.Rdr = true
		FirewallMods.Opts.RdrPort = o.Redirect[0]
	}
	FirewallMods.Opts.Record = o.Record
	FirewallMods.Opts.Mark = o.Mark

	return nil
}

func FirewallAPICall(restOptions *api.RESTOptions, FirewallModel api.FwRuleMod) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}

	return client.Firewall().Create(ctx, FirewallModel)
}
