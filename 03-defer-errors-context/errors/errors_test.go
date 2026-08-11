package errorslab

import (
	"errors"
	"fmt"
	"testing"
)

func TestNotFoundPreservesSentinelIdentity(t *testing.T) {
	err := NotFound("settings")

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("errors.Is(%v, ErrNotFound) = false, want true", err)
	}
	if err.Error() != "load settings: resource not found" {
		t.Fatalf("error = %q, want contextual message", err)
	}
}

func TestFormattingVerbControlsInspection(t *testing.T) {
	wrapped := fmt.Errorf("load settings: %w", ErrNotFound)
	formatted := fmt.Errorf("load settings: %v", ErrNotFound)

	if !IsNotFound(wrapped) {
		t.Fatal("wrapped error is not recognized as ErrNotFound")
	}
	if IsNotFound(formatted) {
		t.Fatal("formatted error is recognized as ErrNotFound; formatted errors must not wrap")
	}
}

func TestSentinelComparisonUsesIdentity(t *testing.T) {
	fresh := errors.New("resource not found")
	if errors.Is(fresh, ErrNotFound) {
		t.Fatal("fresh error unexpectedly matches sentinel")
	}
	if !errors.Is(NotFound("settings"), ErrNotFound) {
		t.Fatal("wrapped sentinel was not found through errors.Is")
	}
}

func TestFieldErrorPreservesCauseAndContext(t *testing.T) {
	err := InvalidField("age", ErrInvalidInput)

	fieldErr, ok := FieldErrorFrom(err)
	if !ok {
		t.Fatal("FieldErrorFrom returned false")
	}
	if fieldErr.Field != "age" {
		t.Fatalf("field = %q, want %q", fieldErr.Field, "age")
	}
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatal("FieldError does not preserve its cause")
	}
}

func TestFieldErrorCanBeInspectedWithAs(t *testing.T) {
	err := fmt.Errorf("decode request: %w", InvalidField("age", ErrInvalidInput))

	var fieldErr *FieldError
	if !errors.As(err, &fieldErr) {
		t.Fatal("errors.As returned false")
	}
	if fieldErr == nil || fieldErr.Field != "age" {
		t.Fatalf("field error = %#v, want field age", fieldErr)
	}
}

func TestAsTypeReturnsTypedErrorWithoutTargetPointer(t *testing.T) {
	err := fmt.Errorf("decode request: %w", InvalidField("age", ErrInvalidInput))

	fieldErr, ok := errors.AsType[*FieldError](err)
	if !ok {
		t.Fatal("errors.AsType returned false")
	}
	if fieldErr == nil || fieldErr.Field != "age" {
		t.Fatalf("field error = %#v, want field age", fieldErr)
	}

	fromHelper, ok := FieldErrorFrom(err)
	if !ok || fromHelper.Field != fieldErr.Field {
		t.Fatal("FieldErrorFrom did not return the expected typed error")
	}
	if !errors.Is(fromHelper, ErrInvalidInput) {
		t.Fatal("FieldErrorFrom did not preserve the expected cause")
	}
}

func TestAsTypeReturnsZeroValueWhenTypeIsAbsent(t *testing.T) {
	jobErr, ok := errors.AsType[*JobError](InvalidField("age", ErrInvalidInput))
	if ok {
		t.Fatal("errors.AsType found an absent JobError")
	}
	if jobErr != nil {
		t.Fatalf("job error = %#v, want nil", jobErr)
	}

	if fieldErr, ok := FieldErrorFrom(nil); ok || fieldErr != nil {
		t.Fatalf("FieldErrorFrom(nil) = (%#v, %t), want (nil, false)", fieldErr, ok)
	}
}

func TestJobErrorPreservesOperationAndCause(t *testing.T) {
	err := FailedJob(42, "decode", ErrInvalidInput)

	jobErr, ok := JobErrorFrom(err)
	if !ok {
		t.Fatal("JobErrorFrom returned false")
	}
	if jobErr.JobID != 42 || jobErr.Op != "decode" {
		t.Fatalf("job error = %#v, want job ID 42 and operation decode", jobErr)
	}
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatal("JobError does not preserve its cause")
	}
}

func TestAsTypeUsesFirstDepthFirstMatchInErrorTree(t *testing.T) {
	first := InvalidField("name", ErrInvalidInput)
	second := InvalidField("age", ErrInvalidInput)
	err := errors.Join(
		fmt.Errorf("first: %w", first),
		fmt.Errorf("second: %w", second),
	)

	fieldErr, ok := errors.AsType[*FieldError](err)
	if !ok {
		t.Fatal("errors.AsType returned false for joined errors")
	}
	if fieldErr.Field != "name" {
		t.Fatalf("field = %q, want first field %q", fieldErr.Field, "name")
	}
}

func TestErrorsIsSearchesJoinedErrorTree(t *testing.T) {
	err := errors.Join(
		fmt.Errorf("permission: %w", ErrPermission),
		fmt.Errorf("missing: %w", ErrNotFound),
	)

	tests := map[string]struct {
		target error
		want   bool
	}{
		"permission": {target: ErrPermission, want: true},
		"not found":  {target: ErrNotFound, want: true},
		"invalid":    {target: ErrInvalidInput, want: false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got := errors.Is(err, test.target); got != test.want {
				t.Fatalf("errors.Is(%v, %v) = %t, want %t", err, test.target, got, test.want)
			}
		})
	}
}

func TestErrorMessagesHandleNilTypedReceivers(t *testing.T) {
	var fieldErr *FieldError
	var jobErr *JobError

	if got := fieldErr.Error(); got != "<nil>" {
		t.Fatalf("nil FieldError.Error() = %q, want %q", got, "<nil>")
	}
	if got := jobErr.Error(); got != "<nil>" {
		t.Fatalf("nil JobError.Error() = %q, want %q", got, "<nil>")
	}
}
