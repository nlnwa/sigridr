// Copyright 2018 National Library of Norway
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

package worker

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/nlnwa/sigridr/api"
)

type ApiClient struct {
	address      string
	cc           *grpc.ClientConn
	workerClient api.WorkerClient
}

func NewApiClient(address string) *ApiClient {
	return &ApiClient{address: address}
}

func (c *ApiClient) Dial() (err error) {
	if c.cc, err = grpc.Dial(c.address, grpc.WithInsecure()); err != nil {
		return errors.Wrapf(err, "failed to dial: %s", c.address)
	} else {
		c.workerClient = api.NewWorkerClient(c.cc)
		return
	}
}

func (c *ApiClient) Hangup() error {
	if c.cc != nil {
		return c.cc.Close()
	} else {
		return nil
	}
}

func (c *ApiClient) Do(ctx context.Context, queuedSeed *api.QueuedSeed) (*api.WorkReply, error) {
	work := &api.WorkRequest{QueuedSeed: queuedSeed}

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	if reply, err := c.workerClient.Do(ctx, work); err != nil {
		return nil, err
	} else {
		return reply, nil
	}
}
