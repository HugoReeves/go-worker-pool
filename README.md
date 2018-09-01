# Works Package

The works package aims to provide simple, prepackaged Golang concurrency patterns.
The initial pattern is inspired by this blog post from [marcio.io](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/).

# Usage

Works is easy to use.

To begin, create a queue using works.NewQueue(Number of Workers, Queue Buffer Length).
The queue contains a job queue channel and a dispatcher.
Start the dispatcher to launch the worker go routines and start working.
You can send jobs to the dispatcher using *Queue.QueueJob(Job{})

Jobs require a payload, with the method Exec.
This is how your code is able to be executed by the worker.
The payload can contain any data the exec function requires.
The exec function can return an error.

```go
type Payload struct {
	Num int
}

func (j Payload) Exec() error {
	fmt.Printf("Job No. %d\n", j.Num)
	time.Sleep(100 * time.Millisecond)
	return nil
}

func main() {
	q := works.NewQueue(2, 100)

	q.Activate()

	for i := 0; i < 10; i++ {
		q.QueueJob(Payload{i})
	}

	time.Sleep(1 * time.Second)

}
```