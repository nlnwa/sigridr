package types

import (
	"time"

	"github.com/nlnwa/sigridr/api/sigridr"
	"github.com/golang/protobuf/ptypes"
	log "github.com/sirupsen/logrus"
)

type Job struct {
	Id             string        `json:"id"`
	Meta           *sigridr.Meta `json:"meta,omitempty"`
	CronExpression string        `json:"cron_expression"`
	ValidFrom      time.Time     `json:"valid_from"`
	ValidTo        time.Time     `json:"valid_to"`
	Disabled       bool          `json:"disabled"`
	Seeds          []Seed        `json:"seeds,omitempty"`
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
	j.Meta = job.Meta
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
		Meta:           job.Meta,
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
