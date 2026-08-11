// Package errorslab demonstrates error identity, wrapping, typed errors, and
// inspection of error chains and trees.
package errorslab

import (
	"errors"
	"fmt"
)

// Sentinel errors represent stable causes that callers may inspect with errors.Is.
var (
	ErrNotFound      = errors.New("resource not found")
	ErrPermission    = errors.New("permission denied")
	ErrInvalidInput  = errors.New("invalid input")
	ErrAlreadyExists = errors.New("resource already exists")
)

// FieldError identifies the input field that caused a validation or decoding
// failure. Cause remains available to errors.Is and errors.As through Unwrap.
type FieldError struct {
	Field string
	Cause error
}

func (e *FieldError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Cause == nil {
		return fmt.Sprintf("field %q is invalid", e.Field)
	}
	return fmt.Sprintf("field %q: %v", e.Field, e.Cause)
}

func (e *FieldError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// JobError adds operation and job identity to a failure while preserving its
// underlying cause.
type JobError struct {
	JobID int
	Op    string
	Cause error
}

func (e *JobError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Cause == nil {
		return fmt.Sprintf("job %d: %s", e.JobID, e.Op)
	}
	return fmt.Sprintf("job %d %s: %v", e.JobID, e.Op, e.Cause)
}

func (e *JobError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// NotFound returns an error with resource-specific context while preserving
// ErrNotFound for errors.Is.
func NotFound(resource string) error {
	return fmt.Errorf("load %s: %w", resource, ErrNotFound)
}

// InvalidField returns a typed validation error. A nil cause is allowed when
// the field itself is the complete explanation.
func InvalidField(field string, cause error) error {
	return &FieldError{Field: field, Cause: cause}
}

// FailedJob returns a typed operation error that preserves cause inspection.
func FailedJob(jobID int, operation string, cause error) error {
	return &JobError{JobID: jobID, Op: operation, Cause: cause}
}

// IsNotFound reports whether err's chain or tree contains ErrNotFound.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// FieldErrorFrom returns the first FieldError found in err's chain or tree.
func FieldErrorFrom(err error) (*FieldError, bool) {
	return errors.AsType[*FieldError](err)
}

// JobErrorFrom returns the first JobError found in err's chain or tree.
func JobErrorFrom(err error) (*JobError, bool) {
	return errors.AsType[*JobError](err)
}
