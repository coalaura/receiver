package main

import "sync"

type Worker struct {
	wg   sync.WaitGroup
	jobs chan func()
}

func NewWorker(routines int) *Worker {
	wk := &Worker{
		jobs: make(chan func(), 64),
	}

	for range routines {
		wk.wg.Add(1)

		go wk.Work()
	}

	return wk
}

func (w *Worker) Schedule(fn func()) {
	w.jobs <- fn
}

func (w *Worker) Work() {
	defer w.wg.Done()

	for fn := range w.jobs {
		fn()
	}
}

func (w *Worker) Stop() {
	close(w.jobs)

	w.wg.Wait()
}
