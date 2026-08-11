package contextlab

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestRunExecutesStepsInOrder(t *testing.T) {
	var order []int

	err := Run(context.Background(),
		func(context.Context) error {
			order = append(order, 1)
			return nil
		},
		func(context.Context) error {
			order = append(order, 2)
			return nil
		},
		func(context.Context) error {
			order = append(order, 3)
			return nil
		},
	)

	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if fmt.Sprint(order) != "[1 2 3]" {
		t.Fatalf("step order = %v, want [1 2 3]", order)
	}
}

func TestRunStopsBeforeFirstStepWhenContextIsCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	called := false
	err := Run(ctx, func(context.Context) error {
		called = true
		return nil
	})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Run error = %v, want context.Canceled", err)
	}
	if called {
		t.Fatal("step ran after context cancellation")
	}
}

func TestRunStopsBetweenSteps(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	secondCalled := false
	err := Run(ctx,
		func(context.Context) error {
			cancel()
			return nil
		},
		func(context.Context) error {
			secondCalled = true
			return nil
		},
	)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Run error = %v, want context.Canceled", err)
	}
	if secondCalled {
		t.Fatal("second step ran after cancellation")
	}
}

func TestStepCanObserveCancellationWhileRunning(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	started := make(chan struct{})
	done := make(chan error, 1)

	go func() {
		done <- Run(ctx, func(ctx context.Context) error {
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
			t.Fatalf("Run error = %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Run did not finish after cancellation")
	}
}

func TestRunWrapsStepErrorWithIndex(t *testing.T) {
	root := errors.New("database unavailable")

	err := Run(context.Background(),
		func(context.Context) error { return nil },
		func(context.Context) error { return root },
	)

	if !errors.Is(err, root) {
		t.Fatalf("Run error = %v, want wrapped root error", err)
	}
	stepErr, ok := errors.AsType[*StepError](err)
	if !ok {
		t.Fatal("Run error is not a StepError")
	}
	if stepErr.Index != 1 {
		t.Fatalf("step index = %d, want 1", stepErr.Index)
	}
}

func TestRunRejectsInvalidInputs(t *testing.T) {
	tests := map[string]struct {
		ctx  context.Context
		step Step
		want error
	}{
		"nil context": {ctx: nil, want: ErrNilContext},
		"nil step":    {ctx: context.Background(), step: nil, want: ErrNilStep},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := Run(test.ctx, test.step)
			if !errors.Is(err, test.want) {
				t.Fatalf("Run error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestWaitReturnsWhenSignalArrives(t *testing.T) {
	signal := make(chan struct{})
	close(signal)

	if err := Wait(context.Background(), signal); err != nil {
		t.Fatalf("Wait returned error: %v", err)
	}
}

func TestWaitReturnsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := Wait(ctx, make(chan struct{})); !errors.Is(err, context.Canceled) {
		t.Fatalf("Wait error = %v, want context.Canceled", err)
	}
}

func TestWaitReturnsDeadlineExceeded(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	if err := Wait(ctx, make(chan struct{})); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Wait error = %v, want context.DeadlineExceeded", err)
	}
}

func TestParentCancellationPropagatesToChild(t *testing.T) {
	parent, cancelParent := context.WithCancel(context.Background())
	child, cancelChild := context.WithCancel(parent)
	defer cancelChild()

	cancelParent()

	if !errors.Is(parent.Err(), context.Canceled) {
		t.Fatalf("parent error = %v, want context.Canceled", parent.Err())
	}
	if !errors.Is(child.Err(), context.Canceled) {
		t.Fatalf("child error = %v, want context.Canceled", child.Err())
	}
}
