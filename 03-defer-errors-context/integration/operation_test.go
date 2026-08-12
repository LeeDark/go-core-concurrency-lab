package integration

import (
	"context"
	"errors"
	"testing"
	"time"

	contextlab "github.com/LeeDark/go-core-concurrency-lab/03-defer-errors-context/context"
	errorslab "github.com/LeeDark/go-core-concurrency-lab/03-defer-errors-context/errors"
)

type fakeResource struct {
	closeCalls int
	closeErr   error
}

func (r *fakeResource) Close() error {
	r.closeCalls++
	return r.closeErr
}

func TestRunOperationSuccessClosesResource(t *testing.T) {
	resource := &fakeResource{}
	var order []int

	err := RunOperation(context.Background(), 10, resource,
		func(context.Context) error {
			order = append(order, 1)
			return nil
		},
		func(context.Context) error {
			order = append(order, 2)
			return nil
		},
	)

	if err != nil {
		t.Fatalf("RunOperation returned error: %v", err)
	}
	if resource.closeCalls != 1 {
		t.Fatalf("Close calls = %d, want 1", resource.closeCalls)
	}
	if len(order) != 2 || order[0] != 1 || order[1] != 2 {
		t.Fatalf("step order = %v, want [1 2]", order)
	}
}

func TestRunOperationWrapsStepAndOperationErrors(t *testing.T) {
	root := errors.New("database unavailable")
	resource := &fakeResource{}
	thirdCalled := false

	err := RunOperation(context.Background(), 42, resource,
		func(context.Context) error { return nil },
		func(context.Context) error { return root },
		func(context.Context) error {
			thirdCalled = true
			return nil
		},
	)

	if !errors.Is(err, root) {
		t.Fatalf("error = %v, want root cause", err)
	}
	operationErr, ok := errors.AsType[*errorslab.OperationError](err)
	if !ok || operationErr.OperationID != 42 {
		t.Fatalf("operation error = %#v, want operation ID 42", operationErr)
	}
	stepErr, ok := errors.AsType[*contextlab.StepError](err)
	if !ok || stepErr.Index != 1 {
		t.Fatalf("step error = %#v, want step index 1", stepErr)
	}
	if thirdCalled {
		t.Fatal("step after failure was called")
	}
	if resource.closeCalls != 1 {
		t.Fatalf("Close calls = %d, want 1", resource.closeCalls)
	}
}

func TestRunOperationStopsBeforeFirstStepWhenCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resource := &fakeResource{}
	called := false

	err := RunOperation(ctx, 1, resource, func(context.Context) error {
		called = true
		return nil
	})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
	if called {
		t.Fatal("step ran after cancellation")
	}
	if resource.closeCalls != 1 {
		t.Fatalf("Close calls = %d, want 1", resource.closeCalls)
	}
}

func TestRunOperationStopsInsideStepWhenCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resource := &fakeResource{}
	started := make(chan struct{})
	done := make(chan error, 1)

	go func() {
		done <- RunOperation(ctx, 2, resource, func(ctx context.Context) error {
			close(started)
			<-ctx.Done()
			return ctx.Err()
		})
	}()

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("step did not start")
	}
	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("error = %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("operation did not finish after cancellation")
	}
	if resource.closeCalls != 1 {
		t.Fatalf("Close calls = %d, want 1", resource.closeCalls)
	}
}

func TestRunOperationReturnsDeadlineExceeded(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	resource := &fakeResource{}

	err := RunOperation(ctx, 3, resource, func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %v, want context.DeadlineExceeded", err)
	}
	if resource.closeCalls != 1 {
		t.Fatalf("Close calls = %d, want 1", resource.closeCalls)
	}
}

func TestRunOperationReturnsCloseErrorOnSuccess(t *testing.T) {
	closeErr := errors.New("flush failed")
	resource := &fakeResource{closeErr: closeErr}

	err := RunOperation(context.Background(), 4, resource)

	if !errors.Is(err, closeErr) {
		t.Fatalf("error = %v, want close error", err)
	}
	if resource.closeCalls != 1 {
		t.Fatalf("Close calls = %d, want 1", resource.closeCalls)
	}
}

func TestRunOperationJoinsStepAndCloseErrors(t *testing.T) {
	stepErr := errors.New("step failed")
	closeErr := errors.New("close failed")
	resource := &fakeResource{closeErr: closeErr}

	err := RunOperation(context.Background(), 5, resource, func(context.Context) error {
		return stepErr
	})

	if !errors.Is(err, stepErr) {
		t.Fatalf("joined error = %v, missing step error", err)
	}
	if !errors.Is(err, closeErr) {
		t.Fatalf("joined error = %v, missing close error", err)
	}
	operationErr, ok := errors.AsType[*errorslab.OperationError](err)
	if !ok || operationErr.OperationID != 5 {
		t.Fatalf("operation error = %#v, want operation ID 5", operationErr)
	}
	stepFailure, ok := errors.AsType[*contextlab.StepError](err)
	if !ok || stepFailure.Index != 0 {
		t.Fatalf("step error = %#v, want step index 0", stepFailure)
	}
	if resource.closeCalls != 1 {
		t.Fatalf("Close calls = %d, want 1", resource.closeCalls)
	}
}

func TestRunOperationClosesAcceptedResourceWhenContextIsNil(t *testing.T) {
	resource := &fakeResource{}

	if err := RunOperation(nil, 6, resource); !errors.Is(err, ErrNilContext) {
		t.Fatalf("nil context error = %v, want ErrNilContext", err)
	}
	if resource.closeCalls != 1 {
		t.Fatalf("Close calls after nil context = %d, want 1", resource.closeCalls)
	}
}

func TestRunOperationRejectsNilResources(t *testing.T) {
	if err := RunOperation(context.Background(), 7, nil); !errors.Is(err, ErrNilResource) {
		t.Fatalf("nil resource error = %v, want ErrNilResource", err)
	}

	var typedNil *fakeResource
	if err := RunOperation(context.Background(), 8, typedNil); !errors.Is(err, ErrNilResource) {
		t.Fatalf("typed nil resource error = %v, want ErrNilResource", err)
	}
}
