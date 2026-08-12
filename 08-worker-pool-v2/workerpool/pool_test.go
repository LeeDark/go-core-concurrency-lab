package workerpool

import (
	"context"
	"errors"
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

func TestRunPassesCancellationToHandler(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan Job)
	handlerStarted := make(chan struct{})
	results := Run(ctx, 1, jobs, func(handlerCtx context.Context, job Job) Result {
		close(handlerStarted)

		<-handlerCtx.Done()
		return Result{JobID: job.ID, Err: handlerCtx.Err()}
	})

	go func() {
		defer close(jobs)
		jobs <- Job{ID: 1}
	}()

	select {
	case <-handlerStarted:
	case <-time.After(resultWaitTimeout):
		t.Fatal("handler did not start before timeout")
	}

	cancel()

	select {
	case result, ok := <-results:
		if !ok {
			t.Fatal("results closed before canceled handler returned")
		}
		if result.JobID != 1 {
			t.Fatalf("result job ID = %d, want 1", result.JobID)
		}
		if !errors.Is(result.Err, context.Canceled) {
			t.Fatalf("result error = %v, want %v", result.Err, context.Canceled)
		}
	case <-time.After(resultWaitTimeout):
		t.Fatal("canceled handler did not return a result before timeout")
	}

	select {
	case _, ok := <-results:
		if ok {
			t.Fatal("received more than one result")
		}
	case <-time.After(resultWaitTimeout):
		t.Fatal("results channel did not close after handler exit")
	}
}

func TestRunStopsWorkerBlockedOnResultSendAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan Job, 1)
	jobs <- Job{ID: 1}
	close(jobs)

	handlerFinished := make(chan struct{})
	results := Run(ctx, 1, jobs, func(_ context.Context, job Job) Result {
		close(handlerFinished)
		return Result{JobID: job.ID}
	})

	select {
	case <-handlerFinished:
	case <-time.After(resultWaitTimeout):
		t.Fatal("handler did not finish before timeout")
	}

	// No consumer receives the result, so the worker must be blocked on results.
	cancel()

	select {
	case _, ok := <-results:
		if ok {
			t.Fatal("received a result after cancellation from an abandoned consumer")
		}
	case <-time.After(resultWaitTimeout):
		t.Fatal("results channel did not close after canceling blocked result send")
	}
}
