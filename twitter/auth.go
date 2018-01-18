package twitter

import (
	"context"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
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
	return config.Token(ctx)
}
