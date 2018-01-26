package agent

import (
	"time"

	"github.com/nlnwa/sigridr/api"
	"github.com/nlnwa/sigridr/types"
)

type status struct {
	*types.Execution
}

func (s *status) setState(state api.Execution_State) *status {
	s.Execution.State = state.String()
	return s
}

func (s *status) fetching(when time.Time) *status {
	s.StartTime = when
	return s.setState(api.Execution_FETCHING)
}

func (s *status) failed(err error) *status {
	s.Execution.Error = err.Error()
	return s.setState(api.Execution_FAILED)
}

func (s *status) finished(when time.Time) *status {
	s.EndTime = when
	return s.setState(api.Execution_FINISHED)
}

func (s *status) fromExecution(execution *types.Execution) *status {
	s.Execution = execution
	return s
}

type StatusOption func(*status)

func withSeed(seed *types.Seed) StatusOption {
	return func(s *status) {
		s.SeedId = seed.Id
	}
}

func withJob(job *types.Job) StatusOption {
	return func(s *status) {
		s.JobId = job.Id
	}
}

func withDefaultState() StatusOption {
	return func(s *status) {
		s.State = api.Execution_CREATED.String()
	}
}

func newStatus(opts ...StatusOption) *status {
	s := &status{
		Execution: new(types.Execution),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}
