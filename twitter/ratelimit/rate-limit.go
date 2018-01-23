package ratelimit

import (
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"

	"github.com/nlnwa/sigridr/api"
)

const headerPrefix = "X-Rate-Limit-"

type RateLimit struct {
	Limit     int
	Remaining int
	Reset     time.Time
}

func New() *RateLimit {
	return &RateLimit{Limit: 450, Remaining: 450}
}

func (rl *RateLimit) ToProto() (*api.RateLimit, error) {
	reset, err := ptypes.TimestampProto(rl.Reset)
	if err != nil {
		return nil, err
	}
	return &api.RateLimit{
		Limit:     int32(rl.Limit),
		Remaining: int32(rl.Remaining),
		Reset_:    reset,
	}, nil
}

func (rl *RateLimit) FromProto(rateLimit *api.RateLimit) (*RateLimit, error) {
	reset, err := ptypes.Timestamp(rateLimit.GetReset_())
	if err != nil {
		return nil, errors.Wrap(err, "failed creating ratelimit from protobuf struct")
	}
	rl.Limit = int(rateLimit.GetLimit())
	rl.Remaining = int(rateLimit.GetRemaining())
	rl.Reset = reset
	return rl, nil
}

func (rl *RateLimit) Timeout() time.Duration {
	return rl.Reset.Sub(time.Now().UTC())
}

func (rl *RateLimit) WithReset(reset time.Time) *RateLimit {
	rl.Reset = reset
	return rl
}

// NewRateLimit creates an instance of RateLimit based on HTTP Headers
func (rl *RateLimit) FromHttpHeaders(header map[string][]string) (*RateLimit, error) {
	for key, value := range header {
		if strings.HasPrefix(key, headerPrefix) {
			n, err := strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "failed creating ratelimit from http headers")
			}
			switch strings.TrimPrefix(key, headerPrefix) {
			case "Reset":
				rl.Reset = time.Unix(n, 0).UTC()
			case "Limit":
				rl.Limit = int(n)
			case "Remaining":
				rl.Remaining = int(n)
			}
		}
	}
	return rl, nil
}
