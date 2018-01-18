package main

import (
	"os"

	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type config struct {
	consumerKey    string
	consumerSecret string
	Token          *oauth2.Token `json:"token"`
}

// Write config file
//
// Writes an oauth2 access token to config file
func writeConfig() {
	token := viper.Get("token").(*oauth2.Token)
	c := config{Token: token}

	y, err := yaml.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(viper.ConfigFileUsed())
	if err != nil {
		log.WithError(err).Fatal()
	}
	defer f.Close()
	f.Chmod(0600)
	f.Write(y)
}
