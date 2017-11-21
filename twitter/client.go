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
	"log"
	"fmt"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
)

type Client struct {
	lib *twitter.Client
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{lib: twitter.NewClient(httpClient)}
}

func (client *Client) Search(query string) {
	searchTweetParams := &twitter.SearchTweetParams{
		Query:     query,
		TweetMode: "extended",
		Count:     1000,
	}

	search, _, err := client.lib.Search.Tweets(searchTweetParams)
	if err != nil {
		log.Println(err)
	}

	for index, tweet := range search.Statuses {
		fmt.Println(index, tweet.FullText, "\n")
	}
}
