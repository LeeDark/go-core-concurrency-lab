package workerpool

import (
	"errors"
	"fmt"
	"testing"
)

func TestPoolProcessesAllJobs(t *testing.T) {
	jobs := make(chan Job)
	results := Run(3, jobs, func(job Job) Result {
		return Result{
			JobID: job.ID,
			Value: job.Payload,
		}
	})

	go func() {
		defer close(jobs)

		for i := 1; i <= 10; i++ {
			jobs <- Job{
				ID:      i,
				Payload: fmt.Sprintf("job-%d", i),
			}
		}
	}()

	received := make(map[int]Result)
	for result := range results {
		received[result.JobID] = result
	}

	if len(received) != 10 {
		t.Fatalf("received %d results, want 10", len(received))
	}

	for i := 1; i <= 10; i++ {
		result, ok := received[i]
		if !ok {
			t.Fatalf("missing result for job %d", i)
		}

		want := fmt.Sprintf("job-%d", i)
		if result.Value != want {
			t.Fatalf("result value for job %d = %q, want %q", i, result.Value, want)
		}
	}
}

func TestPoolHandlesJobErrors(t *testing.T) {
	jobErr := errors.New("job failed")

	jobs := make(chan Job)
	results := Run(2, jobs, func(job Job) Result {
		result := Result{
			JobID: job.ID,
			Value: job.Payload,
		}

		if job.ID == 2 {
			result.Err = jobErr
		}

		return result
	})

	go func() {
		defer close(jobs)

		for i := 1; i <= 3; i++ {
			jobs <- Job{
				ID:      i,
				Payload: fmt.Sprintf("job-%d", i),
			}
		}
	}()

	var failed Result
	for result := range results {
		if result.JobID == 2 {
			failed = result
		}
	}

	if !errors.Is(failed.Err, jobErr) {
		t.Fatalf("result error = %v, want %v", failed.Err, jobErr)
	}
}

func TestPoolClosesResults(t *testing.T) {
	jobs := make(chan Job)
	results := Run(2, jobs, func(job Job) Result {
		return Result{JobID: job.ID}
	})

	go func() {
		defer close(jobs)

		jobs <- Job{ID: 1}
		jobs <- Job{ID: 2}
	}()

	count := 0
	for range results {
		count++
	}

	if count != 2 {
		t.Fatalf("received %d results, want 2", count)
	}

	_, ok := <-results
	if ok {
		t.Fatal("results channel is still open")
	}
}

func TestPoolKeepsResultsOpenUntilWorkerExits(t *testing.T) {
	jobs := make(chan Job)
	workerStarted := make(chan struct{})
	releaseWorker := make(chan struct{})

	results := Run(1, jobs, func(job Job) Result {
		close(workerStarted)
		<-releaseWorker

		return Result{JobID: job.ID}
	})

	go func() {
		defer close(jobs)
		jobs <- Job{ID: 1}
	}()

	<-workerStarted

	select {
	case _, ok := <-results:
		if !ok {
			t.Fatal("results closed before the worker exited")
		}
		t.Fatal("received a result before the worker was released")
	default:
		// The worker is still handling the job, so results must remain open and empty.
	}

	close(releaseWorker)

	result, ok := <-results
	if !ok {
		t.Fatal("results closed before receiving the worker result")
	}
	if result.JobID != 1 {
		t.Fatalf("result job ID = %d, want 1", result.JobID)
	}

	_, ok = <-results
	if ok {
		t.Fatal("results channel is still open after the worker exited")
	}
}

func TestPoolNormalizesNonPositiveWorkerCount(t *testing.T) {
	jobs := make(chan Job)
	results := Run(0, jobs, func(job Job) Result {
		return Result{JobID: job.ID}
	})

	go func() {
		defer close(jobs)
		jobs <- Job{ID: 1}
	}()

	count := 0
	for range results {
		count++
	}

	if count != 1 {
		t.Fatalf("received %d results, want 1", count)
	}
}
