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

package agent

import (
	"context"

	pb "github.com/golang/protobuf/ptypes/empty"

	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/types"
)

type agentApi struct {
	store *agentStore
}

func NewApi(c Config) api.AgentServer {
	return &agentApi{
		store: newStore(c),
	}
}

func (a *agentApi) Do(ctx context.Context, req *api.DoJobRequest) (*pb.Empty, error) {
	seed, err := new(types.Seed).FromProto(req.Seed)
	if err != nil {
		return nil, err
	}
	if seed.Meta.Name == "" {
		return new(pb.Empty), nil
	}
	job, err := new(types.Job).FromProto(req.Job)
	if err != nil {
		return nil, err
	}

	if err := a.store.connect(); err != nil {
		return nil, err
	}

	s := newStatus(
		withJob(job),
		withSeed(seed),
		withDefaultState())

	id, err := a.store.saveStatus(s)
	if err != nil {
		return nil, err
	}

	queuedSeed := &api.QueuedSeed{
		ExecutionId: id,
		SeedId:      seed.Id,
		Parameter:   &api.Parameter{Query: seed.Meta.Name},
	}
	if err := a.store.enqueueSeed(queuedSeed); err != nil {
		return nil, err
	}

	return new(pb.Empty), a.store.disconnect()
}
