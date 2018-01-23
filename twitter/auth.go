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
