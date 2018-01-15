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
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

func HttpClient(token interface{}) *http.Client {
	return httpClient(unmarshalOauth2Token(token))
}

func httpClient(token *oauth2.Token) *http.Client {
	return new(oauth2.Config).Client(context.Background(), token)
}

// UnmarshalOauth2Token returns an oauth2 token from a source that can be either
// an oauth2 token, a map with access token in "access_token" key, or from a string that is the access_token
func unmarshalOauth2Token(token interface{}) *oauth2.Token {
	switch token.(type) {
	case *oauth2.Token:
		return token.(*oauth2.Token)
	case map[string]interface{}:
		ok := false
		oauth2Token := &oauth2.Token{}
		t := token.(map[string]interface{})

		accessToken, ok := t["access_token"].(string)
		if ok {
			oauth2Token.AccessToken = accessToken
		}
		expiry, ok := t["expiry"].(string)
		if ok {
			expiryAsTime, _ := time.Parse(time.RFC3339, expiry)
			oauth2Token.Expiry = expiryAsTime
		}
		tokenType, ok := t["token_type"].(string)
		if ok {
			oauth2Token.TokenType = tokenType
		}
		return oauth2Token
	case string:
		return &oauth2.Token{AccessToken: token.(string)}
	default:
		return nil
	}
}
