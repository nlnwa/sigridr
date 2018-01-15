package twitter

import (
	"context"

	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

const TWITTER_OAUTH2_TOKEN_URL = "https://api.twitter.com/oauth2/token"

func Oauth2Token(key string, secret string) (*oauth2.Token, error) {
	tokenAcquisitionClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, tokenAcquisitionClient)

	config := &clientcredentials.Config{
		ClientID:     key,
		ClientSecret: secret,
		TokenURL:     TWITTER_OAUTH2_TOKEN_URL,
	}
	return config.Token(ctx)
}
