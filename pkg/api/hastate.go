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
	"sort"
)

type HAState struct {
	CommonAPI
}

// HAStateGetEntry c i status get entry
//
// swagger:model CIStatusGetEntry
type HAStateGet struct {
	HAStateAttr []HAStateInfo `json:"Attr"`
}

type HAStateInfo struct {
	// Instance name
	Instance string `json:"instance,omitempty"`

	// Current Cluster Instance State
	State string `json:"state,omitempty"`

	// Sync - sync state
	// Required: true
	Sync *int64 `json:"sync"`

	// Instance Virtual IP address
	Vip string `json:"vip,omitempty"`
}

func (haState HAStateGet) Sort() {
	sort.Slice(haState.HAStateAttr, func(i, j int) bool {
		return haState.HAStateAttr[i].Instance < haState.HAStateAttr[j].Instance
	})
}
