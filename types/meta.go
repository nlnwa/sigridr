package types

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	log "github.com/sirupsen/logrus"

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

func (m *Meta) ToProto() *api.Meta {
	label := make([]*api.Label, 0)
	for _, l := range m.Label {
		label = append(label, l.ToProto())
	}
	created, err := ptypes.TimestampProto(m.Created)
	if err != nil {
		log.WithError(err).Error()
	}
	lastModified, err := ptypes.TimestampProto(m.LastModified)
	if err != nil {
		log.WithError(err).Error()
	}

	return &api.Meta{
		Name:           m.Name,
		Description:    m.Description,
		Created:        created,
		CreatedBy:      m.CreatedBy,
		LastModified:   lastModified,
		LastModifiedBy: m.LastModifiedBy,
		Label:          label,
	}
}

func (m *Meta) FromProto(meta *api.Meta) *Meta {
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
		log.WithError(err).Error()
	}
	m.Created = created

	lastModified, err := ptypes.Timestamp(meta.LastModified)
	if err != nil {
		log.WithError(err).Error()
	}
	m.LastModified = lastModified

	return m
}
