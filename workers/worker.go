package workers

import (
	"log"
	"sync"
)

const (
	MAP_WORKER = "map_worker"
)

type Worker interface {
	ID() string
	AddJob(Job) error
	Do() error
	Stop()
}

type WorkerPool struct {
	mu      sync.Mutex
	max     int
	workers map[string]Worker
	working map[string]bool
}

func NewWorkerPool(max int) *WorkerPool {
	log.Println("creating worker pool with max workers: ", max)
	return &WorkerPool{
		max:     max,
		workers: make(map[string]Worker),
		working: make(map[string]bool),
	}
}

func (p *WorkerPool) Add(w Worker) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.workers) >= p.max {
		log.Printf("worker pool is full, cannot add worker: %s", w.ID())
		return
	}
	log.Printf("adding worker: %s", w.ID())
	p.workers[w.ID()] = w

	log.Println("workers added to pool: ", len(p.workers))
}

func (p *WorkerPool) Remove(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.workers, id)
}

func (p *WorkerPool) Get(id string) Worker {
	p.mu.Lock()
	copy := p.workers
	p.mu.Unlock()
	for _, w := range copy {
		if w.ID() == id {
			return w
		}
	}
	return nil
}

func (p *WorkerPool) StartAll() {
	log.Printf("starting workers")
	for _, w := range p.workers {
		log.Println("starting worker: ", w.ID())
		go func(p *WorkerPool, w Worker) {
			err := w.Do()
			if err != nil {
				log.Println("error starting worker: ", w.ID(), " error: ", err)
				panic(err)
			}
			p.mu.Lock()
			p.working[w.ID()] = true
			p.mu.Unlock()
		}(p, w)
	}
}

func (p *WorkerPool) StopAll() {
	p.mu.Lock()
	defer p.mu.Unlock()

	log.Printf("stopping workers")
	for _, w := range p.workers {
		if p.working[w.ID()] {
			log.Println("stopping worker: ", w.ID())
			w.Stop()
			delete(p.working, w.ID())
		}
	}
}

func (p *WorkerPool) StopWorker(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	log.Println("stopping worker: ", id)
	w := p.workers[id]
	if w != nil {
		w.Stop()
		delete(p.working, id)
	}
}
