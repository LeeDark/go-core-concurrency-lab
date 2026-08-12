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

func TestRunPreservesJobLevelErrors(t *testing.T) {
	jobErr := errors.New("job failed")
	jobs := make(chan Job, 3)
	for i := 1; i <= 3; i++ {
		jobs <- Job{ID: i}
	}
	close(jobs)

	results := Run(context.Background(), 2, jobs, func(_ context.Context, job Job) Result {
		result := Result{JobID: job.ID}
		if job.ID == 2 {
			result.Err = jobErr
		}
		return result
	})

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

func TestRunClosesResultsForClosedJobs(t *testing.T) {
	jobs := make(chan Job)
	close(jobs)

	results := Run(context.Background(), 2, jobs, func(_ context.Context, job Job) Result {
		return Result{JobID: job.ID}
	})

	select {
	case _, ok := <-results:
		if ok {
			t.Fatal("received a result for an empty jobs channel")
		}
	case <-time.After(resultWaitTimeout):
		t.Fatal("results channel did not close for a closed jobs channel")
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
	handlerCanceled := make(chan struct{})
	results := Run(ctx, 1, jobs, func(handlerCtx context.Context, job Job) Result {
		close(handlerStarted)

		<-handlerCtx.Done()
		close(handlerCanceled)
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
	case <-handlerCanceled:
	case <-time.After(resultWaitTimeout):
		t.Fatal("handler did not observe cancellation before timeout")
	}

	timer := time.NewTimer(resultWaitTimeout)
	defer timer.Stop()
	for {
		select {
		case result, ok := <-results:
			if !ok {
				return
			}
			if result.JobID != 1 {
				t.Fatalf("result job ID = %d, want 1", result.JobID)
			}
			if !errors.Is(result.Err, context.Canceled) {
				t.Fatalf("result error = %v, want %v", result.Err, context.Canceled)
			}
		case <-timer.C:
			t.Fatal("results did not close or receive a result before timeout")
		}
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

func TestRunClosesResultsOnlyAfterAllWorkersExit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const workerCount = 2
	jobs := make(chan Job, workerCount)
	for i := 1; i <= workerCount; i++ {
		jobs <- Job{ID: i}
	}
	close(jobs)

	handlerStarted := make(chan struct{}, workerCount)
	releaseHandlers := make(chan struct{})
	results := Run(ctx, workerCount, jobs, func(_ context.Context, job Job) Result {
		handlerStarted <- struct{}{}
		<-releaseHandlers
		return Result{JobID: job.ID}
	})

	for i := range workerCount {
		select {
		case <-handlerStarted:
		case <-time.After(resultWaitTimeout):
			t.Fatalf("handler %d did not start before timeout", i+1)
		}
	}

	cancel()

	select {
	case _, ok := <-results:
		if !ok {
			t.Fatal("results closed while handlers were still running")
		}
		t.Fatal("received a result before handlers were released")
	default:
	}

	close(releaseHandlers)

	timer := time.NewTimer(resultWaitTimeout)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-results:
			if !ok {
				return
			}
		case <-timer.C:
			t.Fatal("results did not close after all workers exited")
		}
	}
}

func TestRunSupportsWholeOperationTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	jobs := make(chan Job, 1)
	jobs <- Job{ID: 1}
	close(jobs)

	handlerErr := make(chan error, 1)
	results := Run(ctx, 1, jobs, func(handlerCtx context.Context, job Job) Result {
		<-handlerCtx.Done()
		handlerErr <- handlerCtx.Err()
		return Result{JobID: job.ID, Err: handlerCtx.Err()}
	})

	timer := time.NewTimer(resultWaitTimeout)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-results:
			if !ok {
				select {
				case err := <-handlerErr:
					if !errors.Is(err, context.DeadlineExceeded) {
						t.Fatalf("handler context error = %v, want %v", err, context.DeadlineExceeded)
					}
				case <-time.After(resultWaitTimeout):
					t.Fatal("handler did not report its context error")
				}
				return
			}
		case <-timer.C:
			t.Fatal("results did not close after whole-operation timeout")
		}
	}
}

func TestRunCannotStopHandlerThatIgnoresContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan Job, 1)
	jobs <- Job{ID: 1}
	close(jobs)

	handlerStarted := make(chan struct{})
	releaseHandler := make(chan struct{})
	results := Run(ctx, 1, jobs, func(_ context.Context, job Job) Result {
		close(handlerStarted)
		<-releaseHandler
		return Result{JobID: job.ID}
	})

	select {
	case <-handlerStarted:
	case <-time.After(resultWaitTimeout):
		t.Fatal("handler did not start before timeout")
	}

	cancel()

	select {
	case _, ok := <-results:
		if !ok {
			t.Fatal("results closed while handler ignored cancellation")
		}
		t.Fatal("received a result while handler was still blocked")
	default:
	}

	close(releaseHandler)

	timer := time.NewTimer(resultWaitTimeout)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-results:
			if !ok {
				return
			}
		case <-timer.C:
			t.Fatal("results did not close after handler was released")
		}
	}
}

func TestRunNormalizesNonPositiveWorkerCount(t *testing.T) {
	jobs := make(chan Job, 1)
	jobs <- Job{ID: 1}
	close(jobs)

	results := Run(context.Background(), 0, jobs, func(_ context.Context, job Job) Result {
		return Result{JobID: job.ID}
	})

	count := 0
	for range results {
		count++
	}

	if count != 1 {
		t.Fatalf("received %d results, want 1", count)
	}
}

func TestProducerStopsSubmittingAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan Job)
	results := Run(ctx, 1, jobs, func(_ context.Context, job Job) Result {
		return Result{JobID: job.ID}
	})

	producerDone := make(chan struct{})
	go func() {
		defer close(producerDone)
		defer close(jobs)

		for i := 1; i <= 2; i++ {
			select {
			case jobs <- Job{ID: i}:
			case <-ctx.Done():
				return
			}
		}
	}()

	cancel()

	select {
	case <-producerDone:
	case <-time.After(resultWaitTimeout):
		t.Fatal("producer did not stop after cancellation")
	}

	timer := time.NewTimer(resultWaitTimeout)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-results:
			if !ok {
				return
			}
		case <-timer.C:
			t.Fatal("results did not close after producer and pool stopped")
		}
	}
}

func TestRunLimitsConcurrentHandlers(t *testing.T) {
	const workerCount = 3

	jobs := make(chan Job, workerCount)
	for i := 1; i <= workerCount; i++ {
		jobs <- Job{ID: i}
	}
	close(jobs)

	started := make(chan struct{}, workerCount)
	release := make(chan struct{})
	var active atomic.Int32
	var maxActive atomic.Int32

	results := Run(context.Background(), workerCount, jobs, func(_ context.Context, job Job) Result {
		current := active.Add(1)
		for {
			previous := maxActive.Load()
			if current <= previous || maxActive.CompareAndSwap(previous, current) {
				break
			}
		}

		started <- struct{}{}
		<-release
		active.Add(-1)
		return Result{JobID: job.ID}
	})

	for i := range workerCount {
		select {
		case <-started:
		case <-time.After(resultWaitTimeout):
			t.Fatalf("only %d handlers started, want %d", i, workerCount)
		}
	}

	if got := maxActive.Load(); got != workerCount {
		t.Fatalf("maximum concurrent handlers = %d, want %d", got, workerCount)
	}

	close(release)

	count := 0
	for range results {
		count++
	}
	if count != workerCount {
		t.Fatalf("received %d results, want %d", count, workerCount)
	}
}
