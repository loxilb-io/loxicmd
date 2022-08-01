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
