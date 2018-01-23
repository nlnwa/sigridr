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

package types

import (
	"github.com/nlnwa/sigridr/api"
)

type Seed struct {
	Id       string   `json:"id,omitempty"`
	Meta     *Meta    `json:"meta,omitempty"`
	EntityId string   `json:"entityId,omitempty"`
	JobId    []string `json:"jobId,omitempty"`
	Disabled bool     `json:"disabled,omitempty"`
}

func (s *Seed) ToProto() (*api.Seed, error) {
	meta, err := s.Meta.ToProto()
	if err != nil {
		return nil, err
	}
	return &api.Seed{
		Id:       s.Id,
		Meta:     meta,
		EntityId: s.EntityId,
		JobId:    s.JobId,
		Disabled: s.Disabled,
	}, nil
}

func (s *Seed) FromProto(seed *api.Seed) (*Seed, error) {
	meta, err := new(Meta).FromProto(seed.Meta)
	if err != nil {
		return nil, err
	}
	s.Id = seed.Id
	s.Meta = meta
	s.EntityId = seed.EntityId
	s.JobId = seed.JobId
	s.Disabled = seed.Disabled

	return s, nil
}
