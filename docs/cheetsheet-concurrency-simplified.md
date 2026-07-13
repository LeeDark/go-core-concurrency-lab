# Concurrency & Parallelism

**Concurrency** is a way to structure a program as independent parts that coordinate with each other. Those parts may be interleaved on one CPU or run at the same time on several CPUs.

**Parallelism** is executing multiple pieces of work at the same time. It depends on available CPUs, the Go runtime, `GOMAXPROCS`, blocking, and synchronization.

In short: **concurrency is about program design; parallelism is about execution.** A concurrent program is not automatically faster than a sequential one.

The coffee-shop analogy:

- Concurrency: split work into taking orders, grinding coffee, and brewing coffee.
- Parallelism: add another barista and machine so two drinks are made at once.

Concurrency applies on a single-core machine, a multi-core machine, or across distributed systems. It helps model independent activities, but it also increases the number of possible program states and makes deadlocks and races easier to create.

Go is influenced by Communicating Sequential Processes (CSP): independent sequential processes communicate by passing messages. An unbuffered channel is a synchronous handoff: the sender continues only when a receiver accepts the value.

## Performance and scheduler

Goroutines are cheaper than OS threads, but they still cost memory, scheduling time, synchronization, and communication overhead. Concurrency is useful for independent CPU work, I/O, pipelines, worker pools, and request handling. It is usually a poor fit for tiny tasks, strictly sequential work, or work with heavy shared-state coordination.

The Go scheduler uses:

```text
G — goroutine
M — OS thread
P — logical processor used by the Go runtime
```

`GOMAXPROCS` limits how many OS threads can execute Go code simultaneously. It usually defaults to the number of available CPUs. The runtime may still create more OS threads, for example when threads are blocked in system calls.

A goroutine is usually:

```text
running   — currently executing
runnable  — ready to execute, waiting for CPU time
waiting   — blocked on I/O, a channel, a mutex, a syscall, or another event
```

The runtime uses local and global run queues plus work stealing. Since Go 1.14, it can preempt long-running goroutines so other goroutines can run. Before that, scheduling was more cooperative.

# Goroutines

A **process** is an operating-system-managed program instance with resources such as memory, file handles, and one or more threads. Processes communicate through inter-process communication mechanisms.

A **thread** is an OS-managed execution context. Its stack holds function calls, local values, and execution state. The OS scheduler assigns CPU time to threads.

A **goroutine** is a lightweight execution unit managed by the Go runtime. Start one with `go`:

```go
go work()
```

It can receive arguments, but it does not return a value directly to the caller. Goroutines start with a small, dynamically growing stack, so they are cheaper than OS threads; they are not free and should not be created without a lifecycle plan.

## Lifecycle and leaks

A goroutine starts when its `go` statement runs and ends when its function returns. If `main` returns, the entire program exits, including unfinished goroutines.

A **goroutine leak** occurs when a goroutine never finishes—for example, it is blocked forever on a channel, lock, I/O operation, or missing cancellation signal. Leaks retain memory and resources.

`sync.WaitGroup` waits for a set of goroutines to finish. It coordinates completion only: it does not cancel work or collect errors.

`context.Context` carries cancellation, deadlines, timeouts, and request-scoped values across API boundaries. It is commonly used to let goroutines stop when their work is no longer needed.

## Communication and shared state

Goroutines communicate through channels or shared memory protected by synchronization primitives such as mutexes. “Share memory by communicating” is a useful Go guideline, but shared memory is also common and must be synchronized.

A **data race** is concurrent access to the same memory location where at least one access is a write and the accesses are not properly synchronized.

A closure captures variables from its surrounding scope. Combining closures, goroutines, and shared mutable variables is a common source of races. Each goroutine has its own stack; values that must outlive a stack frame may be allocated on the heap, as decided by compiler escape analysis.

For network I/O, Go commonly parks the goroutine and uses its netpoller internally. File I/O is generally closer to blocking the underlying OS thread. In both cases, I/O appears synchronous to the goroutine.

# Channels

Channels are Go-runtime-managed conduits for data and synchronization. A channel delivers values in FIFO order. A send blocks until the channel can accept a value; a receive blocks until a value is available.

## Unbuffered and buffered channels

An **unbuffered** channel has no capacity. Send and receive must meet: each blocks until the other is ready. This makes it a synchronization point and an atomic handoff.

A **buffered** channel has a fixed capacity. Sends proceed until the buffer is full; receives proceed while it contains values. A buffer changes blocking behavior, but it does not remove the need for synchronization or ownership rules.

## Close and nil channels

Closing a channel means no more values will be sent. Receivers first drain buffered values; later receives return the zero value and `ok == false`:

```go
value, ok := <-ch
```

Sending to a closed channel panics. Closing is not required for garbage collection; it is a signal to receivers. A normal send goes to one receiver, while closing unblocks all receivers waiting on that channel.

The zero value of a channel is `nil`. Sending to or receiving from a nil channel blocks forever. In a `select`, a nil-channel case is disabled.

## Ownership and direction

The sender usually owns channel closure because it knows when no more values will be sent. Receivers should not close a channel they do not own.

Directional channels clarify intent in APIs:

```go
var receiveOnly <-chan int // receive only
var sendOnly chan<- int    // send only
```

Multiple receivers may read from one channel, but each value is received by only one of them. This is useful for distributing jobs among workers. When several operations are ready, `select` chooses one pseudo-randomly; it blocks when none are ready unless it has a `default` case.

Use `default` for intentionally non-blocking operations. Used carelessly, it can cause busy loops, high CPU usage, skipped work, or lost backpressure.

`select` also multiplexes channels: one goroutine can wait for whichever operation becomes ready first. A timeout commonly selects between the operation and `time.After` (or, preferably for request work, a context deadline). Cancellation commonly selects on `ctx.Done()` so the goroutine can stop when its parent operation is cancelled.

## Common mistakes

- Sending to a closed channel.
- Closing a channel from the receiver side.
- Ranging over a channel that no sender will close.
- Sending to an unbuffered channel without a receiver, or receiving without a sender.
- Waiting on a `WaitGroup` whose counter never reaches zero.
- Using `time.Sleep` instead of synchronization.

# Worker Pool

Worker Pool v1 uses a fixed number of goroutines to process a stream of jobs:

```text
producer -> jobs channel -> N workers -> results channel -> consumer
```

It provides **bounded concurrency**: at most `workerCount` jobs are handled at once. This avoids creating one goroutine per job when throughput or downstream resources must be controlled.

## API

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

`Run` treats a non-positive `workerCount` as one worker. It starts the workers and returns the results channel immediately.

## Lifecycle and ownership

- The caller sends jobs and closes `jobs`.
- Each worker ranges over `jobs`, handles a job, and sends one `Result`.
- Closing `jobs` lets workers finish after any buffered jobs are received.
- Workers send results; the caller receives them and ranges until `results` is closed.
- A separate coordinator goroutine waits on a `sync.WaitGroup` and closes `results` after every worker has exited.
- Workers must not close `results`: no individual worker knows whether the other workers have finished sending.

The intended consumer flow is:

```text
create jobs -> start Run -> send jobs -> close jobs -> range over results
```

## Backpressure and errors

`results` is unbuffered in v1. If the consumer stops reading it, workers block while sending, cannot finish, and `results` cannot be closed. The consumer must drain results unless a later version adds cancellation.

Job-level failures are returned in `Result.Err`; v1 does not define advanced error handling or stop the whole pool on one job error.

## Scope of v1

V1 deliberately focuses on goroutines, channels, ownership, a fixed worker count, and `WaitGroup`. It does **not** include `context.Context`, cancellation, deadlines, timeouts, graceful shutdown, goroutine-leak checks, or advanced error policy. Those belong to Worker Pool v2.

## Tests

The current tests verify that the pool processes all jobs, preserves job errors, closes `results`, and normalizes a non-positive worker count.

```bash
go test ./06-worker-pool-v1/workerpool
go test -race ./06-worker-pool-v1/workerpool
```

# Race Detector / Memory Model / Scheduler

## Data race
## Happens-before
## Mutex
## Atomic
## GOMAXPROCS
## CPU-bound vs I/O-bound

# Interview
