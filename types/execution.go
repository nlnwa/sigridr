package types

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/nlnwa/sigridr/api"
	"github.com/pkg/errors"
)

type Execution struct {
	Id        string    `json:"id,omitempty"`
	State     string    `json:"state,omitempty"`
	JobId     string    `json:"job_id,omitempty"`
	SeedId    string    `json:"seed_id,omitempty"`
	StartTime time.Time `json:"start_time,omitempty"`
	EndTime   time.Time `json:"end_time,omitempty"`
	Statuses  int32     `json:"statuses,omitempty"`
	Error     string    `json:"error,omitempty"`
}

func (e *Execution) FromProto(exec *api.Execution) (*Execution, error) {
	startTime, err := ptypes.Timestamp(exec.GetStartTime())
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert startTime from proto timestamp to time")
	}
	endTime, err := ptypes.Timestamp(exec.GetEndTime())
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert endTime from proto timestamp to time")
	}
	e.Id = exec.Id
	e.State = exec.GetState().String()
	e.JobId = exec.GetJobId()
	e.SeedId = exec.GetSeedId()
	e.StartTime = startTime
	e.EndTime = endTime
	e.Statuses = exec.GetStatuses()
	e.Error = exec.GetError()

	return e, nil
}

func (e *Execution) ToProto() (*api.Execution, error) {
	startTime, err := ptypes.TimestampProto(e.StartTime)
	if err != nil {
		return nil, err
	}
	endTime, err := ptypes.TimestampProto(e.EndTime)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert from time to proto timestamp")
	}
	state, ok := api.Execution_State_value[e.State]
	if !ok {
		state = int32(api.Execution_UNDEFINED)
	}
	return &api.Execution{
		Id:        e.Id,
		SeedId:    e.SeedId,
		JobId:     e.JobId,
		StartTime: startTime,
		EndTime:   endTime,
		Statuses:  e.Statuses,
		State:     api.Execution_State(state),
		Error:     e.Error,
	}, nil

}
