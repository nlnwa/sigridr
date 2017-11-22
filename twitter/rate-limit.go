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

package twitter

import (
	"time"
	"net/http"
	"strings"
	"log"
	"strconv"
)

const headerPrefix = "X-Rate-Limit-"

type RateLimit struct {
	Limit int
	Remaining int
	Reset time.Time
}

// NewRateLimit creates an instance of RateLimit based on HTTP Headers
func NewRateLimit(header *http.Header) *RateLimit {
	rateLimit := new(RateLimit)
	for key, value := range *header {
		if strings.HasPrefix(key, headerPrefix) {
			n, err := strconv.ParseInt(value[0], 10, 32)
			if err != nil {
				log.Println(err)
			}
			switch strings.TrimPrefix(key, headerPrefix) {
			case "Reset":
				rateLimit.Reset = time.Unix(n, 0)
			case "Limit":
				rateLimit.Limit = int(n)
			case "Remaining":
				rateLimit.Remaining = int(n)
			}
		}
	}
	return rateLimit
}
