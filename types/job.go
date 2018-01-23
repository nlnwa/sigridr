package types

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"

	"github.com/nlnwa/sigridr/api"
)

type Job struct {
	Id             string    `json:"id,omitempty"`
	Meta           *Meta     `json:"meta,omitempty"`
	CronExpression string    `json:"cronExpression,omitempty"`
	ValidFrom      time.Time `json:"validFrom,omitempty"`
	ValidTo        time.Time `json:"validTo,omitempty"`
	Disabled       bool      `json:"disabled,omitempty"`
}

func (j *Job) FromProto(job *api.Job) (*Job, error) {
	validTo, err := ptypes.Timestamp(job.ValidTo)
	if err != nil {
		return nil, err
	}
	validFrom, err := ptypes.Timestamp(job.ValidFrom)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert from proto timestamp to time")
	}
	j.Id = job.Id
	if j.Meta, err = new(Meta).FromProto(job.Meta); err != nil {
		return nil, errors.Wrap(err, "failed to convert from proto timestamp to time")
	}
	j.CronExpression = job.CronExpression
	j.ValidTo = validTo
	j.ValidFrom = validFrom
	j.Disabled = job.Disabled

	return j, nil
}

func (j *Job) ToProto() (*api.Job, error) {
	validTo, err := ptypes.TimestampProto(j.ValidTo)
	if err != nil {
		return nil, err
	}
	validFrom, err := ptypes.TimestampProto(j.ValidFrom)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert from time to proto timestamp")
	}
	meta, err := j.Meta.ToProto()
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert from time to proto timestamp")
	}
	return &api.Job{
		Id:             j.Id,
		Meta:           meta,
		CronExpression: j.CronExpression,
		ValidTo:        validTo,
		ValidFrom:      validFrom,
		Disabled:       j.Disabled,
	}, nil

}

func (j *Job) IsValid() bool {
	now := time.Now().UTC()
	return now.After(j.ValidFrom) && now.Before(j.ValidTo)
}
