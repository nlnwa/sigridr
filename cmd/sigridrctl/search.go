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

package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"

	"github.com/nlnwa/sigridr/auth"
	"github.com/nlnwa/sigridr/twitter"
	"github.com/nlnwa/sigridr/twitter/ratelimit"
)

var (
	cfgFile        string
	consumerKey    string
	consumerSecret string
	accessToken    string
	debug          bool
	count          int
)

var searchCmd = &cobra.Command{
	Use:   "search query ...",
	Short: "Query Twitter's Search API",
	Long:  `Query Twitter's Search API`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := search(cmd, args); err != nil {
			panic(err)
		}
	},
}

func search(_ *cobra.Command, args []string) error {
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
		return fmt.Errorf("failed searching twitter: %v", err)
	}

	if debug {
		fmt.Printf("Response:\n\tProtocol=%s\n\tStatus=%s\n\tCode=%d\n",
			response.Protocol,
			response.Status,
			response.Code)

		// HTTP Headers
		fmt.Println("HTTP headers:")
		for k, v := range response.Header {
			switch k {
			default:
				fmt.Printf("\t%s=%s\n", k, v)
			}
		}
		fmt.Printf("Metadata:\n\t%+v\n", result.Metadata)

		// Rate limits
		if rl, err := ratelimit.New().FromHttpHeaders(response.Header); err != nil {
			return err
		} else {
			fmt.Printf("Ratelimit:\n\tlimit=%d\n\tremaining=%d\n\treset=%s\n", rl.Limit, rl.Remaining, rl.Reset)
		}
	}

	for index, tweet := range result.Statuses {
		fmt.Printf("\n%3d: %s\n", index+1, tweet.FullText)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(searchCmd)

	cobra.OnInitialize(initConfig)

	searchCmd.Flags().IntVarP(&count, "count", "", 100, "number of results")
	searchCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.sigridr.yaml)")
	searchCmd.Flags().StringVarP(&consumerSecret, "consumer-secret", "s", "", "consumer secret")
	searchCmd.Flags().StringVarP(&consumerKey, "consumer-key", "k", "", "consumer key")
	searchCmd.Flags().StringVarP(&accessToken, "access-token", "a", "", "access token")
	searchCmd.Flags().BoolVar(&debug, "debug", false, "enable debugging output")

	viper.BindPFlag("config", searchCmd.Flags().Lookup("config"))
	viper.BindPFlag("consumer-secret", searchCmd.Flags().Lookup("consumer-secret"))
	viper.BindPFlag("consumer-key", searchCmd.Flags().Lookup("consumer-key"))
	viper.BindPFlag("access-token", searchCmd.Flags().Lookup("access-token"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".sigridr" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".sigridr")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it
	if err := viper.ReadInConfig(); err == nil {
		if debug {
			fmt.Printf("Using configuration file %s", viper.ConfigFileUsed())
			fmt.Println("Config:")
			for _, key := range viper.AllKeys() {
				fmt.Printf("\t%s=%v\n", key, viper.Get(key))
			}
		}
	} else {
		// no config file found - set default config file
		viper.SetConfigFile(home + "/.sigridr.yaml")
	}

	// If consumer key and consumer secret provided fetch oauth2 token and store it in config file
	if ck, cs := viper.GetString("consumer-key"), viper.GetString("consumer-secret"); ck != "" && cs != "" {
		token, err := twitter.Oauth2Token(ck, cs)
		if err != nil {
			panic(err)
		}
		viper.Set("token", token)
		if err := writeConfig(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		return
	}

	// If access token provided, use it and store it in config file
	if accessToken := viper.Get("access-token"); accessToken != "" {
		viper.Set("token", &oauth2.Token{AccessToken: accessToken.(string)})
		if err := writeConfig(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		return
	}
}
