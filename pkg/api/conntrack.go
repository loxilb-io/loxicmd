package api

type Conntrack struct {
	CommonAPI
}

type CtInformationGet struct {
	CtInfo []ConntrackInformation `json:"ctAttr"`
}

type ConntrackInformation struct {
	Dip    string `json:"destinationIP"`
	Sip    string `json:"sourceIP"`
	Dport  uint16 `json:"destinationPort"`
	Sport  uint16 `json:"sourcePort"`
	Proto  string `json:"protocol"`
	CState string `json:"conntrackState"`
	CAct   string `json:"conntrackAct"`
	Pkts   uint64 `json:"packets"`
	Bytes  uint64 `json:"bytes"`
}
