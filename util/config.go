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

package util

import (
	"os"

	"golang.org/x/oauth2"
	"github.com/spf13/viper"
	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
)

type config struct {
	consumerKey    string       `json:"consumer-key"`
	consumerSecret string       `json:"consumer-secret"`
	Token          oauth2.Token `json:"token"`
}

func WriteConfig() {
	token := viper.Get("token").(*oauth2.Token)
	c := config{Token: *token}

	y, err := yaml.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(viper.ConfigFileUsed())
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.Chmod(0600)
	f.Write(y)
}
