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

type LoadBalancer struct {
	CommonAPI
}

type EpSelect uint
type LbMode int32
type LbRuleModGet struct {
	LbRules []LoadBalancerModel `json:"lbAttr"`
}

type LoadBalancerModel struct {
	Service   LoadBalancerService    `json:"serviceArguments" yaml:"serviceArguments"`
	SecIPs    []LoadBalancerSecIp    `json:"secIPs" yaml:"secIPs"`
	Endpoints []LoadBalancerEndpoint `json:"endpoints" yaml:"endpoints"`
}

type LoadBalancerService struct {
	ExternalIP string   `json:"externalIP" yaml:"externalIP"`
	Port       uint16   `json:"port"           yaml:"port" `
	Protocol   string   `json:"protocol"       yaml:"protocol"`
	Sel        EpSelect `json:"sel"            yaml:"sel"`
	Mode       LbMode   `json:"mode"           yaml:"mode"`
	BGP        bool     `json:"BGP"            yaml:"BGP"`
	Monitor    bool     `json:"Monitor"        yaml:"Monitor"`
	Timeout    uint32   `json:"inactiveTimeOut" yaml:"inactiveTimeOut"`
	Block      uint16   `json:"block"          yaml:"block"`
}

type LoadBalancerEndpoint struct {
	EndpointIP string `json:"endpointIP" yaml:"endpointIP"`
	TargetPort uint16 `json:"targetPort" yaml:"targetPort"`
	Weight     uint8  `json:"weight"     yaml:"weight"`
	State      string `json:"state"      yaml:"state"`
}

type LoadBalancerSecIp struct {
	SecIP string `json:"secIP" yaml:"secIP"`
}

type ConfigurationLBFile struct {
	TypeMeta   `yaml:",inline"`
	ObjectMeta `yaml:"metadata,omitempty"`
	Spec       LoadBalancerModel `yaml:"spec"`
}
