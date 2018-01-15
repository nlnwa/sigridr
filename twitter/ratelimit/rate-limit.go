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

package ratelimit

import (
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	log "github.com/sirupsen/logrus"

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

func (rl *RateLimit) ToProto() *api.RateLimit {
	reset, err := ptypes.TimestampProto(rl.Reset)
	if err != nil {
		log.WithError(err).Error()
	}
	return &api.RateLimit{int32(rl.Limit), int32(rl.Remaining), reset}
}

func (rl *RateLimit) FromProto(rateLimit *api.RateLimit) *RateLimit {
	reset, err := ptypes.Timestamp(rateLimit.GetReset_())
	if err != nil {
		log.WithError(err).Error()
	}
	rl.Limit = int(rateLimit.GetLimit())
	rl.Remaining = int(rateLimit.GetRemaining())
	rl.Reset = reset
	return rl
}

func (rl *RateLimit) Timeout() time.Duration {
	return rl.Reset.Sub(time.Now().UTC())
}

func (rl *RateLimit) WithReset(reset time.Time) *RateLimit {
	rl.Reset = reset
	return rl
}

// NewRateLimit creates an instance of RateLimit based on HTTP Headers
func (rl *RateLimit) FromHttpHeaders(header map[string][]string) *RateLimit {
	for key, value := range header {
		if strings.HasPrefix(key, headerPrefix) {
			n, err := strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				log.Println(err)
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
	return rl
}
