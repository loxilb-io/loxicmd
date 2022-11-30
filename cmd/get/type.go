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
	LOADBALANCER_TITLE      = []string{"External IP", "Port", "Protocol", "Select", "Mode", "# of Endpoints"}
	LOADBALANCER_WIDE_TITLE = []string{"External IP", "Port", "Protocol", "Select", "Mode", "Endpoint IP", "Target Port", "Weight"}
	SESSION_TITLE           = []string{"ident", "session IP"}
	SESSION_WIDE_TITLE      = []string{"ident", "session IP", "access Network Tunnel", "core Network Tunnel"}
	PORT_WIDE_TITLE         = []string{"index", "portname", "MAC", "link/state", "mtu", "isActive/bpf", "Statistics", "L3Info", "L2Info", "Sync"}
	PORT_TITLE              = []string{"index", "portname", "MAC", "link/state", "L3Info", "L2Info"}
	ULCL_TITLE              = []string{"ident", "ulcl IP", "qfi"}
	POLICY_TITLE            = []string{"Ident", "peakInfoRate", "committedInfoRate"}
	POLICY_WIDE_TITLE       = []string{"Ident", "peakInfoRate", "committedInfoRate", "excessBlkSize", "committedBlkSize", "policyType", "ColorAware", "polObjName", "attachment"}
)
