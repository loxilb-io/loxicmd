package api

import (
	"context"
	_ "fmt"
	"net/http"
	"net/url"
)

type Conntrack struct {
	restClient  *RESTClient
	requestInfo RequestInfo
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
}

func (l *Conntrack) GetUrlString() string {
	ctURL := url.URL{
		Scheme: l.restClient.GetProcotol(),
		Host:   l.restClient.GetHost(),
		Path:   l.requestInfo.GetBaseURL(),
	}

	return ctURL.String()
}

func (l *Conntrack) Get(ctx context.Context) (*http.Response, error) {
	getURL := l.GetUrlString()
	return l.restClient.GET(ctx, getURL)
}

func (l *Conntrack) SubResources(resourceList []string) *Conntrack {
	l.requestInfo.subResource = append(l.requestInfo.subResource, resourceList...)
	return l
}

func (l *Conntrack) SetUrl(url string) *Conntrack {
	l.requestInfo.resource = url
	return l
}
