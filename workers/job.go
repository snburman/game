package workers

import "github.com/google/uuid"

type Job struct {
	id   string
	data any
}

func NewJob[T any](data T) Job {
	return Job{
		id:   "job_" + uuid.New().String(),
		data: data,
	}
}

func (j Job) ID() string {
	return j.id
}

func (j Job) Run() error {
	return nil
}

func (j Job) Data() any {
	return j.data
}
