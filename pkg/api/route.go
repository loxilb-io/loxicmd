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

type Route struct {
	CommonAPI
}

type RouteModGet struct {
	RouteAttr []Routev4Get `json:"routeAttr"`
}

// RouteGetEntryStatistic - Info about an route statistic
type RouteGetEntryStatistic struct {
	// Statistic of the ingress port bytes.
	Bytes int `json:"bytes"`
	// Statistic of the egress port bytes.
	Packets int `json:"packets"`
}

// Routev4Get - Info about an route
type Routev4Get struct {
	// Flags - flag type
	Flags string `json:"flags" yaml:"flags"`
	// Gw - gateway information if any
	Gw string `json:"gateway" yaml:"gateway"`
	// Dst - ip addr
	Dst string `json:"destinationIPNet" yaml:"destinationIPNet"`
	// index of the route
	HardwareMark int `json:"hardwareMark" yaml:"hardwareMark"`
	// statistic
	Statistic RouteGetEntryStatistic `json:"statistic" yaml:"statistic"`
}

type ConfigurationRouteFile struct {
	TypeMeta   `yaml:",inline"`
	ObjectMeta `yaml:"metadata,omitempty"`
	Spec       Routev4Get `yaml:"spec"`
}

func (routeresp RouteModGet) Sort() {
	sort.Slice(routeresp.RouteAttr, func(i, j int) bool {
		return routeresp.RouteAttr[i].Dst < routeresp.RouteAttr[j].Dst
	})
}
