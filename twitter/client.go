// Copyright 2018 National Library of Norway
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package twitter

import (
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/pkg/errors"
)

type Client interface {
	Search(*Params) (*Result, *Response, error)
}

// Wrap go-twitter client
type lib struct {
	client *twitter.Client
}

// NewClient creates a new Client using the provided httpClient
func New(httpClient *http.Client) Client {
	return &lib{twitter.NewClient(httpClient)}
}

// Search searches using Twitter's Search API
func (l *lib) Search(params *Params) (*Result, *Response, error) {
	p := twitter.SearchTweetParams(*params)

	if search, response, err := l.client.Search.Tweets(&p); err != nil {
		return nil, nil, errors.Wrap(err, "failed to search twitter")
	} else {
		return search, new(Response).FromHttpResponse(response), nil
	}
}
