package workers

import (
	"log"
	"reflect"

	"github.com/snburman/magicgame/objects"
)

type MapWorker struct {
	id   string
	jobs chan Job
	done chan bool
	err  chan error
}

func NewMapWorker(objects *objects.ObjectManager) MapWorker {
	return MapWorker{
		id:   MAP_WORKER,
		jobs: make(chan Job, 1028),
		done: make(chan bool),
		err:  make(chan error),
	}
}

func (w MapWorker) ID() string {
	return w.id
}

func (w MapWorker) AddJob(j Job) error {
	w.jobs <- j
	return nil
}

func (w MapWorker) Do() error {
	log.Println("starting ", w.id)
	go func(w MapWorker) {
		for {
			select {
			case job := <-w.jobs:
				w.Process(job)
			case <-w.done:
				return
			}
		}
	}(w)
	return nil
}

func (w MapWorker) Stop() {
	w.done <- true
}

func (w MapWorker) Process(j Job) {
	switch j.data.(type) {
	case Draw:
		err := j.Data().(Draw).Run()
		if err != nil {
			log.Printf("error running job: %s", err)
			w.err <- err
		}
	default:
		t := reflect.TypeOf(j.Data())
		log.Printf("unsupported job type: %s", t)
	}
}
