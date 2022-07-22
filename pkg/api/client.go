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
