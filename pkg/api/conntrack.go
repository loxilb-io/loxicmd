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

type Conntrack struct {
	CommonAPI
}

type CtInformationGet struct {
	CtInfo []ConntrackInformation `json:"ctAttr"`
}

type ConntrackInformation struct {
	Dip      string `json:"destinationIP"`
	Sip      string `json:"sourceIP"`
	Dport    uint16 `json:"destinationPort"`
	Sport    uint16 `json:"sourcePort"`
	Proto    string `json:"protocol"`
	CState   string `json:"conntrackState"`
	CAct     string `json:"conntrackAct"`
	Pkts     uint64 `json:"packets"`
	Bytes    uint64 `json:"bytes"`
	ServName string `json:"servName"`
}

func (ct ConntrackInformation) Key() string {
	return fmt.Sprintf("%s|%s|%05d|%05d|%s", ct.Dip, ct.Sip, ct.Dport, ct.Sport, ct.Proto)
}

func (ctresp CtInformationGet) Sort() {
	sort.Slice(ctresp.CtInfo, func(i, j int) bool {
		return ctresp.CtInfo[i].Key() < ctresp.CtInfo[j].Key()
	})
}
