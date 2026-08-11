package workerpool

import "sync"

type Job struct {
	ID      int
	Payload string
}

type Result struct {
	JobID int
	Value string
	Err   error
}

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
