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
	"context"
	"strings"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"

	"github.com/nlnwa/sigridr/cmd/sigridrctl"
	"github.com/nlnwa/sigridr/pkg/twitter/auth"
)

var (
	cfgFile        string
	consumerKey    string
	consumerSecret string
	accessToken    string
	debug          bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "sigridr",
	Short: "Twitter API client",
	Long:  `Twitter API client`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.WithError(err).Fatal()
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sigridr.yaml)")

	RootCmd.PersistentFlags().StringVarP(&consumerSecret, "consumer-secret", "s", "", "Consumer secret")
	RootCmd.PersistentFlags().StringVarP(&consumerKey, "consumer-key", "k", "", "Consumer key")
	RootCmd.PersistentFlags().StringVarP(&accessToken, "access-token", "a", "", "Access token")
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "Turn on debugging")

	viper.BindPFlag("consumer-secret", RootCmd.PersistentFlags().Lookup("consumer-secret"))
	viper.BindPFlag("consumer-key", RootCmd.PersistentFlags().Lookup("consumer-key"))
	viper.BindPFlag("access-token", RootCmd.PersistentFlags().Lookup("access-token"))
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		log.WithError(err).Fatal()
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".sigridr" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".sigridr")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it
	if err := viper.ReadInConfig(); err == nil {
		log.WithFields(log.Fields{"path": viper.ConfigFileUsed()}).Debugln("Configuration file used")

		for _, key := range viper.AllKeys() {
			log.WithField(key, viper.Get(key)).Debugln("Configuration value")
		}
	} else {
		// no config file found - set default config file
		viper.SetConfigFile(home + "/.sigridr.yaml")
	}

	// If consumer key and consumer secret provided fetch oauth2 token and store it in config file
	if ck, cs := viper.GetString("consumer-key"), viper.GetString("consumer-secret"); ck != "" && cs != "" {
		token, err := twitter.Oauth2Token(ck, cs)
		if err != nil {
			log.WithError(err).Fatal()
		}
		viper.Set("token", token)
		sigridrctl.WriteConfig()
	}

	// If access token provided, use it and store it in config file
	if accessToken := viper.Get("access-token"); accessToken != "" {
		viper.Set("token", &oauth2.Token{AccessToken: accessToken.(string)})
		sigridrctl.WriteConfig()
	}
}
