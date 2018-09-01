package works

// worker represents the listening go routine that process jobs
// The worker registers its job channel into the worker pool
// This enables jobs to be sent down the job channel
// The worker waits for jobs on the job channel then executes
// and blocks until they're ready
// NOTE: None of these should be public exposed as the user
// Should never deal with workers (ie. start or stop them)

type worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	Quit       chan bool
}

// newWorker creates a new worker registeres witb the parent dispatcher's worker pool
func (d *Dispatcher) newWorker() worker {
	return worker{
		WorkerPool: d.WorkerPool,
		JobChannel: make(chan Job),
		Quit:       make(chan bool),
	}
}

// start runs a go routine forever or until a quit is signalled
// It waits for jobs to be sent down the channel and blocks while
// The payload is executed
func (w worker) start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			// Wait for a job to be available on the job channel
			// Then execute the job
			case job := <-w.JobChannel:
				err := job.Payload.Exec()

				if err != nil {
					job.Payload.OnError(err)
				}

			case <-w.Quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w worker) stop() {
	go func() {
		w.Quit <- true
	}()
}
