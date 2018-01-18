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
		oauth2Token := &oauth2.Token{}
		t := token.(map[string]interface{})

		accessToken, ok := t["access_token"]
		if ok {
			oauth2Token.AccessToken = accessToken.(string)
		}
		expiry, ok := t["expiry"]
		if ok {
			expiryAsTime, _ := time.Parse(time.RFC3339, expiry.(string))
			oauth2Token.Expiry = expiryAsTime
		}
		tokenType, ok := t["token_type"]
		if ok {
			oauth2Token.TokenType = tokenType.(string)
		}
		return oauth2Token
	case string:
		return &oauth2.Token{AccessToken: token.(string)}
	default:
		return nil
	}
}
