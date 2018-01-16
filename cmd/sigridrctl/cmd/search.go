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

package cmd

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nlnwa/sigridr/twitter"
	"github.com/nlnwa/sigridr/auth"
	"github.com/nlnwa/sigridr/twitter/ratelimit"
)

var (
	count int
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search query ...",
	Short: "Query Twitter's Search API",
	Long:  `Query Twitter's Search API`,
	Run: func(cmd *cobra.Command, args []string) {
		query := strings.Join(args, " ")

		params := &twitter.Params{
			Query:      query,
			Count:      count,
			ResultType: "recent",
			TweetMode:  "extended",
		}

		// Get authorized httpClient and set timeout
		httpClient := auth.HttpClient(viper.Get("token"))
		httpClient.Timeout = 10 * time.Second

		// Get twitter client
		client := twitter.New(httpClient)

		// Search twitter using params
		result, response, err := client.Search(params)
		if err != nil {
			log.WithError(err).Errorln("Searching twitter")
		}

		for index, tweet := range result.Statuses {
			log.WithFields(log.Fields{
				"n":        index,
				"fullText": tweet.FullText,
				"id":       tweet.ID,
			}).Println("Tweet")
		}

		if log.GetLevel() == log.DebugLevel {
			log.WithFields(log.Fields{
				"Protocol":   response.Protocol,
				"Status":     response.Status,
				"StatusCode": response.Code,
			}).Debugln("Response")

			// HTTP Headers
			for k, v := range response.Header {
				switch k {
				default:
					log.WithField(k, v).Debugln("HTTP Header")
				}
			}

			// Twitter Search API Metadata
			log.WithFields(log.Fields{
				"Count":       result.Metadata.Count,
				"SinceID":     result.Metadata.SinceID,
				"SinceIDStr":  result.Metadata.SinceIDStr,
				"MaxID":       result.Metadata.MaxID,
				"MaxIDStr":    result.Metadata.MaxIDStr,
				"RefreshURL":  result.Metadata.RefreshURL,
				"NextResults": result.Metadata.NextResults,
				"CompletedIn": result.Metadata.CompletedIn,
				"Query":       result.Metadata.Query,
			}).Debugln("Metadata")

			// Rate limits
			rl := ratelimit.New().FromHttpHeaders(response.Header)
			log.WithFields(log.Fields{
				"limit":     rl.Limit,
				"remaining": rl.Remaining,
				"reset":     rl.Reset,
			}).Debugln("Rate limit")
		}
	},
}

func init() {
	RootCmd.AddCommand(searchCmd)

	searchCmd.PersistentFlags().IntVarP(&count, "count", "", 100, "Limit number of results")
}
