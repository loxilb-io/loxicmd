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

import "sort"

type Mirror struct {
	CommonAPI
}
type MirrorGet struct {
	Mirrors []MirrGetMod `json:"mirrAttr"`
}

const (
	// MirrTypeSpan - simple SPAN
	MirrTypeSpan = 0 // Default
	// MirrTypeRspan - type RSPAN
	MirrTypeRspan = 1
	// MirrTypeErspan - type ERSPAN
	MirrTypeErspan = 2
)

// MirrInfo - information related to a mirror entry
type MirrInfo struct {
	// MirrType - one of MirrTypeSpan, MirrTypeRspan or MirrTypeErspan
	MirrType int `json:"type" yaml:"type"`
	// MirrPort - port where mirrored traffic needs to be sent
	MirrPort string `json:"port" yaml:"port"`
	// MirrVlan - for RSPAN we may need to send tagged mirror traffic
	MirrVlan int `json:"vlan" yaml:"vlan"`
	// MirrRip - RemoteIP. For ERSPAN we may need to send tunnelled mirror traffic
	MirrRip string `json:"remoteIP" yaml:"remoteIP"`
	// MirrRip - SourceIP. For ERSPAN we may need to send tunnelled mirror traffic
	MirrSip string `json:"sourceIP" yaml:"sourceIP"`
	// MirrTid - mirror tunnel-id. For ERSPAN we may need to send tunnelled mirror traffic
	MirrTid int `json:"tunnelID" yaml:"tunnelID"`
}

// MirrObjType - type of mirror attachment
type MirrObjType int

const (
	// MirrAttachPort - mirror attachment to a port
	MirrAttachPort MirrObjType = 1 << iota
	// MirrAttachRule - mirror attachment to a lb rule
	MirrAttachRule
)

// MirrObj - information of object attached to mirror
type MirrObj struct {
	// MirrObjName - object name to be attached to mirror
	MirrObjName string `json:"mirrObjName" yaml:"mirrObjName"`
	// AttachMent - one of MirrAttachPort or MirrAttachRule
	AttachMent MirrObjType `json:"attachment" yaml:"attachment"`
}

// MirrMod - information related to a  mirror entry
type MirrMod struct {
	// Ident - unique identifier for the mirror
	Ident string `json:"mirrorIdent" yaml:"mirrorIdent"`
	// Info - information about the mirror
	Info MirrInfo `json:"mirrorInfo" yaml:"mirrorInfo"`
	// Target - information about object to which mirror needs to be attached
	Target MirrObj `json:"targetObject" yaml:"targetObject"`
}

// MirrGetMod - information related to Get a mirror entry
type MirrGetMod struct {
	// Ident - unique identifier for the mirror
	Ident string `json:"mirrorIdent"`
	// Info - information about the mirror
	Info MirrInfo `json:"mirrorInfo"`
	// Target - information about object to which mirror needs to be attached
	Target MirrObj `json:"targetObject"`
	// Sync - sync state
	Sync DpStatusT `json:"sync"`
}
type ConfigurationMirrorFile struct {
	TypeMeta   `yaml:",inline"`
	ObjectMeta `yaml:"metadata,omitempty"`
	Spec       MirrMod `yaml:"spec"`
}

func (Mirrorresp MirrorGet) Sort() {
	sort.Slice(Mirrorresp.Mirrors, func(i, j int) bool {
		return Mirrorresp.Mirrors[i].Ident < Mirrorresp.Mirrors[j].Ident
	})
}
