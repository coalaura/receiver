package main

type Worker struct {
	jobs chan func()
}

func NewWorker(routines int) *Worker {
	wk := &Worker{
		jobs: make(chan func(), 64),
	}

	for range routines {
		go wk.Work()
	}

	return wk
}

func (w *Worker) Schedule(fn func()) {
	w.jobs <- fn
}

func (w *Worker) Work() {
	for {
		fn := <-w.jobs

		fn()
	}
}
