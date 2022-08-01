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

import "net"

type Session struct {
	CommonAPI
}

type SessionInformationGet struct {
	SessionInfo []SessionMod `json:"sessionAttr"`
}

type SessionMod struct {
	Ident string  `json:"ident"`
	Ip    net.IP  `json:"sessionIP"`
	AnTun SessTun `json:"accessNetworkTunnel"`
	CnTun SessTun `json:"connectionNetworkTunnel"`
}

type SessTun struct {
	TeID uint32 `json:"teID"`
	Addr net.IP `json:"tunnelIP"`
}
