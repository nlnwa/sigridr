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
	"os"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
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
func writeConfig() error {
	token := viper.Get("token").(*oauth2.Token)
	c := config{Token: token}

	y, err := yaml.Marshal(c)
	if err != nil {
		errors.Wrap(err, "failed to marshal yaml")
	}

	f, err := os.Create(viper.ConfigFileUsed())
	if err != nil {
		return errors.Wrap(err, "failed to create config file")
	}
	defer f.Close()

	f.Chmod(0600)
	_, err = f.Write(y)
	if err != nil {
		return errors.Wrap(err, "failed to write config file")
	}
	return nil
}
