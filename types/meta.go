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
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"

	"github.com/nlnwa/sigridr/api"
)

type Meta struct {
	Name           string    `json:"name,omitempty"`
	Description    string    `json:"description,omitempty"`
	Created        time.Time `json:"created,omitempty"`
	CreatedBy      string    `json:"createdBy,omitempty"`
	LastModified   time.Time `json:"lastModified,omitempty"`
	LastModifiedBy string    `json:"lastModifiedBy,omitempty"`
	Label          []*Label  `json:"label,omitempty"`
}

func (m *Meta) ToProto() (*api.Meta, error) {
	label := make([]*api.Label, 0)
	for _, l := range m.Label {
		label = append(label, l.ToProto())
	}
	created, err := ptypes.TimestampProto(m.Created)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert from time to proto timestamp")
	}
	lastModified, err := ptypes.TimestampProto(m.LastModified)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert from time to proto timestamp")
	}

	return &api.Meta{
		Name:           m.Name,
		Description:    m.Description,
		Created:        created,
		CreatedBy:      m.CreatedBy,
		LastModified:   lastModified,
		LastModifiedBy: m.LastModifiedBy,
		Label:          label,
	}, nil
}

func (m *Meta) FromProto(meta *api.Meta) (*Meta, error) {
	m.Name = meta.Name
	m.Description = meta.Description
	m.CreatedBy = meta.CreatedBy
	m.LastModifiedBy = meta.LastModifiedBy

	label := make([]*Label, 0)
	for _, l := range meta.Label {
		label = append(label, new(Label).FromProto(l))
	}
	m.Label = label

	created, err := ptypes.Timestamp(meta.Created)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert from proto timestamp to time")
	}
	m.Created = created

	lastModified, err := ptypes.Timestamp(meta.LastModified)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert from proto timestamp to time")
	}
	m.LastModified = lastModified

	return m, nil
}
