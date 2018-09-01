package works

// Queue is the top level struct
// It provides access to both the Job Queue and
// Stores the child dispatcher wich pulls jobs off the queue
type Queue struct {
	JobQueue   chan Job
	Dispatcher *Dispatcher
}

// NewQueue creates a new queue object
// With a max number of workers and a max queue buffer
func NewQueue(maxWorkers int, maxQueue int) *Queue {
	q := make(chan Job, maxQueue)
	return &Queue{
		q,
		&Dispatcher{
			JobQueue:   q,
			WorkerPool: make(chan chan Job, maxWorkers),
			MaxWorkers: maxWorkers,
		},
	}
}

// QueueJob adds a job onto the end of the job queue
// This is blocking and can take time
func (q *Queue) QueueJob(p Payload) {
	q.JobQueue <- Job{Payload: p}
}

// Activate will start the queue dispatcher
func (q *Queue) Activate() {
	q.Dispatcher.Start()
}
