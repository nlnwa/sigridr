package types

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	log "github.com/sirupsen/logrus"

	"github.com/nlnwa/sigridr/api/sigridr"
)

type Label struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func (l *Label) ToProto() *sigridr.Label {
	label := sigridr.Label(*l)

	return &label
}

func (l *Label) FromProto(label *sigridr.Label) *Label {
	m := Label(*label)
	l = &m

	return l
}

type Meta struct {
	Name           string    `json:"name,omitempty"`
	Description    string    `json:"description,omitempty"`
	Created        time.Time `json:"created,omitempty"`
	CreatedBy      string    `json:"createdBy,omitempty"`
	LastModified   time.Time `json:"lastModified,omitempty"`
	LastModifiedBy string    `json:"lastModifiedBy,omitempty"`
	Label          []*Label  `json:"label,omitempty"`
}

func (s *Meta) ToProto() *sigridr.Meta {
	label := make([]*sigridr.Label, 0)
	for _, l := range s.Label {
		label = append(label, l.ToProto())
	}
	created, err := ptypes.TimestampProto(s.Created)
	if err != nil {
		log.WithError(err).Error()
	}
	lastModified, err := ptypes.TimestampProto(s.LastModified)
	if err != nil {
		log.WithError(err).Error()
	}

	return &sigridr.Meta{
		Name:           s.Name,
		Description:    s.Description,
		Created:        created,
		CreatedBy:      s.CreatedBy,
		LastModified:   lastModified,
		LastModifiedBy: s.LastModifiedBy,
		Label:          label,
	}
}

func (m *Meta) FromProto(meta *sigridr.Meta) *Meta {
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

type Seed struct {
	Id       string   `json:"id,omitempty"`
	Meta     *Meta    `json:"meta,omitempty"`
	EntityId string   `json:"entityId,omitempty"`
	JobId    []string `json:"jobId,omitempty"`
	Disabled bool     `json:"disabled,omitempty"`
}

func (s *Seed) ToProto() *sigridr.Seed {
	return &sigridr.Seed{
		Id:       s.Id,
		Meta:     s.Meta.ToProto(),
		EntityId: s.EntityId,
		JobId:    s.JobId,
		Disabled: s.Disabled,
	}
}

func (s *Seed) FromProto(seed *sigridr.Seed) *Seed {
	s.Id = seed.Id
	s.Meta = new(Meta).FromProto(seed.Meta)
	s.EntityId = seed.EntityId
	s.JobId = seed.JobId
	s.Disabled = seed.Disabled

	return s
}

type Entity struct {
	Id   string `json:"id,omitempty"`
	Meta *Meta  `json:"meta,omitempty"`
}

type Job struct {
	Id             string    `json:"id,omitempty"`
	Meta           *Meta     `json:"meta,omitempty"`
	CronExpression string    `json:"cronExpression,omitempty"`
	ValidFrom      time.Time `json:"validFrom,omitempty"`
	ValidTo        time.Time `json:"validTo,omitempty"`
	Disabled       bool      `json:"disabled,omitempty"`
	Seeds          []Seed    `json:"seeds,omitempty"`
}

func (j *Job) FromProto(job *sigridr.Job) *Job {
	validTo, err := ptypes.Timestamp(job.ValidTo)
	if err != nil {
		log.WithError(err).Error()
	}
	validFrom, err := ptypes.Timestamp(job.ValidFrom)
	if err != nil {
		log.WithError(err).Error()
	}
	j.Id = job.Id
	j.Meta = new(Meta).FromProto(job.Meta)
	j.CronExpression = job.CronExpression
	j.ValidTo = validTo
	j.ValidFrom = validFrom
	j.Disabled = job.Disabled

	return j
}

func (job *Job) ToProto() *sigridr.Job {
	validTo, err := ptypes.TimestampProto(job.ValidTo)
	if err != nil {
		log.WithError(err).Error()
	}
	validFrom, err := ptypes.TimestampProto(job.ValidFrom)
	if err != nil {
		log.WithError(err).Error()
	}
	return &sigridr.Job{
		Id:             job.Id,
		Meta:           job.Meta.ToProto(),
		CronExpression: job.CronExpression,
		ValidTo:        validTo,
		ValidFrom:      validFrom,
		Disabled:       job.Disabled,
	}

}

func (j *Job) IsValid() bool {
	now := time.Now().UTC();
	isValid := now.After(j.ValidFrom) && now.Before(j.ValidTo)
	if !isValid {
		log.WithFields(log.Fields{
			"validFrom": j.ValidFrom,
			"validTo":   j.ValidTo,
			"now":       now,
		}).Debugln("Job not valid")
	}
	return isValid
}
