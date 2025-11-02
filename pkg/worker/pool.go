package worker

import (
	"context"
	"sync"
)

// Job represents a unit of work
type Job struct {
	ID      string
	Payload interface{}
}

// Result represents the result of a job
type Result struct {
	Job   Job
	Value interface{}
	Err   error
}

// Worker function type
type WorkerFunc func(ctx context.Context, job Job) (interface{}, error)

// Pool represents a worker pool
type Pool struct {
	workers    int
	jobs       chan Job
	results    chan Result
	workerFunc WorkerFunc
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewPool creates a new worker pool
func NewPool(workers int, workerFunc WorkerFunc) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	return &Pool{
		workers:    workers,
		jobs:       make(chan Job, workers*2),
		results:    make(chan Result, workers*2),
		workerFunc: workerFunc,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start starts the worker pool
func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

// worker is the worker goroutine
func (p *Pool) worker(id int) {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case job, ok := <-p.jobs:
			if !ok {
				return
			}

			value, err := p.workerFunc(p.ctx, job)
			p.results <- Result{
				Job:   job,
				Value: value,
				Err:   err,
			}
		}
	}
}

// Submit submits a job to the pool
func (p *Pool) Submit(job Job) {
	p.jobs <- job
}

// Results returns the results channel
func (p *Pool) Results() <-chan Result {
	return p.results
}

// Stop stops the worker pool
func (p *Pool) Stop() {
	close(p.jobs)
	p.wg.Wait()
	close(p.results)
}

// StopWithCancel stops the pool and cancels all running jobs
func (p *Pool) StopWithCancel() {
	p.cancel()
	p.Stop()
}
