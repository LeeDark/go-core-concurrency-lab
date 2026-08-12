package workerpool

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

const resultWaitTimeout = time.Second

func TestRunProcessesJobsAndClosesResults(t *testing.T) {
	jobs := make(chan Job)
	results := Run(context.Background(), 2, jobs, func(_ context.Context, job Job) Result {
		return Result{JobID: job.ID, Value: job.Payload}
	})

	go func() {
		defer close(jobs)
		jobs <- Job{ID: 1, Payload: "first"}
		jobs <- Job{ID: 2, Payload: "second"}
	}()

	received := make(map[int]string)
	for result := range results {
		received[result.JobID] = result.Value
	}

	if len(received) != 2 {
		t.Fatalf("received %d results, want 2", len(received))
	}
	if received[1] != "first" || received[2] != "second" {
		t.Fatalf("received results = %#v, want both submitted jobs", received)
	}
}

func TestRunStopsWorkersWaitingForJobsAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	jobs := make(chan Job)
	var handlerRan atomic.Bool
	results := Run(ctx, 3, jobs, func(context.Context, Job) Result {
		handlerRan.Store(true)
		return Result{}
	})

	cancel()

	select {
	case _, ok := <-results:
		if ok {
			t.Fatal("received a result without a job")
		}
	case <-time.After(resultWaitTimeout):
		t.Fatal("results channel did not close after cancellation")
	}

	if handlerRan.Load() {
		t.Fatal("handler ran without a job")
	}
}
