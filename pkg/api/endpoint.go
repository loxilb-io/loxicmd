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

type EndPoint struct {
	CommonAPI
}

type EPInformationGet struct {
	EPInfo []EndPointGetEntry `json:"Attr"`
}

// EndPointState - Info related to a end-point state
type EndPointGetEntry struct {
	// Name - hostname in CIDR
	Name string `json:"hostName"`
	// Desc - host specific description
	Desc string `json:"description"`
	// InActTries - No. of inactive probes to mark
	// an end-point inactive
	InActTries int `json:"inactiveReTries"`
	// ProbeType - Type of probe : "icmp","connect-tcp", "connect-udp", "connect-sctp", "http"
	ProbeType string `json:"probeType"`
	// ProbeReq - Request string in case of http probe
	ProbeReq string `json:"probeReq"`
	// ProbeResp - Response string in case of http probe
	ProbeResp string `json:"probeResp"`
	// ProbeDuration - How frequently (in seconds) to check activity
	ProbeDuration uint32 `json:"probeDuration"`
	// ProbePort - Port to probe for connect type
	ProbePort uint16 `json:"probePort"`
	// MinDelay - Minimum delay in this end-point
	MinDelay string `json:"minDelay"`
	// AvgDelay - Average delay in this end-point
	AvgDelay string `json:"avgDelay"`
	// MaxDelay - Max delay in this end-point
	MaxDelay string `json:"maxDelay"`
	// CurrState - Current state of this end-point
	CurrState string `json:"currState"`
}

type EPConfig struct {
	EPInfo []EndPointMod `json:"Attr"`
}

// EndPointMod - Info related to a end-point config entry
type EndPointMod struct {
	// Name - hostname in CIDR
	Name string `json:"hostName" yaml:"hostName"`
	// Desc - host specific description
	Desc string `json:"description" yaml:"description"`
	// InActTries - No. of inactive probes to mark
	// an end-point inactive
	InActTries int `json:"inactiveReTries" yaml:"inactiveReTries"`
	// ProbeType - Type of probe : "icmp","connect-tcp", "connect-udp", "connect-sctp", "http"
	ProbeType string `json:"probeType" yaml:"probeType"`
	// ProbeReq - Request string in case of http probe
	ProbeReq string `json:"probeReq" yaml:"probeReq"`
	// ProbeResp - Response string in case of http probe
	ProbeResp string `json:"probeResp" yaml:"probeResp"`
	// ProbeDuration - How frequently (in seconds) to check activity
	ProbeDuration uint32 `json:"probeDuration" yaml:"probeDuration"`
	// ProbePort - Port to probe for connect type
	ProbePort uint16 `json:"probePort" yaml:"probePort"`
}

type ConfigurationEndPointFile struct {
	TypeMeta   `yaml:",inline"`
	ObjectMeta `yaml:"metadata,omitempty"`
	Spec       EndPointMod `yaml:"spec"`
}
