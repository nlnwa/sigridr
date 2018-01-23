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

package twitter

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const Oauth2TokenUrl = "https://api.twitter.com/oauth2/token"

func Oauth2Token(key string, secret string) (*oauth2.Token, error) {
	tokenAcquisitionClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, tokenAcquisitionClient)

	config := &clientcredentials.Config{
		ClientID:     key,
		ClientSecret: secret,
		TokenURL:     Oauth2TokenUrl,
	}
	if token, err := config.Token(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to fetch oauth2 token")
	} else {
		return token, nil
	}
}
