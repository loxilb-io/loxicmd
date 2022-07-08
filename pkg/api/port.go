package api

type Port struct {
	CommonAPI
}

type DpStatusT uint8

type PortProp uint8

type PortGet struct {
	Ports []PortDump `json:"portAttr"`
}

type PortDump struct {
	Name   string         `json:"portName"`
	PortNo int            `json:"portNo"`
	Zone   string         `json:"zone"`
	SInfo  PortSwInfo     `json:"portSoftwareInformation"`
	HInfo  PortHwInfo     `json:"portHardwareInformation"`
	Stats  PortStatsInfo  `json:"portStatisticInformation"`
	L3     PortLayer3Info `json:"portL3Information"`
	L2     PortLayer2Info `json:"portL2Information"`
	Sync   DpStatusT      `json:"DataplaneSync"`
}

type PortStatsInfo struct {
	RxBytes   uint64 `json:"rxBytes"`
	TxBytes   uint64 `json:"txBytes"`
	RxPackets uint64 `json:"rxPackets"`
	TxPackets uint64 `json:"txPackets"`
	RxError   uint64 `json:"rxErrors"`
	TxError   uint64 `json:"txErrors"`
}

type PortHwInfo struct {
	MacAddr    [6]byte `json:"rawMacAddress"`
	MacAddrStr string  `json:"macAddress"`
	Link       bool    `json:"link"`
	State      bool    `json:"state"`
	Mtu        int     `json:"mtu"`
	Master     string  `json:"master"`
	Real       string  `json:"real"`
	TunId      uint32  `json:"tunnelId"`
}

type PortLayer3Info struct {
	Routed     bool     `json:"routed"`
	Ipv4_addrs []string `json:"IPv4Address"`
	Ipv6_addrs []string `json:"IPv6Address"`
}

type PortSwInfo struct {
	OsId       int       `json:"osId"`
	PortType   int       `json:"portType"`
	PortProp   PortProp  `json:"portProp"`
	PortActive bool      `json:"portActive"`
	PortReal   *PortDump `json:"portReal"`
	PortOvl    *PortDump `json:"portOvl"`
	BpfLoaded  bool      `json:"bpfLoaded"`
}

type PortLayer2Info struct {
	IsPvid bool `json:"isPvid"`
	Vid    int  `json:"vid"`
}
