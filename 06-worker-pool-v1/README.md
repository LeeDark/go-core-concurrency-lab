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
- keeping `results` open while a worker is still handling a job;
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

## Lab 2: Close Rules

This lab is about what closing a channel means and who is allowed to do it. It does not add cancellation, timeouts, or a new pool API.

### Rules

1. `close(ch)` means that no more values will be sent on `ch`. It is a completion signal, not resource cleanup.
2. The sending side usually closes a channel, because it knows when production is finished.
3. A receiver should not close a channel it does not own.
4. A receive from a closed channel first drains buffered values. Later receives return the zero value with `ok == false`.
5. Sending to a closed channel panics. Closing an already closed channel also panics.
6. Sending to or receiving from a `nil` channel blocks forever.

Closing is not needed for garbage collection. A channel that is no longer reachable can be collected whether it is open or closed.

### Small examples

The producer closes the channel it owns after sending all values:

```go
jobs := make(chan Job)

go func() {
	defer close(jobs)

	for _, job := range submittedJobs {
		jobs <- job
	}
}()
```

Use the comma-ok form when a zero value must be distinguished from a closed channel:

```go
value, ok := <-ch
if !ok {
	// ch is closed and drained
}
```

`range` uses the same rule: it receives buffered values first and then stops when the channel is closed and drained.

The following example intentionally panics. `recover` is used only to demonstrate the rule; it is not a normal way to handle channel ownership mistakes:

```go
ch := make(chan int)
close(ch)

func() {
	defer func() {
		fmt.Println(recover()) // send on closed channel
	}()

	ch <- 1
}()
```

The following examples intentionally block forever and must not be used as normal program flow:

```go
ch := make(chan int)
ch <- 1 // no receiver: deadlock

var nilCh chan int
<-nilCh // nil channel: blocks forever
```

### Applying the rules to this pool

```text
caller      -> sends jobs and closes jobs
workers     -> receive jobs; never close jobs
workers     -> send results
coordinator -> waits for every worker, then closes results
caller      -> receives results; never closes results
```

When the caller closes `jobs`, each worker finishes its `range jobs` loop after draining any buffered jobs. A worker must not close `results`: it cannot know whether the other workers still have results to send.

The coordinator is the one place that knows all workers have finished. It calls `wg.Wait()` and then closes `results` exactly once. The caller can therefore safely use `for result := range results`.

If multiple goroutines send to the same channel, they must coordinate channel closure. No individual sender should close the channel unless it can prove that every sender is finished.

### Interview checkpoints

1. Who closes `jobs`?

   The caller, because it owns sending jobs.

2. Who closes `results`?

   The pool coordinator, after `wg.Wait()` confirms that every worker has exited.

3. Why should a worker not close `results`?

   One worker cannot know whether another worker will send another result.

4. What happens when a worker ranges over closed `jobs`?

   It receives buffered jobs first and then exits the loop.

5. What happens if the consumer stops reading `results`?

   Workers can block while sending. They cannot exit, so the `WaitGroup` cannot finish and `results` cannot be closed. Worker Pool v1 documents this behavior; cancellation belongs to v2.

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
