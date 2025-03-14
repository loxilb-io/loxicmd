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
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
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
	ServiceName string
	Token       string
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

func (r *RESTClient) GET(ctx context.Context, getURL string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, err
	}
	r.getTokens()
	// move RESTOptions
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.Options.Token)
	return r.Client.Do(req)
}

func (r *RESTClient) POST(ctx context.Context, postURL string, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	r.getTokens()
	// move RESTOptions
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.Options.Token)
	return r.Client.Do(req)
}

func (r *RESTClient) DELETE(ctx context.Context, deleteURL string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, deleteURL, nil)
	if err != nil {
		return nil, err
	}
	r.getTokens()
	// move RESTOptions
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.Options.Token)
	return r.Client.Do(req)
}

func (r *RESTClient) getTokens() {
	if r.Options.Token == "" {
		token, err := os.ReadFile("/tmp/loxilbtoken")
		if err != nil {
			return
		}
		r.Options.Token = string(token)
	}
}
