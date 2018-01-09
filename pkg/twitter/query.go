package twitter

import (
	"net/url"
	"strings"
)

func SinceId(metadata *Metadata) (string, error) {
	refreshUrl := strings.TrimPrefix(metadata.RefreshURL, "?")
	m, err := url.ParseQuery(refreshUrl)
	if err != nil {
		return "", err
	}
	return m.Get("since_id"), nil
}

func MaxId(metadata *Metadata) (string, error) {
	nextResults := strings.TrimPrefix(metadata.NextResults, "?")
	m, err := url.ParseQuery(nextResults)
	if err != nil {
		return "", err
	}
	return m.Get("max_id"), nil
}
