package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
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
			log.WithError(err).Error()
			os.Exit(2)
		}
	},
}

func search(cmd *cobra.Command, args []string) error {
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

	for index, tweet := range result.Statuses {
		log.WithFields(log.Fields{
			"n":        index,
			"fullText": tweet.FullText,
			"id":       tweet.ID,
		}).Println("Status")
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

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it
	if err := viper.ReadInConfig(); err == nil {
		log.WithFields(log.Fields{"path": viper.ConfigFileUsed()}).Debugln("using configuration file")

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
		writeConfig()
		return
	}

	// If access token provided, use it and store it in config file
	if accessToken := viper.Get("access-token"); accessToken != "" {
		viper.Set("token", &oauth2.Token{AccessToken: accessToken.(string)})
		writeConfig()
		return
	}
}
