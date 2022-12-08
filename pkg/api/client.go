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
	"net/http"
	"time"
)

const (
	loxiProvider             = "netlox"
	loxiApiVersion           = "v1"
	loxiLoadBalancerResource = "config/loadbalancer"
	loxiConntrackResource    = "config/conntrack/all"
	loxiPortResource         = "config/port/all"
	loxiSessionResource      = "config/session"
	loxiSessionUlClResource  = "config/sessionulcl"
	loxiPolicyResource       = "config/policy"
	loxiRouteResource        = "config/route"
	loxiIPv4AddressResource  = "config/ipv4address"
	loxiNeighborResource     = "config/neighbor"
	loxiFDBResource          = "config/fdb"
	loxiVlanResource         = "config/vlan"
	loxiVxlanResource        = "config/tunnel/vxlan"
	loxiMirrorResource       = "config/mirror"
	loxiFirewallResource     = "config/firewall"
	loxiStatusResource       = "status"
)

type LoxiClient struct {
	restClient RESTClient
}

func NewLoxiClient(o *RESTOptions) *LoxiClient {
	return &LoxiClient{
		restClient: RESTClient{
			Options: *o,
			Client: &http.Client{
				Timeout: time.Second * time.Duration(o.Timeout),
			},
		},
	}
}

func (l *LoxiClient) LoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		CommonAPI: CommonAPI{
			restClient: &l.restClient,
			requestInfo: RequestInfo{
				provider:   loxiProvider,
				apiVersion: loxiApiVersion,
				resource:   loxiLoadBalancerResource,
			},
		},
	}
}

func (l *LoxiClient) Conntrack() *Conntrack {
	return &Conntrack{
		CommonAPI: CommonAPI{
			restClient: &l.restClient,
			requestInfo: RequestInfo{
				provider:   loxiProvider,
				apiVersion: loxiApiVersion,
				resource:   loxiConntrackResource,
			},
		},
	}
}

func (l *LoxiClient) Port() *Port {
	return &Port{
		CommonAPI: CommonAPI{
			restClient: &l.restClient,
			requestInfo: RequestInfo{
				provider:   loxiProvider,
				apiVersion: loxiApiVersion,
				resource:   loxiPortResource,
			},
		},
	}
}

func (l *LoxiClient) Session() *Session {
	return &Session{
		CommonAPI: CommonAPI{
			restClient: &l.restClient,
			requestInfo: RequestInfo{
				provider:   loxiProvider,
				apiVersion: loxiApiVersion,
				resource:   loxiSessionResource,
			},
		},
	}
}

func (l *LoxiClient) SessionUlCL() *SessionUlCl {
	return &SessionUlCl{
		CommonAPI: CommonAPI{
			restClient: &l.restClient,
			requestInfo: RequestInfo{
				provider:   loxiProvider,
				apiVersion: loxiApiVersion,
				resource:   loxiSessionUlClResource,
			},
		},
	}
}

func (l *LoxiClient) Policy() *Policy {
	return &Policy{
		CommonAPI: CommonAPI{
			restClient: &l.restClient,
			requestInfo: RequestInfo{
				provider:   loxiProvider,
				apiVersion: loxiApiVersion,
				resource:   loxiPolicyResource,
			},
		},
	}
}

func (l *LoxiClient) Route() *Route {
	return &Route{
		CommonAPI: CommonAPI{
			restClient: &l.restClient,
			requestInfo: RequestInfo{
				provider:   loxiProvider,
				apiVersion: loxiApiVersion,
				resource:   loxiRouteResource,
			},
		},
	}
}

func (l *LoxiClient) IPv4Address() *IPv4Address {
	return &IPv4Address{
		CommonAPI: CommonAPI{
			restClient: &l.restClient,
			requestInfo: RequestInfo{
				provider:   loxiProvider,
				apiVersion: loxiApiVersion,
				resource:   loxiIPv4AddressResource,
			},
		},
	}
}

func (l *LoxiClient) Neighbor() *Neighbor {
	return &Neighbor{
		CommonAPI: CommonAPI{
			restClient: &l.restClient,
			requestInfo: RequestInfo{
				provider:   loxiProvider,
				apiVersion: loxiApiVersion,
				resource:   loxiNeighborResource,
			},
		},
	}
}

func (l *LoxiClient) FDB() *FDB {
	return &FDB{
		CommonAPI: CommonAPI{
			restClient: &l.restClient,
			requestInfo: RequestInfo{
				provider:   loxiProvider,
				apiVersion: loxiApiVersion,
				resource:   loxiFDBResource,
			},
		},
	}
}

func (l *LoxiClient) Vlan() *Vlan {
	return &Vlan{
		CommonAPI: CommonAPI{
			restClient: &l.restClient,
			requestInfo: RequestInfo{
				provider:   loxiProvider,
				apiVersion: loxiApiVersion,
				resource:   loxiVlanResource,
			},
		},
	}
}

func (l *LoxiClient) Vxlan() *Vxlan {
	return &Vxlan{
		CommonAPI: CommonAPI{
			restClient: &l.restClient,
			requestInfo: RequestInfo{
				provider:   loxiProvider,
				apiVersion: loxiApiVersion,
				resource:   loxiVxlanResource,
			},
		},
	}
}

func (l *LoxiClient) Status() *Status {
	return &Status{
		CommonAPI: CommonAPI{
			restClient: &l.restClient,
			requestInfo: RequestInfo{
				provider:   loxiProvider,
				apiVersion: loxiApiVersion,
				resource:   loxiStatusResource,
			},
		},
	}
}

func (l *LoxiClient) Firewall() *Firewall {
	return &Firewall{
		CommonAPI: CommonAPI{
			restClient: &l.restClient,
			requestInfo: RequestInfo{
				provider:   loxiProvider,
				apiVersion: loxiApiVersion,
				resource:   loxiFirewallResource,
			},
		},
	}
}

func (l *LoxiClient) Mirror() *Mirror {
	return &Mirror{
		CommonAPI: CommonAPI{
			restClient: &l.restClient,
			requestInfo: RequestInfo{
				provider:   loxiProvider,
				apiVersion: loxiApiVersion,
				resource:   loxiMirrorResource,
			},
		},
	}
}
