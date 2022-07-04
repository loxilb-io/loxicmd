package api

import (
	"context"
	"encoding/json"
	_ "fmt"
	"net/http"
	"net/url"
)

type LoadBalancer struct {
	restClient  *RESTClient
	requestInfo RequestInfo
}

type EpSelect uint

type LbRuleModGet struct {
	LbRules []LoadBalancerModel    `json:"lbAttr"`
	CtInfos []ConntrackInformation `json:"conntrackAttr`
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

type ConntrackInformation struct {
	Dip    string `json:"destinationIP"`
	Sip    string `json:"sourceIP"`
	Dport  uint16 `json:"destinationPort"`
	Sport  uint16 `json:"sourcePort"`
	Proto  string `json:"protocol"`
	CState string `json:"conntrackState"`
	CAct   string `json:"conntrackAct"`
}

func (l *LoadBalancer) GetUrlString() string {
	lbURL := url.URL{
		Scheme: l.restClient.GetProcotol(),
		Host:   l.restClient.GetHost(),
		Path:   l.requestInfo.GetBaseURL(),
	}

	return lbURL.String()
}

func (l *LoadBalancer) Create(ctx context.Context, lbModel LoadBalancerModel) (*http.Response, error) {
	body, err := json.Marshal(lbModel)
	if err != nil {
		// need validation check
		return nil, err
	}
	createURL := l.GetUrlString()
	return l.restClient.POST(ctx, createURL, body)
}

func (l *LoadBalancer) Delete(ctx context.Context) (*http.Response, error) {
	deleteURL := l.GetUrlString()
	return l.restClient.DELETE(ctx, deleteURL)
}

func (l *LoadBalancer) Get(ctx context.Context) (*http.Response, error) {
	getURL := l.GetUrlString()
	return l.restClient.GET(ctx, getURL)
}

func (l *LoadBalancer) SubResources(resourceList []string) *LoadBalancer {
	l.requestInfo.subResource = append(l.requestInfo.subResource, resourceList...)
	return l
}

func (l *LoadBalancer) SetUrl(url string) *LoadBalancer {
	l.requestInfo.resource = url
	return l
}
