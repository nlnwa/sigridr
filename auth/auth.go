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

package auth

import (
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

const TWITTER_OAUTH2_TOKEN_URL = "https://api.twitter.com/oauth2/token"

func GetTwitterOauth2Token(key string, secret string) (*oauth2.Token, error) {
	config := &clientcredentials.Config{
		ClientID:     key,
		ClientSecret: secret,
		TokenURL:     TWITTER_OAUTH2_TOKEN_URL,
	}
	return config.Token(context.TODO())
}

func HttpClient(token interface{}) *http.Client {
	return httpClient(convertToOauth2Token(token))
}

func httpClient(token *oauth2.Token) *http.Client {
	config := &oauth2.Config{}
	// token := &oauth2.Token{AccessToken: accessToken}
	return config.Client(context.TODO(), token)
}

func convertToOauth2Token(token interface{}) *oauth2.Token {
	switch token.(type) {
	case *oauth2.Token:
		return token.(*oauth2.Token)
	default:
		t := token.(map[string]interface{})
		accessToken := t["access_token"].(string)
		expiry, _ := time.Parse(time.RFC3339, t["expiry"].(string))
		tokenType := t["token_type"].(string)
		return &oauth2.Token{
			AccessToken: accessToken,
			Expiry:      expiry,
			TokenType:   tokenType,
		}
	}
}
