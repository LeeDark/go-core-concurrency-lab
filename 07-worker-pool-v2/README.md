# Worker Pool v2

## Goal

Extend Worker Pool v1 with production-style lifecycle control:

```text
producer -> jobs channel -> N workers -> results channel -> consumer
                |              |
                +-- context ---+
```

Worker Pool v2 is about cancellation, timeouts, error policy, and avoiding goroutine leaks.

Start v2 only after Worker Pool v1 is clear.

## Starting Point

Worker Pool v1 should already cover:

- fixed number of workers;
- `jobs` channel;
- `results` channel;
- `Job` and `Result`;
- `sync.WaitGroup`;
- correct `results` closing;
- clear channel ownership.

Do not solve v2 by rewriting the whole lab from scratch. Add lifecycle control on top of the v1 mental model.

## V2 Scope

Worker Pool v2 includes:

- `context.Context`;
- cancellation while workers are waiting for jobs;
- cancellation while workers are sending results;
- timeout for the whole operation;
- optional timeout per job;
- explicit error handling policy;
- goroutine leak reasoning;
- race detector validation when explicitly requested;
- graceful stop semantics.

Worker Pool v2 does not need to include:

- HTTP server shutdown;
- OS signal handling;
- persistent queues;
- retry policy;
- rate limiting;
- metrics;
- tracing.

Those belong to later labs.

## Suggested Structure

```text
07-worker-pool-v2/
  README.md
  workerpool/
    pool.go
```

If v1 already has a clean `workerpool` package, copy the idea and evolve the API here.

## Possible API

Start with this shape:

```go
type Job struct {
	ID      int
	Payload string
}

type Result struct {
	JobID int
	Value string
	Err   error
}

func Run(
	ctx context.Context,
	workerCount int,
	jobs <-chan Job,
	handle func(context.Context, Job) Result,
) <-chan Result
```

Why this API:

- caller controls cancellation;
- workers can stop when `ctx.Done()` is closed;
- handler receives the same context;
- pool still owns `results`;
- caller still owns `jobs`.

## Cancellation Rules

Workers should stop when:

- `jobs` is closed;
- `ctx.Done()` is closed.

Workers must not get stuck forever trying to send a result after the consumer has stopped reading.

That means sends to `results` should usually happen inside `select`:

```text
select:
  send result to results
  or stop when ctx.Done() is closed
```

## Implementation Plan

### Step 1: Add Context To Run

Add `ctx context.Context` to `Run`.

Workers should select between:

- receiving from `jobs`;
- cancellation from `ctx.Done()`.

### Step 2: Add Context To Handler

Change handler shape:

```go
handle func(context.Context, Job) Result
```

This lets slow or blocking handlers observe cancellation.

### Step 3: Protect Result Sends

When a worker finishes a job, it should not blindly block forever on:

```go
results <- result
```

Use cancellation-aware sending.

### Step 4: Close Results Correctly

Keep the v1 rule:

```text
wait for all workers
close(results)
```

Cancellation does not change ownership.

The pool still closes `results` only after all worker goroutines exit.

### Step 5: Define Error Policy

Pick one simple policy for v2:

- every job produces a `Result`;
- `Result.Err` contains job-level failure;
- context cancellation may stop remaining jobs;
- the results channel closes when workers exit.

Do not introduce separate `errors` channel yet unless there is a strong reason.

## Timeout Plan

Whole operation timeout:

```go
ctx, cancel := context.WithTimeout(parent, timeout)
defer cancel()
```

Per-job timeout can be added later inside the handler:

```go
jobCtx, cancel := context.WithTimeout(ctx, jobTimeout)
defer cancel()
```

For the first v2 pass, prefer whole-operation timeout only.

## Leak Risks To Understand

Common leak cases:

- worker waits forever on `jobs`;
- worker waits forever sending to `results`;
- producer keeps sending after cancellation;
- consumer stops reading before workers finish;
- handler blocks and ignores context.

V2 should explain which of these are handled by the design and which remain caller responsibility.

## Channel Ownership Rules

`jobs`:

- caller sends jobs;
- caller closes `jobs`;
- pool receives jobs;
- pool does not close `jobs`.

`results`:

- pool sends results;
- caller receives results;
- pool closes `results`;
- caller does not close `results`.

`ctx.Done()`:

- caller creates context;
- caller cancels context;
- workers observe cancellation;
- workers do not close `ctx.Done()`.

## Learning Checkpoints

Before moving past v2, be able to answer:

1. How does context cancellation stop workers waiting for jobs?
2. How does context cancellation stop workers blocked on result sends?
3. Why does cancellation not mean the caller should close `results`?
4. What is the difference between closing `jobs` and canceling context?
5. What happens to queued jobs after cancellation?
6. Should the pool drain jobs after cancellation?
7. Should every job always produce a result?
8. What can still leak if `handle` ignores context?
9. What does `WaitGroup` still do in v2?
10. Where should timeout be created: caller, pool, or handler?

## Interview Angle

Useful wording:

> In a worker pool, cancellation has to cover both receiving work and publishing results. Otherwise a worker can still leak even if it listens to context while reading jobs.

For v2, focus the explanation on lifecycle:

```text
The caller owns cancellation through context.
Workers select on jobs and ctx.Done().
When publishing results, workers also select on ctx.Done() so they do not block forever.
A WaitGroup tracks worker exit, and the pool closes results after all workers are done.
```

## Relation To Later Labs

Worker Pool v2 prepares for:

- pipeline cancellation;
- fan-out/fan-in;
- rate-limited processing;
- graceful shutdown;
- shared-state alternatives.

Keep v2 focused on lifecycle control, not service architecture.
