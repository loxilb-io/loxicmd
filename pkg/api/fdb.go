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

type FDB struct {
	CommonAPI
}

type FDBModGet struct {
	FdbAttr []FDBMod `json:"fdbAttr"`
}

type FDBMod struct {
	// Dev - name of the related device
	Dev string `json:"dev" yaml:"dev"`
	// MacAddress - Hardware address
	MacAddress string `json:"macAddress" yaml:"macAddress"`
}

type ConfigurationFDBFile struct {
	TypeMeta   `yaml:",inline"`
	ObjectMeta `yaml:"metadata,omitempty"`
	Spec       FDBMod `yaml:"spec"`
}

func (fdb FDBMod) Key() string {
	return fmt.Sprintf("%s|%s", fdb.Dev, fdb.MacAddress)
}

func (FDBresp FDBModGet) Sort() {
	sort.Slice(FDBresp.FdbAttr, func(i, j int) bool {
		return FDBresp.FdbAttr[i].Key() < FDBresp.FdbAttr[j].Key()
	})
}
