# Worker Pool v1

## Goal

Build the first minimal worker pool:

```text
producer -> jobs channel -> N workers -> results channel -> consumer
```

This version is only about the basic mechanics:

- goroutines;
- channel ownership;
- fixed worker count;
- `sync.WaitGroup`;
- closing `results` correctly.

Keep v1 small enough to reason about line by line.

## Scope

Include:

- `Job` struct;
- `Result` struct;
- `Run` function;
- fixed number of workers;
- input `jobs` channel;
- output `results` channel;
- handler function;
- `sync.WaitGroup`;
- one goroutine that closes `results` after all workers finish.

Do not include:

- `context.Context`;
- cancellation;
- timeouts;
- graceful shutdown;
- goroutine leak checks;
- race-detector work;
- advanced error handling;
- tests unless explicitly requested.

Those topics belong to Worker Pool v2.

## Structure

Use this structure:

```text
06-worker-pool-v1/
  README.md
  primitives/
    main.go
  workerpool/
    pool.go
```

`primitives/main.go` stays as the scratchpad for goroutine and channel examples.

`workerpool/pool.go` is the clean implementation.

## First API

Start with this API:

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

func Run(workerCount int, jobs <-chan Job, handle func(Job) Result) <-chan Result
```

Why this shape:

- caller owns sending jobs;
- caller closes `jobs`;
- pool owns sending results;
- pool closes `results`;
- `jobs <-chan Job` means the pool can only receive jobs;
- returned `<-chan Result` means the caller can only receive results.

## Implementation Steps

### Step 1: Define Types

Create:

- `Job`;
- `Result`.

Keep fields simple. This lab is about concurrency mechanics, not business logic.

### Step 2: Create Results Channel

Inside `Run`:

```text
results := make(chan Result)
```

Return this channel immediately after starting workers and the closer goroutine.

### Step 3: Start Workers

Start `workerCount` goroutines.

Each worker should:

```text
range over jobs
handle each job
send result to results
return when jobs is closed
```

### Step 4: Wait For Workers

Use `sync.WaitGroup`:

```text
Add before starting each worker
Done when each worker exits
Wait in a separate closer goroutine
```

### Step 5: Close Results

Close `results` only after all workers finish:

```text
wg.Wait()
close(results)
```

Do not close `results` from individual workers.

### Step 6: Add Simple Tests

The current test file is:

```text
workerpool/pool_test.go
```

It covers:

- processing all submitted jobs;
- preserving job-level errors in `Result.Err`;
- closing `results` after workers finish;
- normalizing non-positive `workerCount` to one worker.

Use only targeted commands for this lab:

```bash
go test ./06-worker-pool-v1/workerpool
go test -race ./06-worker-pool-v1/workerpool
```

## Channel Ownership

`jobs`:

- caller sends jobs;
- caller closes `jobs`;
- workers receive jobs;
- workers never close `jobs`.

`results`:

- workers send results;
- caller receives results;
- pool closes `results`;
- caller never closes `results`;
- individual workers never close `results`.

## Small Demo Flow

The first demo should be this simple:

```text
create jobs channel
start pool
send jobs
close jobs
range over results
print results
```

Do not add cancellation or timeout to this demo.

## Learning Checkpoints

Before moving to Worker Pool v2, answer:

1. Who closes `jobs`?

   The caller closes `jobs`, because the caller owns sending jobs.

2. Who closes `results`?

   The pool closes `results`, because the pool owns sending results.

3. Why should workers not close `results`?

   Multiple workers send to `results`, so no single worker knows that all sends are finished.

4. Why is `sync.WaitGroup` needed?

   It lets the closer goroutine wait until all workers exit before closing `results`.

5. What happens when a worker ranges over a closed `jobs` channel?

   The worker receives any buffered jobs first, then the `range` loop ends.

6. What happens if nobody reads from `results`?

   Workers block while sending results, so the pool cannot finish.

7. What happens when sending to a closed channel?

   The program panics.

8. What happens when reading from a closed channel?

   The receive succeeds immediately with the zero value and `ok == false`.

9. Why is this bounded concurrency?

   The number of concurrent workers is fixed by `workerCount`.

10. Why does v1 avoid `context.Context`?

   V1 focuses on channel ownership, worker lifecycle, and `WaitGroup` before adding cancellation.

## Interview Angle

Useful wording:

> I usually avoid unbounded goroutine creation for message or request processing. I prefer bounded worker pools or semaphores when throughput and external dependencies need to be controlled.

Short v1 explanation:

```text
I start a fixed number of workers instead of one goroutine per job.
The jobs channel distributes work across workers.
The results channel collects processed output.
A WaitGroup lets the pool close results after every worker exits.
```

## Stop Line

Stop v1 when the basic worker pool is clear.

Move to `07-worker-pool-v2` for:

- context cancellation;
- timeout;
- cancellation-aware result sending;
- leak reasoning;
- race detector;
- graceful stop semantics.
