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
	log "github.com/sirupsen/logrus"
)

// Wrap go-twitter
type Client struct {
	lib *twitter.Client
}

// Alias go-twitter SearchTweetParams
type SearchParams = twitter.SearchTweetParams

// NewClient creates a new Client using the provided httpClient
func NewClient(httpClient *http.Client) *Client {
	return &Client{lib: twitter.NewClient(httpClient)}
}

// Twitter Search API search
func (client *Client) Search(params *SearchParams) []twitter.Tweet {
	params.Count = 1000
	params.TweetMode = "extended"

	search, response, err := client.lib.Search.Tweets(params)
	if err != nil {
		log.WithError(err).Errorln("Searching twitter")
	}

	// DEBUG
	if log.GetLevel() == log.DebugLevel {
		// Protocol
		log.WithField("Protocol", response.Proto).Debugln("Protocol")

		// HTTP Headers
		for k, v := range response.Header {
			switch k {
			default:
				log.WithField(k, v).Debugln("HTTP Header")
			}
		}

		// Twitter Search API Metadata
		log.WithFields(log.Fields{
			"Count":       search.Metadata.Count,
			"SinceID":     search.Metadata.SinceID,
			"SinceIDStr":  search.Metadata.SinceIDStr,
			"MaxID":       search.Metadata.MaxID,
			"MaxIDStr":    search.Metadata.MaxIDStr,
			"RefreshURL":  search.Metadata.RefreshURL,
			"NextResults": search.Metadata.NextResults,
			"CompletedIn": search.Metadata.CompletedIn,
			"Query":       search.Metadata.Query,
		}).Debugln("Metadata describes a search result")

		// Rate rimits
		rl := NewRateLimit(&response.Header)
		log.WithFields(log.Fields{
			"limit":     rl.Limit,
			"remaining": rl.Remaining,
			"reset":     rl.Reset,
		}).Debugln("Rate limit")
	}
	return search.Statuses
}
