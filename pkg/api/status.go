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

type Status struct {
	CommonAPI
}

type ProcessGet struct {
	ProcessAttr []Process `json:"processAttr"`
}

type FilesystemGet struct {
	FilesystemAttr []Filesystem `json:"filesystemAttr"`
}

type DeviceGet struct {
	HostName     string
	MachineID    string
	BootID       string
	OS           string
	Kernel       string
	Architecture string
	Uptime       string
}

type Process struct {
	Pid          string
	User         string
	Priority     string
	Nice         string
	VirtMemory   string
	ResidentSize string
	SharedMemory string
	Status       string
	CPUUsage     string
	MemoryUsage  string
	ProcessTime  string `json:"time"`
	Command      string
}

type Filesystem struct {
	FileSystem string
	Fstype     string `json:"type"`
	Size       string
	Used       string
	Avail      string
	UsePercent string
	MountedOn  string
}
