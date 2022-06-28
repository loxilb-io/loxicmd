package api

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
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
	subResource map[string]string
	queryArgs   map[string]string
}

func (r *RequestInfo) makeBaseURL() string {
	p := path.Join(r.provider, r.apiVersion)
	if len(r.resource) != 0 {
		p = path.Join(p, r.resource)
	}

	if len(r.subResource) != 0 {
		subP := []string{p}
		for key, value := range r.subResource {
			subP = append(subP, key, value)
		}
		p = path.Join(subP...)
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
	Protocol   string
	ServerIP   string
	ServerPort int16
	Timeout    int16
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

func (r *RESTClient) POST(ctx context.Context, postURL string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	// move RESTOptions
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("post success: %v", string(result))
	return nil
}

func (r *RESTClient) DELETE(ctx context.Context, deleteURL string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, deleteURL, nil)
	if err != nil {
		return err
	}

	// move RESTOptions
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("delete success: %v", string(result))
	return nil
}
