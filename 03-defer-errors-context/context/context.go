// Package contextlab demonstrates cancellation and deadline propagation in
// small, testable workflow helpers.
package contextlab

import (
	"context"
	"errors"
	"fmt"
)

var (
	// ErrNilContext reports an invalid nil context argument.
	ErrNilContext = errors.New("nil context")
	// ErrNilStep reports a nil workflow step.
	ErrNilStep = errors.New("nil workflow step")
)

// Step is one cancellable unit of a workflow.
type Step func(context.Context) error

// StepError adds the index of a failed step while preserving its cause.
type StepError struct {
	Index int
	Cause error
}

func (e *StepError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("workflow step %d: %v", e.Index, e.Cause)
}

func (e *StepError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// Run executes steps in order.
//
// It checks ctx before every step, passes ctx to each step, and returns the
// context error unchanged when cancellation happens between steps. Errors
// returned by a step are wrapped in StepError so callers can inspect both the
// step index and the original cause with errors.Is or errors.AsType.
func Run(ctx context.Context, steps ...Step) error {
	if ctx == nil {
		return ErrNilContext
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	for index, step := range steps {
		if err := ctx.Err(); err != nil {
			return err
		}
		if step == nil {
			return &StepError{Index: index, Cause: ErrNilStep}
		}

		if err := step(ctx); err != nil {
			return &StepError{Index: index, Cause: err}
		}
	}

	return nil
}

// Wait waits for signal or for ctx cancellation. The signal channel is
// receive-only because the caller owns its lifecycle and decides when the
// event is published or the channel is closed.
func Wait(ctx context.Context, signal <-chan struct{}) error {
	if ctx == nil {
		return ErrNilContext
	}

	select {
	case <-signal:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
