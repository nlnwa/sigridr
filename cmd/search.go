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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/nlnwa/sigridr/auth"
	"github.com/nlnwa/sigridr/twitter"
	"github.com/nlnwa/sigridr/db"
	"fmt"
)

var save bool

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search query ...",
	Short: "Query Twitter's Search API",
	Long:  `Query Twitter's Search API`,
	Run: func(cmd *cobra.Command, args []string) {
		query := strings.Join(args, " ")
		params := &twitter.SearchParams{Query: query}

		httpClient := auth.HttpClient(viper.Get("token"))
		client := twitter.NewClient(httpClient)
		tweets := client.Search(params)

		if save {
			db.Connect(db.Options{Database: "sigridr"})
			defer db.Disconnect()
			db.CreateTable("result")
			db.Insert("result", tweets)
		} else {
			for index, tweet := range tweets {
				fmt.Printf("[%v] : %v\n", index, tweet.FullText)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(searchCmd)

	searchCmd.PersistentFlags().BoolVarP(&save, "save", "", false, "Save result")
}
