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

func (s *Seed) ToProto() *api.Seed {
	return &api.Seed{
		Id:       s.Id,
		Meta:     s.Meta.ToProto(),
		EntityId: s.EntityId,
		JobId:    s.JobId,
		Disabled: s.Disabled,
	}
}

func (s *Seed) FromProto(seed *api.Seed) *Seed {
	s.Id = seed.Id
	s.Meta = new(Meta).FromProto(seed.Meta)
	s.EntityId = seed.EntityId
	s.JobId = seed.JobId
	s.Disabled = seed.Disabled

	return s
}
