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

type BFDSession struct {
	CommonAPI
}

// HAStateGetEntry c i status get entry
//
// swagger:model CIStatusGetEntry
type BFDSessionGet struct {
	BFDSessionAttr []BFDSessionInfo `json:"Attr"`
}

type BFDSessionInfo struct {
	// Instance name
	Instance string `json:"instance" yaml:"instance"`

	// RemoteIP - Remote IP for BFD session
	RemoteIP string `json:"remoteIp" yaml:"remoteIp"`

	// Interval - Tx Interval between BFD packets
	SourceIP string `json:"sourceIp" yaml:"sourceIp"`

	// Port - BFD session port
	Port uint16 `json:"port" yaml:"port"`

	// Interval - Tx Interval between BFD packets
	Interval uint64 `json:"interval" yaml:"interval"`

	// RetryCount - Retry Count for detecting failure
	RetryCount uint8 `json:"retryCount" yaml:"retryCount"`

	// Current BFD State
	State string `json:"state" yaml:"state"`
}

type ConfigurationBFDFile struct {
	TypeMeta   `yaml:",inline"`
	ObjectMeta `yaml:"metadata,omitempty"`
	Spec       BFDSessionInfo `yaml:"spec"`
}
