package api

import (
	_ "fmt"
)

type LoadBalancer struct {
	CommonAPI
}

type EpSelect uint

type LbRuleModGet struct {
	LbRules []LoadBalancerModel `json:"lbAttr"`
}

type LoadBalancerModel struct {
	Service   LoadBalancerService    `json:"serviceArguments"`
	Endpoints []LoadBalancerEndpoint `json:"endpoints"`
}

type LoadBalancerService struct {
	ExternalIP string   `json:"externalIP"`
	Port       uint16   `json:"port"`
	Protocol   string   `json:"protocol"`
	Sel        EpSelect `json:"sel"`
}

type LoadBalancerEndpoint struct {
	EndpointIP string `json:"endpointIP"`
	TargetPort uint16 `json:"targetPort"`
	Weight     uint8  `json:"weight"`
}
