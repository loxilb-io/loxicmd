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

type Vlan struct {
	CommonAPI
}
type VlanGet struct {
	Vlans []VlanDump `json:"vlanAttr"`
}

type VlanDump struct {
	Vid       int             `json:"vid"`
	Dev       string          `json:"dev"`
	Member    []VlanMemberMod `json:"member"`
	Statistic VlanStat        `json:"vlanStatistic"`
}

// vlanStat - statistics for vlan interface
type VlanStat struct {
	InBytes    uint64
	InPackets  uint64
	OutBytes   uint64
	OutPackets uint64
}

// VlanBridgerMod - Info about an Vlan bridge
type VlanBridgeMod struct {
	// Vid - Virtual LAN ID
	Vid int `json:"vid"`
}

// VlanMemberMod - Info about an Vlan bridge member
type VlanMemberMod struct {
	// Dev - name of the related device
	Dev string `json:"dev"`
	// Tagged - Tagging status of the device
	Tagged bool `json:"Tagged"`
}
