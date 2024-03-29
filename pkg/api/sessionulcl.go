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
	"net"
	"sort"
)

type SessionUlCl struct {
	CommonAPI
}

type UlclInformationGet struct {
	UlclInfo []SessionUlClMod `json:"ulclAttr"`
}

type SessionUlClMod struct {
	Ident string  `json:"ulclIdent" yaml:"ulclIdent"`
	Args  UlClArg `json:"ulclArgument" yaml:"ulclArgument"`
}

type UlClArg struct {
	Addr net.IP `json:"ulclIP" yaml:"ulclIP"`
	Qfi  uint8  `json:"qfi" yaml:"qfi"`
}

type ConfigurationSessionUlclFile struct {
	TypeMeta   `yaml:",inline"`
	ObjectMeta `yaml:"metadata,omitempty"`
	Spec       SessionUlClMod `yaml:"spec"`
}

func (ulclresp UlclInformationGet) Sort() {
	sort.Slice(ulclresp.UlclInfo, func(i, j int) bool {
		return ulclresp.UlclInfo[i].Ident < ulclresp.UlclInfo[j].Ident
	})
}
