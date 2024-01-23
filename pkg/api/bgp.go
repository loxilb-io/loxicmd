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

import (
	"fmt"
	"sort"
)

type BGPNeighbor struct {
	CommonAPI
}

type BGPNeighborModGet struct {
	BGPAttr []BGPNeighborEntry `json:"bgpNeiAttr"`
}

type BGPNeighborMod struct {
	// IPaddress - Neighbor IP address
	IPaddress string `json:"ipAddress" yaml:"ipAddress"`
	// RemoteAs - ASN of remote site
	RemoteAs int `json:"remoteAs" yaml:"remoteAs"`
	// RemotePort - BGP port
	RemotePort int `json:"remotePort" yaml:"remotePort"`
	// SetMultiHop - BGP Multihop enable
	SetMultiHop bool `json:"setMultiHop" yaml:"setMultiHop"`
}

type BGPNeighborEntry struct {
	// IPaddress - Neighbor IP address
	IPaddress string `json:"ipAddress" yaml:"ipAddress"`
	// Status - BGP connection status
	State string `json:"state" yaml:"state"`
	// RemoteAs - ASN of remote site
	RemoteAs int `json:"remoteAs" yaml:"remoteAs"`
	// UpDownTime - uptime or down time based on status
	UpDownTime string `json:"updowntime" yaml:"updowntime"`
}

type ConfigurationBGPFile struct {
	TypeMeta   `yaml:",inline"`
	ObjectMeta `yaml:"metadata,omitempty"`
	Spec       BGPNeighborMod `yaml:"spec"`
}

func (nei BGPNeighborEntry) Key() string {
	return fmt.Sprintf("%s|%s|%s", nei.IPaddress, nei.State, nei.UpDownTime)
}

func (BGPsresp BGPNeighborModGet) Sort() {
	sort.Slice(BGPsresp.BGPAttr, func(i, j int) bool {
		return BGPsresp.BGPAttr[i].Key() < BGPsresp.BGPAttr[j].Key()
	})
}
