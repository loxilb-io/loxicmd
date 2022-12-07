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
	"strings"
)

type Port struct {
	CommonAPI
}

type DpStatusT uint8

type PortProp uint8

type PortGet struct {
	Ports []PortDump `json:"portAttr"`
}

type PortDump struct {
	Name   string         `json:"portName"`
	PortNo int            `json:"portNo"`
	Zone   string         `json:"zone"`
	SInfo  PortSwInfo     `json:"portSoftwareInformation"`
	HInfo  PortHwInfo     `json:"portHardwareInformation"`
	Stats  PortStatsInfo  `json:"portStatisticInformation"`
	L3     PortLayer3Info `json:"portL3Information"`
	L2     PortLayer2Info `json:"portL2Information"`
	Sync   DpStatusT      `json:"DataplaneSync"`
}

type PortStatsInfo struct {
	RxBytes   uint64 `json:"rxBytes"`
	TxBytes   uint64 `json:"txBytes"`
	RxPackets uint64 `json:"rxPackets"`
	TxPackets uint64 `json:"txPackets"`
	RxError   uint64 `json:"rxErrors"`
	TxError   uint64 `json:"txErrors"`
}

type PortHwInfo struct {
	MacAddr    [6]byte `json:"rawMacAddress"`
	MacAddrStr string  `json:"macAddress"`
	Link       bool    `json:"link"`
	State      bool    `json:"state"`
	Mtu        int     `json:"mtu"`
	Master     string  `json:"master"`
	Real       string  `json:"real"`
	TunId      uint32  `json:"tunnelId"`
}

type PortLayer3Info struct {
	Routed     bool     `json:"routed"`
	Ipv4_addrs []string `json:"IPv4Address"`
	Ipv6_addrs []string `json:"IPv6Address"`
}

type PortSwInfo struct {
	OsId       int       `json:"osId"`
	PortType   int       `json:"portType"`
	PortProp   PortProp  `json:"portProp"`
	PortActive bool      `json:"portActive"`
	PortReal   *PortDump `json:"portReal"`
	PortOvl    *PortDump `json:"portOvl"`
	BpfLoaded  bool      `json:"bpfLoaded"`
}

type PortLayer2Info struct {
	IsPvid bool `json:"isPvid"`
	Vid    int  `json:"vid"`
}

const (
	// PortReal - Base port type
	PortReal = 0x1
	// PortBondSif - Bond slave port type
	PortBondSif = 0x2
	// PortBond - Bond port type
	PortBond = 0x4
	// PortVlanSif - Vlan slave port type
	PortVlanSif = 0x8
	// PortVlanBr - Vlan Br port type
	PortVlanBr = 0x10
	// PortVxlanSif - Vxlan slave port type
	PortVxlanSif = 0x20
	// PortVxlanBr - Vxlan br port type
	PortVxlanBr = 0x40
	// PortWg - Wireguard port type
	PortWg = 0x80
	// PortVti - Vti port type
	PortVti = 0x100
)

func (p PortSwInfo) PortTypeToString() string {
	var pStr string
	if p.PortType&PortReal == PortReal {
		pStr += "phy,"
	}
	if p.PortType&PortVlanSif == PortVlanSif {
		pStr += "vlan-sif,"
	}
	if p.PortType&PortVlanBr == PortVlanBr {
		pStr += "vlan,"
	}
	if p.PortType&PortBondSif == PortBondSif {
		pStr += "bond-sif,"
	}
	if p.PortType&PortBondSif == PortBond {
		pStr += "bond,"
	}
	if p.PortType&PortVxlanSif == PortVxlanSif {
		pStr += "vxlan-sif,"
	}
	if p.PortType&PortVti == PortVti {
		pStr += "vti,"
	}
	if p.PortType&PortVxlanBr == PortVxlanBr {
		pStr += "vxlan,"
		if p.PortReal != nil {
			pStr += fmt.Sprintf("(%s)", p.PortReal.Name)
		}
	}
	nStr := strings.TrimSuffix(pStr, ",")
	return nStr
}
