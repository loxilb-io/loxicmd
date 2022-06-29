package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

const (
	DEFAULT_URL_FORMAT = "%s://%s:%d"
)

type RequestInfo struct {
	provider    string
	apiVersion  string
	resource    string
	subResource []string
	queryArgs   map[string]string
}

func (r *RequestInfo) makeBaseURL() string {
	p := path.Join(r.provider, r.apiVersion)
	if len(r.resource) != 0 {
		p = path.Join(p, r.resource)
	}

	if len(r.subResource) != 0 {
		subP := path.Join(r.subResource...)
		p = path.Join(p, subP)
	}
	return p
}

// GetBaseURL return url.URL.Path string
func (r *RequestInfo) GetBaseURL() string {
	return r.makeBaseURL()
}

// GetQueryValue return url.Values for url.URL
func (r *RequestInfo) GetQueryString() string {
	return url.Values{}.Encode()
}

type RESTOptions struct {
	PrintOption string
	Protocol    string
	ServerIP    string
	ServerPort  int16
	Timeout     int16
}

type RESTClient struct {
	Options RESTOptions
	Client  *http.Client
}

func (r *RESTClient) GetProcotol() string {
	return r.Options.Protocol
}

func (r *RESTClient) GetHost() string {
	return fmt.Sprintf("%s:%d", r.Options.ServerIP, int(r.Options.ServerPort))
}

func (r *RESTClient) POST(ctx context.Context, postURL string, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// move RESTOptions
	req.Header.Set("Content-Type", "application/json")
	return r.Client.Do(req)
}

func (r *RESTClient) DELETE(ctx context.Context, deleteURL string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, deleteURL, nil)
	if err != nil {
		return nil, err
	}

	// move RESTOptions
	req.Header.Set("Content-Type", "application/json")
	return r.Client.Do(req)
}
