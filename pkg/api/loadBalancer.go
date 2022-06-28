package api

import (
	"context"
	"encoding/json"
	_ "fmt"
	"net/url"
)

type LoadBalancer struct {
	restClient  *RESTClient
	requestInfo RequestInfo
}

type LoadBalancerModel struct {
	ExternalIP string                 `json:"external_ip_address"`
	Port       int16                  `json:"port"`
	Protocol   string                 `json:"protocol"`
	Endpoints  []LoadBalancerEndpoint `json:"endpoints"`
}

type LoadBalancerEndpoint struct {
	EndpointIP string `json:"endpoint_ip_address"`
	TargetPort int16  `json:"targetPort"`
	Weight     int8   `json:"weight"`
}

func (l *LoadBalancer) GetUrlString() string {
	lbURL := url.URL{
		Scheme: l.restClient.GetProcotol(),
		Host:   l.restClient.GetHost(),
		Path:   l.requestInfo.GetBaseURL(),
	}

	return lbURL.String()
}

func (l *LoadBalancer) Create(ctx context.Context, lbModel LoadBalancerModel) error {
	body, err := json.Marshal(lbModel)
	if err != nil {
		// need validation check
		//fmt.Println("LoadBalancer Create: Failed marshaling. lbModel:")
		//fmt.Println(lbModel)
		return err
	}
	//fmt.Println("LoadBalancer Create: Success marshaling. Body:")
	//fmt.Println(string(body))
	createURL := l.GetUrlString()
	if err := l.restClient.POST(ctx, createURL, body); err != nil {
		return err
	}
	return nil
}

func (l *LoadBalancer) Delete(ctx context.Context) error {
	deleteURL := l.GetUrlString()
	if err := l.restClient.DELETE(ctx, deleteURL); err != nil {
		return err
	}
	return nil
}

func (l *LoadBalancer) Get(ctx context.Context) error {
	return nil
}

func (l *LoadBalancer) SubResources(resourceList []string) *LoadBalancer {
	l.requestInfo.subResource = append(l.requestInfo.subResource, resourceList...)
	return l
}
