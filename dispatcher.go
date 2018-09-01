package works

// Dispatcher is a struct that stores workers and their job channels
type Dispatcher struct {
	JobQueue chan Job
	// A pool of workers channel that are registered with the dispatcher
	WorkerPool chan chan Job
	// The max number of workers
	MaxWorkers int
	// Allows stopping of workers
	Workers []worker
}

// Start will create workers for the dispatcher and start them up
func (d *Dispatcher) Start() {
	// starting n number of workers
	for i := 0; i < d.MaxWorkers; i++ {
		w := d.newWorker()
		w.start()
		d.Workers = append(d.Workers, w)
	}

	go d.dispatch()
}

// Stop cancels all workers and removes them from the job channel
func (d *Dispatcher) Stop() {
	for _, c := range d.Workers {
		c.stop()
		d.Workers = d.Workers[1:]
	}
}

// The private dispatch method
// Pulls jobs off of the job queue and attempts to enqueue
// them onto an available worker
func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.JobQueue:
			/* NOTE: This section used to be wrapped in an anonymous go routine.
			This was removed to prevent the draining of the buffered job queue chan
			And subsequent excesive number of go routines waiting because
			they were blocked by a lack of ready workers */

			// a job request has been received
			// try to obtain a worker job channel that is available.
			// this will block until a worker is idle
			jobChannel := <-d.WorkerPool

			// dispatch the job to the worker job channel
			jobChannel <- job
		}
	}
}
