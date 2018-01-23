package twitter

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func SinceId(metadata *Metadata) (string, error) {
	refreshUrl := strings.TrimPrefix(metadata.RefreshURL, "?")
	if m, err := url.ParseQuery(refreshUrl); err != nil {
		return "", errors.Wrap(err, "failed to parse URL querystring")
	} else {
		return m.Get("since_id"), nil
	}
}

func MaxId(metadata *Metadata) (string, error) {
	nextResults := strings.TrimPrefix(metadata.NextResults, "?")
	if m, err := url.ParseQuery(nextResults); err != nil {
		return "", errors.Wrap(err, "failed to parse URL querystring")
	} else {
		return m.Get("max_id"), nil
	}
}
