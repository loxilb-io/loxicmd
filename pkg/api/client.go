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
)

type LoxiClient struct {
	restClient RESTClient
}

// 나중에는 options 보고 httpClient 생성하도록 바꿉시다.
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
