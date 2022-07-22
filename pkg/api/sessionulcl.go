package api

import "net"

type SessionUlCl struct {
	CommonAPI
}

type UlclInformationGet struct {
	UlclInfo []SessionUlClMod `json:"ulclAttr"`
}

type SessionUlClMod struct {
	Ident string  `json:"ulclIdent"`
	Args  UlClArg `json:"ulclArgument"`
}

type UlClArg struct {
	Addr net.IP `json:"ulclIP"`
	Qfi  uint8  `json:"qfi"`
}
