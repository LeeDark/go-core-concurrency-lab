package workerpool

import "sync"

// Job is one unit of work accepted by the pool.
type Job struct {
	ID      int
	Payload string
}

// Result is the value produced after processing a Job.
type Result struct {
	JobID int
	Value string
	Err   error
}

// Run starts a fixed number of workers and returns a receive-only results channel.
// The caller owns sending to and closing jobs and must keep receiving results until
// the returned channel is closed. A non-positive workerCount is normalized to one.
// The handler must be non-nil and must not panic; Run does not recover handler panics.
// A nil jobs channel makes workers wait indefinitely because v1 has no cancellation.
func Run(workerCount int, jobs <-chan Job, handle func(Job) Result) <-chan Result {
	if workerCount <= 0 {
		workerCount = 1
	}

	results := make(chan Result)

	var wg sync.WaitGroup
	wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		// workers
		go func() {
			defer wg.Done()

			// workers read from jobs
			for job := range jobs {
				// workers send to results
				results <- handle(job)
			}
		}()

		// when jobs is closed, workers exit
	}

	// coordinator goroutine
	go func() {
		// WaitGroup waits for all workers
		wg.Wait()
		// one separate goroutine closes results
		close(results)
	}()

	return results
}
