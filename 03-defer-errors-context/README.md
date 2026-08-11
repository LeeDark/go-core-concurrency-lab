# Defer, Errors, And Context

## Status

Phase 4 Track A is in progress. The first three stages are implemented:

1. defer — core cheatsheet section and executable examples;
2. errors — sentinel errors, typed errors, wrapping, inspection, and tests;
3. context — cancellable workflow helpers, deadline-aware waiting, and tests.

The final integrated core lab is implemented in `integration/` and combines cleanup, error
propagation, and cancellation.

## Goal

Build a clear mental model for three Go mechanisms used together in concurrent code:

```text
defer   -> local cleanup
errors  -> explicit failure and cause inspection
context -> cancellation and deadline propagation
```

This is a learning lab, not a production framework. The examples keep ownership and control flow
visible instead of hiding them behind a large abstraction.

## Structure

```text
03-defer-errors-context/
  README.md
  defer/
    defer_examples_test.go
  errors/
    errors.go
    errors_test.go
  context/
    context.go
    context_test.go
  integration/
    operation.go
    operation_test.go
```

Theory is maintained in the three core cheatsheets:

- [`docs/cheatsheet-core.md`](../docs/cheatsheet-core.md);
- [`docs/cheatsheet-core.ru.md`](../docs/cheatsheet-core.ru.md);
- [`docs/cheatsheet-core.ua.md`](../docs/cheatsheet-core.ua.md).

## Stage 1: Defer

The stage covers execution on normal and early return, LIFO order, immediate argument evaluation,
closures, cleanup after resource acquisition, named return values, panic unwinding, and the resource
retention risk of defer inside a long loop.

Executable examples:

```text
defer/defer_examples_test.go
```

```bash
go test ./03-defer-errors-context/defer
```

## Stage 2: Errors

The `errors` package demonstrates how to preserve operation context and the underlying cause.

It covers:

- sentinel errors and `errors.Is`;
- `%w` wrapping;
- typed errors and `Unwrap`;
- `errors.As` and Go 1.26 `errors.AsType`;
- `errors.Join` and error trees;
- negative cases and table-driven tests.

The package is named `errorslab` to avoid colliding with the standard `errors` package.

```go
type FieldError struct {
	Field string
	Cause error
}

type JobError struct {
	JobID int
	Op    string
	Cause error
}

func NotFound(resource string) error
func InvalidField(field string, cause error) error
func FailedJob(jobID int, operation string, cause error) error
func FieldErrorFrom(err error) (*FieldError, bool)
func JobErrorFrom(err error) (*JobError, bool)
```

Both typed errors implement `Unwrap`, so callers can inspect structured fields and still use
`errors.Is` or `errors.AsType` for the underlying cause.

```bash
go test ./03-defer-errors-context/errors
```

## Stage 3: Context

The `context` package demonstrates cooperative cancellation and deadline propagation in a small
sequential workflow.

The core rules are:

- the caller normally creates and cancels a context;
- a child inherits cancellation from its parent;
- canceling a child does not cancel its parent;
- cancellation closes `ctx.Done()`;
- `ctx.Err()` reports `context.Canceled` or `context.DeadlineExceeded`;
- a step must observe context to stop promptly;
- context cannot forcibly interrupt arbitrary code;
- context values are for request-scoped metadata, not mandatory business arguments.

The package is named `contextlab` to avoid colliding with the standard `context` package.

```go
type Step func(context.Context) error

type StepError struct {
	Index int
	Cause error
}

func Run(ctx context.Context, steps ...Step) error
func Wait(ctx context.Context, signal <-chan struct{}) error
```

`Run` checks context before every step, passes it to the step, and wraps step failures with the step
index. `Wait` returns when either the caller-owned signal channel is ready or context is canceled.

```bash
go test ./03-defer-errors-context/context
```

## Stage 4: Integrated core lab

The `integration` package composes the three mechanisms in one operation lifecycle:

```text
caller
  -> context.WithCancel / context.WithTimeout
  -> RunOperation
       -> owned resource
       -> defer cleanup
       -> contextlab.Run steps
       -> wrapped operation error
       -> resource.Close
```

### API

```go
type Resource interface {
	Close() error
}

func RunOperation(
	ctx context.Context,
	operationID int,
	resource Resource,
	steps ...contextlab.Step,
) (err error)
```

The caller creates the context and resource. After accepting a non-nil resource, `RunOperation`
closes it exactly once on every return path.

The operation preserves structured causes:

- step failures contain `contextlab.StepError` and `errorslab.OperationError`;
- cancellation remains discoverable with `errors.Is`;
- close failures are wrapped;
- a step failure and a close failure are combined with `errors.Join`.

The integration tests use a controlled fake resource rather than files, sockets, or databases. They
cover success, step failure, cancellation before and during a step, deadline expiration, cleanup
errors, joined errors, and invalid input.

Run the focused tests with:

```bash
go test ./03-defer-errors-context/integration
```

## Combined lifecycle model

The three mechanisms answer different questions:

```text
defer   — what must be cleaned up when this function exits?
errors  — what failed, and what was the original cause?
context — should this operation continue, or must it stop?
```

A typical operation uses them together:

```go
func RunOperation(parent context.Context) error {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	if err := contextlab.Run(ctx, stepOne, stepTwo); err != nil {
		return fmt.Errorf("run operation: %w", err)
	}
	return nil
}
```

The caller supplies the parent context, the function releases its child context with `defer`, and
the returned error retains its cause through wrapping.

## Targeted checks

```bash
go test ./03-defer-errors-context/defer
go test ./03-defer-errors-context/errors
go test ./03-defer-errors-context/context
go test ./03-defer-errors-context/integration
```

Combined targeted check:

```bash
go test ./03-defer-errors-context/defer \
	./03-defer-errors-context/errors \
	./03-defer-errors-context/context \
	./03-defer-errors-context/integration
```

Use the race detector after lifecycle examples are stable:

```bash
go test -race ./03-defer-errors-context/context
go test -race ./03-defer-errors-context/integration
```

Do not use `go test ./...` as the default check for this lab.

## Stop line

This lab does not include worker pools, pipelines, retry policy, `errgroup`, metrics, tracing, HTTP
shutdown, OS signals, persistent queues, or a general service framework. Those topics build on the
mechanisms here and belong to later labs or phases.
