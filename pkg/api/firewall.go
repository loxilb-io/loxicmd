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
package api

type Firewall struct {
	CommonAPI
}

type FWInformationGet struct {
	FWInfo []FwRuleMod `json:"fwAttr"`
}

// FwRuleOpts - Information related to Firewall options
type FwOptArg struct {
	// Drop - Drop any matching rule
	Drop bool `json:"drop"`
	// Trap - Trap anything matching rule
	Trap bool `json:"trap"`
	// Redirect - Redirect any matching rule
	Rdr     bool   `json:"redirect"`
	RdrPort string `json:"redirectPortName"`
	// Allow - Allow any matching rule
	Allow bool `json:"allow"`
	Mark int `json:"fwMark"`
}

// FwRuleArg - Information related to firewall rule
type FwRuleArg struct {
	// SrcIP - Source IP in CIDR notation
	SrcIP string `json:"sourceIP"`
	// DstIP - Destination IP in CIDR notation
	DstIP string `json:"destinationIP"`
	// SrcPortMin - Minimum source port range
	SrcPortMin uint16 `json:"minSourcePort"`
	// SrcPortMax - Maximum source port range
	SrcPortMax uint16 `json:"maxSourcePort"`
	// DstPortMin - Minimum destination port range
	DstPortMin uint16 `json:"minDestinationPort"`
	// SrcPortMax - Maximum source port range
	DstPortMax uint16 `json:"maxDestinationPort"`
	// Proto - the protocol
	Proto uint8 `json:"protocol"`
	// InPort - the incoming port
	InPort string `json:"portName"`
	// Pref - User preference for ordering
	Pref uint16 `json:"preference"`
}

// FwRuleMod - Info related to a firewall entry
type FwRuleMod struct {
	// Serv - service argument of type FwRuleArg
	Rule FwRuleArg `json:"ruleArguments"`
	// Opts - firewall options
	Opts FwOptArg `json:"opts"`
}
