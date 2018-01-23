package controller

import (
	"github.com/pkg/errors"
	r "gopkg.in/gorethink/gorethink.v3"

	"github.com/nlnwa/sigridr/database"
	"github.com/nlnwa/sigridr/types"
)

type jobStore struct {
	*database.Rethink
}

func newJobStore(c Config) *jobStore {
	db := database.New(database.WithAddress(c.DatabaseHost, c.DatabasePort), database.WithName(c.DatabaseName))
	db.SetTags("json")

	return &jobStore{db}
}

func (js *jobStore) connect() error {
	return js.Rethink.Connect()
}

func (js *jobStore) disconnect() error {
	return js.Rethink.Disconnect()
}

func (js *jobStore) getJobs() ([]types.Job, error) {
	var jobs []types.Job

	if err := js.ListTable("job", &jobs); err != nil {
		return nil, err
	} else {
		return jobs, nil
	}
}

func (js *jobStore) getSeeds(job *types.Job) ([]types.Seed, error) {
	var seeds []types.Seed

	cursor, err := js.Filter("seed", func(seed r.Term) r.Term {
		return seed.Field("jobId").Contains(job.Id)
	})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(&seeds); err != nil {
		return nil, errors.Wrap(err, "failed to get all seeds from cursor")
	} else {
		return seeds, nil
	}
}
