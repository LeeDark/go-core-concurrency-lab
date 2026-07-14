# Concurrency & Parallelism

## Definitions

**Concurrency** is the ability to structure a program as independent computational units that can make progress out of order or in a partial order without changing the intended result. The units may be interleaved through time-sharing or may run in parallel.

**Parallelism** means executing multiple pieces of work at the same time.

In short: **concurrency is about how a program is designed; parallelism is about how it runs.** This distinction appears in early work such as Edsger Dijkstra's *Solution of a Problem in Concurrent Programming Control* (1965): concurrent programming and parallel execution are related, but they are not the same thing.

The coffee-shop analogy:

- Concurrency: split the process into taking orders, grinding coffee, and brewing coffee so the stages can coordinate.
- Parallelism: add another barista and coffee machine so two customers can be served simultaneously.

The same concurrency principles apply to a program on a single-core machine, a multi-core machine, or multiple machines in a distributed system.

## Concurrency as a programming model

Programming models divide problems in different ways:

- Object-oriented programming divides a problem into structural components that interact.
- Functional programming divides a problem into functional components that call each other.
- Concurrent programming divides a problem into temporal components that communicate and may be interleaved or run in parallel.

Communicating Sequential Processes (CSP) is a mathematical model that influenced Go. A system is composed of independent sequential processes that communicate by passing messages. Communication may be synchronous: a sender cannot continue until a receiver accepts the message. This is how an unbuffered Go channel behaves.

## Complexity and deadlocks

A sequential program state consists roughly of its memory values and current execution location. In a concurrent program, the whole state is the combination of the states of all components.

Independent components may run in many different orders, so the number of possible combined states grows rapidly. This makes exhaustive state analysis difficult and creates room for deadlocks, data races, and ordering bugs.

## Performance

Concurrent code is not automatically faster than sequential code. Actual performance depends on whether work can run in parallel, available CPUs, blocking operations, `GOMAXPROCS`, and overhead from goroutines, scheduling, synchronization, channels, mutexes, and contention.

Goroutines are cheaper than OS threads, but they are not free. Creating, scheduling, switching, synchronizing, and communicating between them all have a cost.

Good candidates for concurrency:

- independent CPU-heavy chunks;
- independent I/O operations;
- pipelines;
- worker pools;
- request handling.

Poor candidates:

- tiny tasks where coordination costs more than the work;
- strongly sequential logic;
- tasks with heavy shared state;
- tasks that require too much synchronization.

## Go scheduler

The Go scheduler is commonly described with the G-M-P model:

```text
G — goroutine
M — OS thread (machine)
P — logical processor required to execute Go code
```

Go uses an M:N scheduler: many goroutines are multiplexed over a smaller or equal set of OS threads. A `P` holds scheduler resources and a runnable queue. The runtime uses local queues, a global queue, and work stealing to distribute runnable goroutines.

`GOMAXPROCS` limits the number of `P`s and therefore how many OS threads can execute user-level Go code simultaneously. It usually defaults to the number of available CPUs. It does **not** limit the total number of OS threads: the runtime may create more threads when some are blocked in system calls.

A goroutine is usually in one of these states:

```text
running   — currently executing
runnable  — ready to execute, waiting for CPU time
waiting   — blocked on I/O, a channel, a mutex, a syscall, or another event
```

A running goroutine may be preempted and moved back to runnable. A waiting goroutine becomes runnable when the event it waits for completes.

Since Go 1.14, the scheduler supports asynchronous preemption, so long-running goroutines can be interrupted to let other goroutines run. Earlier scheduling relied more heavily on cooperative safe points and blocking operations.

# Goroutines

## Process, thread, and goroutine

A **process** is an instance of a program with operating-system-managed resources such as memory, processor time, file handles, and at least one thread. Separate processes communicate through defined inter-process communication mechanisms.

A **thread** is an OS-managed execution context. Its stack records nested function calls, local values, and execution state. The OS scheduler assigns processor time to threads and may preempt one thread to run another.

A **goroutine** is a lightweight concurrent execution unit managed by the Go runtime. The `go` keyword starts a function in a new goroutine:

```go
go work()
```

The function may accept arguments, but it cannot return a value directly to the caller. Results must be communicated through a channel, shared state, a callback, or another coordination mechanism.

## Goroutines vs OS threads

- Goroutines are scheduled by the Go runtime; threads are scheduled by the operating system.
- Goroutines start with small stacks that grow and shrink as needed; OS thread stacks are usually much larger.
- Goroutines have lower creation and context-switching overhead, but they still consume memory and scheduler time.
- OS threads may have priorities; goroutines do not have user-assigned priorities.
- Many goroutines can run on a smaller set of OS threads.

Each function call checks whether the goroutine stack has enough space. The runtime grows the stack when necessary. The exact initial size and growth strategy are implementation details and may change between Go versions.

## Lifecycle

A goroutine starts when its `go` statement executes and finishes when its function returns. If the `main` goroutine returns, the whole program exits, even if other goroutines are unfinished.

Creating a goroutine also creates a lifecycle question: how will it finish? A goroutine should normally have a clear completion condition, input closure, or cancellation signal.

The Go runtime itself starts internal goroutines for work such as garbage collection and runtime management in addition to the `main` goroutine.

## How the runtime schedules goroutines

The runtime tracks goroutines, OS threads, and logical processors. A runnable goroutine is selected and assigned to an OS thread associated with a `P`. It runs until it finishes, blocks, yields, or is preempted.

Blocking behavior depends on the operation:

- Channel and mutex blocking is known to the Go runtime, so it can park only the goroutine and run another one.
- A blocking system call may block the OS thread. The runtime can detach its `P` and use or create another thread to keep Go code running.
- The total OS thread count can therefore be higher than `GOMAXPROCS`, although only `GOMAXPROCS` threads execute user-level Go code at once.
- Network I/O commonly uses the runtime netpoller. The goroutine is parked while the netpoller waits for readiness, so the OS thread usually remains available.
- File I/O is generally closer to a blocking system call and may block its OS thread.

Most Go I/O APIs still look synchronous from the goroutine's point of view: a file read, network read, HTTP call, or database query returns only after it completes. Concurrent I/O is usually expressed by running such calls in goroutines and coordinating them with channels, `WaitGroup`, `errgroup`, or context cancellation.

## Stack, heap, closures, and races

- **Scope** is the source-code region where an identifier is visible. It does not determine whether a value lives on the stack or heap.
- **Stack** memory stores call frames, local values, and execution state. Every goroutine has its own stack.
- **Heap** memory stores values that cannot safely remain in a stack frame or must outlive it. Heap allocations are managed by the garbage collector.
- **Escape analysis** is the compiler analysis that decides whether a value can stay on the stack or must escape to the heap.
- A **closure** captures variables from its surrounding lexical scope. Closures combined with goroutines and shared mutable variables are a common source of bugs.
- A **data race** is concurrent access to the same memory location where at least one access is a write and the accesses are not properly synchronized.

## Leaks and coordination

A **goroutine leak** happens when a goroutine never finishes, commonly because it is blocked on a channel, I/O operation, lock, or missing cancellation signal. Leaked goroutines keep memory and other resources alive and may slowly degrade a program.

`sync.WaitGroup` waits for a group of goroutines to finish. It coordinates completion, but does not cancel goroutines or handle their errors.

`context.Context` propagates cancellation, deadlines, timeouts, and request-scoped values across API boundaries. It is commonly used to let goroutines stop when their work is no longer needed.

Goroutines communicate through channels or through shared memory protected by synchronization primitives such as mutexes. “Share memory by communicating” is a useful Go guideline, but synchronized shared memory is also common and sometimes simpler.

# Channels

## Purpose and blocking behavior

Channels are managed by the Go runtime, not by the operating system. They combine two roles:

- a conduit for typed data;
- a synchronization mechanism between goroutines.

A channel preserves the order of values sent by a single sender. A send blocks until the channel can accept a value, and a receive blocks until a value is available.

## Unbuffered channels

An unbuffered channel has no storage capacity. A send blocks until another goroutine receives the value, and a receive blocks until another goroutine sends one.

This makes an unbuffered channel a synchronization point and an atomic handoff between sender and receiver: both sides must participate before either operation completes.

## Buffered channels

A buffered channel has a fixed capacity:

- a send blocks when the buffer is full;
- a receive blocks when the buffer is empty;
- otherwise, the operation can proceed without waiting for the other side at that exact moment.

A buffer changes blocking behavior, but does not remove the need for synchronization, ownership, or backpressure reasoning.

## Nil channels and garbage collection

The zero value of a channel is `nil`. Sending to or receiving from a nil channel blocks forever. In a `select`, a case using a nil channel is disabled, which can be useful for dynamically turning cases on and off.

Channels do not need to be closed for garbage collection. If no reachable goroutine or value references a channel, it can be collected even when its buffer contains values. Closing has semantic meaning—it signals that no more values will be sent—not cleanup meaning.

## Closing a channel

Closing a channel signals that no more values will be sent.

- Receivers first drain any buffered values.
- Later receives complete immediately with the element type's zero value and `ok == false`.
- Sending to a closed channel panics.
- Closing an already closed channel panics.

Use the comma-ok form when a receiver must distinguish a real zero value from a closed channel:

```go
value, ok := <-ch
```

A normal send is not a broadcast: one sent value is received by one receiver. Closing behaves as a notification to all current and future receivers because receives no longer block. `context.Context` uses this pattern by closing `Done()` on cancellation.

## Ownership

Channel ownership means deciding who creates a channel, who sends and receives, and who closes it. The sending side usually owns closure because only senders know when no more values will be produced.

A channel may have multiple senders or receivers. Each sent value is consumed by only one receiver, making a shared input channel useful for distributing work among worker goroutines. If multiple goroutines can send, they must coordinate closure so nobody sends after the channel is closed.

Directional channel types make intent explicit in APIs:

```go
var receiveOnly <-chan int // receive only
var sendOnly chan<- int    // send only
```

A receive-only reference cannot receive permission to send or close. A send-only reference cannot receive; it may close the channel because closing belongs to the sending side.

## Select

`select` waits on a set of channel sends and receives and executes one operation that can proceed.

- If one case is ready, it runs.
- If multiple cases are ready, one is selected pseudo-randomly.
- If no case is ready and `default` exists, `default` runs.
- If no case is ready and there is no `default`, `select` blocks.

### Multiplexing

`select` lets one goroutine wait on several channel operations and respond to whichever becomes ready first. It is commonly used to merge events, wait for a result or cancellation, or coordinate multiple pipeline stages.

### Timeout

A timeout can be implemented by selecting between an operation and `time.After`, or by using a context with a deadline. This prevents waiting forever for a channel operation.

### Cancellation

Cancellation is commonly handled by selecting on `ctx.Done()`. When the context is cancelled, its `Done` channel closes and all goroutines waiting on it can stop.

### Default case risks

A `default` case makes `select` non-blocking. This is useful when skipping an unavailable operation is intentional, but careless use can create busy loops, high CPU usage, missed backpressure, or silently dropped work.

## Common mistakes and deadlocks

A deadlock occurs when goroutines are blocked forever and no goroutine can make progress. Common channel-related mistakes include:

- sending to a closed channel;
- closing a channel from the receiver side;
- closing a channel while other senders may still use it;
- sending to an unbuffered channel without a receiver;
- receiving from a channel without a sender;
- ranging over a channel that is never closed;
- leaking goroutines blocked on channels;
- waiting on a `WaitGroup` whose counter never reaches zero;
- using `time.Sleep` instead of proper synchronization.

# Worker Pool

## Pattern

Worker Pool v1 uses a fixed number of goroutines to process a stream of jobs:

```text
producer -> jobs channel -> N workers -> results channel -> consumer
```

Instead of creating one goroutine for every job, the pool starts `workerCount` goroutines. Each worker repeatedly receives a job, calls the handler, sends its result, and exits when `jobs` is closed and drained.

This is **bounded concurrency**: at most `workerCount` handlers run simultaneously.

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

The directional types document intent: the pool only receives from `jobs`, and the caller only receives from the returned results channel. The current implementation normalizes a non-positive `workerCount` to one.

## Lifecycle and channel ownership

`jobs` ownership:

- the caller sends jobs;
- the caller closes `jobs`;
- workers receive jobs;
- workers never close `jobs`.

`results` ownership:

- workers send results;
- the caller receives results;
- the pool closes `results`;
- the caller and individual workers never close `results`.

The pool uses `sync.WaitGroup` to close `results` safely:

```text
Add before starting workers
Done when each worker exits
Wait in a separate coordinator goroutine
close(results) after Wait returns
```

No individual worker can close `results`, because it does not know whether the other workers have finished sending.

`results` is not closed merely because `jobs` is closed. It stays open until every worker has exited. The coordinator closes it only after `wg.Wait()` returns.

The normal flow is:

```text
create jobs channel
start the pool
send jobs
close jobs
range over results
```

## When to use

- jobs are independent;
- concurrency must be bounded;
- work arrives as a stream;
- throughput or pressure on a dependency must be controlled;
- results can be collected asynchronously.

## When not to use

- work is small and strictly sequential;
- jobs depend heavily on one another;
- one goroutine is sufficient;
- a semaphore around existing goroutines expresses the limit more clearly;
- durable queuing, retries, or distributed processing are required.

## Backpressure

Both channels are unbuffered in v1. Producers wait until workers can receive jobs, and workers wait until the consumer can receive results.

If nobody reads `results`, workers block while sending. They cannot exit, the `WaitGroup` cannot finish, and the coordinator cannot close `results`. The consumer must drain the channel unless a later version adds cancellation.

## Error handling

Each job may return an error in `Result.Err`. V1 preserves job-level errors but does not cancel the pool, retry work, or define a global error policy.

## Scope of v1

V1 focuses on goroutines, channel ownership, a fixed worker count, `sync.WaitGroup`, and correct result-channel closure. It deliberately excludes `context.Context`, cancellation, deadlines, timeouts, graceful shutdown, goroutine-leak checks, and advanced error handling. Those topics belong to Worker Pool v2.

## Testing

The tests verify:

- all submitted jobs are processed;
- job errors remain in `Result.Err`;
- `results` closes after all workers finish;
- `results` remains open while a worker is still handling a job;
- a non-positive worker count becomes one worker.

Use targeted commands for this lab:

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
