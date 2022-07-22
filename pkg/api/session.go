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
