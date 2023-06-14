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
package get

var (
	CONNTRACK_TITLE         = []string{"destination IP", "source IP", "destination Port", "source Port", "protocol", "state", "act", "packets", "bytes"}
	LOADBALANCER_TITLE      = []string{"External IP", "Port", "Protocol", "Block", "Select", "Mode", "# of Endpoints", "Monitor"}
	LOADBALANCER_WIDE_TITLE = []string{"External IP", "Secondary IPs", "Port", "Protocol", "block", "Select", "Mode", "Endpoint IP", "Target Port", "Weight", "State"}
	SESSION_TITLE           = []string{"ident", "session IP"}
	SESSION_WIDE_TITLE      = []string{"ident", "session IP", "access Network Tunnel", "core Network Tunnel"}
	PORT_WIDE_TITLE         = []string{"index", "portname", "MAC", "link/state", "mtu", "isActive/bpf\nPort type", "Statistics", "L3Info", "L2Info", "Sync"}
	PORT_TITLE              = []string{"index", "portname", "MAC", "link/state", "L3Info", "L2Info"}
	ULCL_TITLE              = []string{"ident", "ulcl IP", "qfi"}
	POLICY_TITLE            = []string{"Ident", "peakInfoRate", "committedInfoRate"}
	POLICY_WIDE_TITLE       = []string{"Ident", "peakInfoRate", "committedInfoRate", "excessBlkSize", "committedBlkSize", "policyType", "ColorAware", "polObjName", "attachment"}
	ROUTE_TITLE             = []string{"destinationIPNet", "gateway", "flag"}
	ROUTE_WIDE_TITLE        = []string{"destinationIPNet", "gateway", "flag", "HardwareMark", "packets", "bytes"}
	IP_TITLE                = []string{"Device Name", "IP Address"}
	FDB_TITLE               = []string{"Device Name", "MAC Address"}
	IP_WIDE_TITLE           = []string{"Device Name", "IP Address", "Sync"}
	VLAN_WIDE_TITLE         = []string{"Device Name", "Vlan ID", "Member", "Statistics"}
	VLAN_TITLE              = []string{"Device Name", "Vlan ID", "Member"}
	VXLAN_TITLE             = []string{"Device Name", "Vxlan ID", "endpoint interface", "Peer IP"}
	NEIGHBOR_TITLE          = []string{"IP Address", "Device Name", "Mac Address"}
	PROCESS_TITLE           = []string{"pid", "user", "priority", "nice", "virtMemory", "residentSize", "sharedMemory", "status", "CPUUsage", "MemoryUsage", "time", "command"}
	DEVICE_TITLE            = []string{"hostName", "machineID", "bootID", "OS", "kernel", "architecture", "uptime"}
	FILESYSTEM_TITLE        = []string{"fileSystem", "type", "size", "used", "avail", "usePercent", "mountedOn"}
	MIRROR_TITLE            = []string{"Mirror Name", "Mirror info", "Target\nAttachment", "target\nName"}
	MIRROR_WIDE_TITLE       = []string{"Mirror Name", "Mirror info", "Target\nAttachment", "target\nName", "Sync"}
	FIREWALL_TITLE          = []string{"Source IP", "destination IP", "min SPort", "max SPort", "min DPort", "max DPort", "protocol", "port Name", "preference", "Option"}
	ENDPOINT_TITLE          = []string{"Host", "Name", "ptype", "port", "duration", "retries", "minDelay", "avgDelay", "maxDelay", "State"}
	PARAM_TITLE             = []string{"Param Name", "Value"}
)
