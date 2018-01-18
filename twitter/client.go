package twitter

import (
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
)

type Client interface {
	Search(*Params) (*Result, *Response, error)
}

// Wrap go-twitter client
type lib struct {
	client *twitter.Client
}

// NewClient creates a new Client using the provided httpClient
func New(httpClient *http.Client) Client {
	return &lib{twitter.NewClient(httpClient)}
}

// Search searches using Twitter's Search API
func (l *lib) Search(params *Params) (*Result, *Response, error) {
	p := twitter.SearchTweetParams(*params)

	if search, response, err := l.client.Search.Tweets(&p); err != nil {
		return nil, nil, err
	} else {
		return search, new(Response).FromHttpResponse(response), nil
	}
}
