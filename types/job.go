package types

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	log "github.com/sirupsen/logrus"

	"github.com/nlnwa/sigridr/api"
)

type Job struct {
	Id             string    `json:"id,omitempty"`
	Meta           *Meta     `json:"meta,omitempty"`
	CronExpression string    `json:"cronExpression,omitempty"`
	ValidFrom      time.Time `json:"validFrom,omitempty"`
	ValidTo        time.Time `json:"validTo,omitempty"`
	Disabled       bool      `json:"disabled,omitempty"`
	Seeds          []Seed    `json:"seeds,omitempty"`
}

func (j *Job) FromProto(job *api.Job) *Job {
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

func (j *Job) ToProto() *api.Job {
	validTo, err := ptypes.TimestampProto(j.ValidTo)
	if err != nil {
		log.WithError(err).Error()
	}
	validFrom, err := ptypes.TimestampProto(j.ValidFrom)
	if err != nil {
		log.WithError(err).Error()
	}
	return &api.Job{
		Id:             j.Id,
		Meta:           j.Meta.ToProto(),
		CronExpression: j.CronExpression,
		ValidTo:        validTo,
		ValidFrom:      validFrom,
		Disabled:       j.Disabled,
	}

}

func (j *Job) IsValid() bool {
	now := time.Now().UTC()
	return now.After(j.ValidFrom) && now.Before(j.ValidTo)
}
