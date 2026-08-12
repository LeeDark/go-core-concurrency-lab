// Package integration composes defer, errors, and context in one small
// operation lifecycle.
package integration

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	contextlab "github.com/LeeDark/go-core-concurrency-lab/03-defer-errors-context/context"
	errorslab "github.com/LeeDark/go-core-concurrency-lab/03-defer-errors-context/errors"
)

var (
	// ErrNilContext reports an invalid nil context argument.
	ErrNilContext = errors.New("nil context")
	// ErrNilResource reports an operation without an owned resource.
	ErrNilResource = errors.New("nil resource")
)

// Resource is an operation-owned resource that must be closed exactly once.
type Resource interface {
	Close() error
}

// RunOperation runs steps with a resource and composes the three Track A
// mechanisms:
//   - context controls whether steps should continue;
//   - errors preserve step and operation causes;
//   - defer guarantees resource cleanup on every return path.
//
// The caller owns ctx and resource creation. RunOperation takes ownership
// after it accepts a usable resource, before validating ctx, and closes that
// resource exactly once. A nil interface and an interface containing a nil
// value are both rejected. If both the operation and cleanup fail, the
// returned error contains both causes through errors.Join.
func RunOperation(
	ctx context.Context,
	operationID int,
	resource Resource,
	steps ...contextlab.Step,
) (err error) {
	if nilResource(resource) {
		return ErrNilResource
	}

	defer func() {
		closeErr := resource.Close()
		if closeErr == nil {
			return
		}

		closeErr = fmt.Errorf("operation %d: close resource: %w", operationID, closeErr)
		if err == nil {
			err = closeErr
			return
		}
		err = errors.Join(err, closeErr)
	}()

	if ctx == nil {
		return ErrNilContext
	}

	if ctxErr := ctx.Err(); ctxErr != nil {
		return errorslab.FailedOperation(operationID, "canceled before start", ctxErr)
	}

	if stepErr := contextlab.Run(ctx, steps...); stepErr != nil {
		return errorslab.FailedOperation(operationID, "run steps", stepErr)
	}

	return nil
}

func nilResource(resource Resource) bool {
	if resource == nil {
		return true
	}

	value := reflect.ValueOf(resource)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
