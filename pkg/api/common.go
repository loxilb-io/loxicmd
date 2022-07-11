package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type CommonAPI struct {
	restClient  *RESTClient
	requestInfo RequestInfo
}

func (l *CommonAPI) GetUrlString() string {
	lbURL := url.URL{
		Scheme: l.restClient.GetProcotol(),
		Host:   l.restClient.GetHost(),
		Path:   l.requestInfo.GetBaseURL(),
	}

	return lbURL.String()
}

func (l *CommonAPI) Create(ctx context.Context, modelbody interface{}) (*http.Response, error) {
	body, err := json.Marshal(modelbody)
	if err != nil {
		// need validation check
		return nil, err
	}
	createURL := l.GetUrlString()
	return l.restClient.POST(ctx, createURL, body)
}

func (l *CommonAPI) Delete(ctx context.Context) (*http.Response, error) {
	deleteURL := l.GetUrlString()
	return l.restClient.DELETE(ctx, deleteURL)
}

func (l *CommonAPI) Get(ctx context.Context) (*http.Response, error) {
	getURL := l.GetUrlString()
	return l.restClient.GET(ctx, getURL)
}

func (l *CommonAPI) SubResources(resourceList []string) *CommonAPI {
	l.requestInfo.subResource = append(l.requestInfo.subResource, resourceList...)
	return l
}

func (l *CommonAPI) SetUrl(url string) *CommonAPI {
	l.requestInfo.resource = url
	return l
}
