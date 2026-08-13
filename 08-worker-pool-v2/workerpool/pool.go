package workerpool

import (
	"context"
	"sync"
	"time"
)

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

// WithJobTimeout wraps a handler with a timeout for each individual job. The
// parent pool context is still honored, and the wrapper does not change channel
// ownership or the pool's whole-operation cancellation policy.
func WithJobTimeout(
	timeout time.Duration,
	handle func(context.Context, Job) Result,
) func(context.Context, Job) Result {
	return func(ctx context.Context, job Job) Result {
		jobCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		result := handle(jobCtx, job)
		if result.Err == nil && jobCtx.Err() != nil {
			result.Err = jobCtx.Err()
		}
		return result
	}
}

// Run starts a fixed number of workers and returns a receive-only results channel.
// The caller owns sending to and closing jobs. The pool owns sending to and closing
// the returned results channel. A non-positive workerCount is normalized to one.
// Workers stop waiting for jobs when ctx is canceled. A handler must observe ctx
// itself if it needs to stop while processing a job.
func Run(
	ctx context.Context,
	workerCount int,
	jobs <-chan Job,
	handle func(context.Context, Job) Result,
) <-chan Result {
	if workerCount <= 0 {
		workerCount = 1
	}

	results := make(chan Result)

	var wg sync.WaitGroup
	wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()

			for {
				select {
				case job, ok := <-jobs:
					if !ok {
						return
					}

					result := handle(ctx, job)
					if ctx.Err() != nil {
						return
					}
					select {
					case results <- result:
					case <-ctx.Done():
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
