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

type Neighbor struct {
	CommonAPI
}

type NeighborModGet struct {
	NeighborAttr []NeighborMod `json:"neighborAttr"`
}

type NeighborMod struct {
	// Dev - name of the related device
	Dev string `json:"dev" yaml:"dev"`
	// IP - Actual IP address
	IP string `json:"ipAddress" yaml:"ipAddress"`
	// MacAddress - Hardware address
	MacAddress string `json:"macAddress" yaml:"macAddress"`
}
type ConfigurationNeighborFile struct {
	TypeMeta   `yaml:",inline"`
	ObjectMeta `yaml:"metadata,omitempty"`
	Spec       NeighborMod `yaml:"spec"`
}

func (nei NeighborMod) Key() string {
	return fmt.Sprintf("%s|%s|%s", nei.IP, nei.Dev, nei.MacAddress)
}

func (Neighborsresp NeighborModGet) Sort() {
	sort.Slice(Neighborsresp.NeighborAttr, func(i, j int) bool {
		return Neighborsresp.NeighborAttr[i].Key() < Neighborsresp.NeighborAttr[j].Key()
	})
}
