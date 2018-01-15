// Copyright © 2017 National Library of Norway
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
)

// Wrap go-twitter client
type Client struct {
	lib *twitter.Client
}

// Alias go-twitter SearchTweetParams
type Params = twitter.SearchTweetParams
type Metadata = twitter.SearchMetadata
type Tweet = twitter.Tweet
type Result = twitter.Search

type Response struct {
	Status   string              `json:"status"`
	Code     int                 `json:"code"`
	Protocol string              `json:"protocol"`
	Header   map[string][]string `json:"header"`
}

func (r *Response) fromHttpResponse(response *http.Response) *Response {
	r.Header = response.Header
	r.Status = response.Status
	r.Code = response.StatusCode
	r.Protocol = response.Proto
	return r
}

// NewClient creates a new Client using the provided httpClient
func New(httpClient *http.Client) *Client {
	return &Client{lib: twitter.NewClient(httpClient)}
}

// Search searches using Twitter's Search API
func (client *Client) Search(params *Params) (*Result, *Response, error) {
	search, response, err := client.lib.Search.Tweets(params)
	if err != nil {
		return nil, nil, err
	}
	return search, new(Response).fromHttpResponse(response), nil
}
