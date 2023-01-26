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

type Policy struct {
	CommonAPI
}

type PolObjType uint

type PolInformationGet struct {
	PolModInfo []PolMod `json:"polAttr"`
}

type PolInfo struct {
	PolType           int    `json:"type" yaml:"type"`
	ColorAware        bool   `json:"colorAware" yaml:"colorAware"`
	CommittedInfoRate uint64 `json:"committedInfoRate" yaml:"committedInfoRate"` // CIR in Mbps
	PeakInfoRate      uint64 `json:"peakInfoRate" yaml:"peakInfoRate"`           // PIR in Mbps
	CommittedBlkSize  uint64 `json:"committedBlkSize" yaml:"committedBlkSize"`   // CBS in bytes
	ExcessBlkSize     uint64 `json:"excessBlkSize" yaml:"excessBlkSize"`         // EBS in bytes
}

type PolObj struct {
	PolObjName string     `json:"polObjName" yaml:"polObjName"`
	AttachMent PolObjType `json:"attachment" yaml:"attachment"`
}

type PolMod struct {
	Ident  string  `json:"policyIdent" yaml:"policyIdent"`
	Info   PolInfo `json:"policyInfo" yaml:"policyInfo"`
	Target PolObj  `json:"targetObject" yaml:"targetObject"`
}
type ConfigurationPolicyFile struct {
	TypeMeta   `yaml:",inline"`
	ObjectMeta `yaml:"metadata,omitempty"`
	Spec       PolMod `yaml:"spec"`
}
