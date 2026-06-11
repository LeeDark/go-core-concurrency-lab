# Concurrency & Parallelism

- Concurrency is not how the real world works. The real world works with parallelism.
- Parallelism is the natural way of thinking about multiple independent things interacting with each other.
- Solution of a Problem in Concurrent Programming Control, Edsger Dijkstra, 1965.
  - Concurrent programming & Parallel execution.
  - Concurrency relates to how programs are written. Parallelism relates to how programs run.
- The same concurrency principles apply to applications running on multiple machines on a distributed system, to applications running on a multi-core processor in a laptop, and to applications that run on a single-core system.
- Communicating Sequential Processes (CSP) is a mathematical model that influenced the design of Go.
  - Systems are composed of multiple sequential processes that are running in parallel.
  - These processes can communicate with each other synchronously, which means that a system sending a message to 
    the other can only continue once the other system receives it (this is exactly how unbuffered channels behave in Go.).
- Can the system deadlock?
  - When multiple systems run together, the possible states of the composite system grow exponentially.
  - Independent components of a concurrent program can run in any order, making it practically impossible to do
    state analysis.
  - Sequential process: the state can be defined as the values in memory together with the current execution location 
    of that program.
  - Multiple processes are running in parallel: the state of the whole system is the combination of the states of its 
    components.
- Concurrency is the ability of different parts of a program to be executed out-of-order or in partial order without affecting the result.
  - Concurrency is about how the program is written.
  - It is about "dealing with multiple things at the same time."
  - Doing multiple things at the same time defines parallelism.
- Concurrent programming is about organizing a problem into computational units that can run using time sharing or 
  that can run in parallel.
  - Time-sharing means sharing a computing resource with multiple users or processes.
  - Context-switching: when multiple threads of executions are created by a program, the processor runs one thread 
    for some time, and then switches to another thread, and so on.
- Concurrency is a programming model like object-oriented programming or functional programming.
  - Object-oriented programming divides a problem into logically related structural components that interact with each other.
  - Functional programming divides a problem into functional components that call each other.
  - Concurrent programming divides a problem into temporal components that send messages to each other, and that can be interleaved or run in parallel.
- Go runtime before 1.14 used cooperative scheduler, after it is using preemptive scheduler.
  - In preemptive threading, a running thread can be stopped at any time during that thread’s execution.
  - In non-preemptive threading (or cooperative threading), a running thread voluntarily gives up execution by performing a blocking operation, a system call, or something else.
- The states a thread/gorouine can be next:
  - **Ready** state: when a thread is **created**.
  - **Running** state: when the scheduler **assigns** it to a processor and starts running.
  - A **Running** thread can be **preempted** and moved back into the **Ready** state.
  - **Blocked** state: when the thread performs an I/O operation or blocks waiting for a lock or channel operation.
  - When the I/O operation completes, the lock is unlocked, or the channel operation is completed, the thread moves 
    back to the **Ready** state, waiting to be scheduled.
- 

# Goroutines

- A **process** is an instance of a program with certain dedicated resources, such as memory space, processor time, 
  file handles (for example, most processes in Linux have stdin, stdout, and stderr), and at least one thread.
  - Any two processes that wish to communicate have to do it through well-defined **inter-process communication** 
    utilities.
- A **thread** is an execution context that contains all the resources required to run a sequence of instructions.
  - A thread is usually managed by the operating system.
  - The stack is necessary to:
    - keep the sequence of nested function calls within that thread.
    - store values declared in the functions executing in that thread.
  - A given function may execute in many different threads, so the local variables used when that function runs in a thread are stored in the stack of that thread.
- A **scheduler** allocates processor time to threads.
  - Some schedulers are preemptive and can stop a thread at any time to switch to another thread.
  - Some schedulers are collaborative and have to wait for the thread to yield to switch to another one.
- A **goroutine** is an execution context that is managed by the Go runtime (as opposed to a thread that is managed by the operating system).
  - A goroutine usually has a much smaller startup overhead than an operating system thread. 
  - A goroutine starts with a small stack that grows as needed.
  - Creating new goroutines is faster and cheaper than creating operation system threads.
- The **Go scheduler** assigns operating system threads to run goroutines.
  - The number of operating system threads used by the Go runtime is equal to the number of processors/cores on the platform.
    - Unless you change this by setting the GOMAXPROCS environment variable or by calling the runtime.GOMAXPROCS function.
    - Anything more than that and the operating system will have to resort to time sharing.
    - With GOMAXPROCS threads running in parallel, there is no context-switching overhead at the operating system level.
  - The Go scheduler assigns goroutines to operating system threads to get more work on each thread as opposed to doing less work on many threads.
  - The Go scheduler performs better because it knows which goroutines to wake up to get more out of them.
- The **go** keyword starts the given function in a new goroutine.
  - The function running as a goroutine can take parameters, but it cannot return a value.
- Differences between threads and goroutines
  - Threads usually have priorities. Goroutines do not have pre-assigned priorities.
  - A goroutine starts with a small stack (Go runtimes after 1.19 use a historical average, earlier versions use 2K). 
    - Every function call checks whether the remaining stack space is sufficient. If not, the stack is resized.
  - An operation system thread usually starts with a much larger stack (in the order of megabytes) that usually does not grow.
- The **Go runtime** starts several goroutines when a program starts.
  - At least one for the **garbage collector** and another for the **main** goroutine.
- Data race — concurrent access to the same memory location from multiple goroutines, where at least one access is a write, and the accesses are not properly synchronized.
  - Closure — a function that captures and uses variables from its surrounding lexical scope.
    - A common source of bugs is combining closures with goroutines and shared mutable state.
  - Scope — the part of the source code where an identifier is visible and can be used.
    - Scope describes visibility in code; it does not define whether a variable is allocated on the stack or on the heap.
  - Stack — memory used for function call frames, local variables, and execution state. Each goroutine in Go has its own stack.
    - Goroutine stacks start small and can grow dynamically, which makes goroutines lightweight.
    - The compiler determines actual variable placement through escape analysis.
  - Heap — memory used for values that need to outlive the current stack frame or cannot be safely kept on the stack.
    - Heap allocations are managed by the Garbage Collector.
  - Escape analysis — a compiler analysis that determines whether a value can be allocated on the stack or must escape to the heap.
- Problems with goroutines
  - The ability to create concurrent execution blocks.
  - How to terminate them responsibly?
- How Go runtime manages goroutines
  - Go uses an M:N scheduler that runs M goroutines on N OS threads.
  - Go runtime keeps track of the OS threads and the goroutines.
  - When an OS thread is ready to execute a goroutine, the scheduler selects one that is ready to run and assigns it to the thread.
  - The OS thread runs that goroutine until it blocks, yields, or is preempted.
    - Blocking by channel operations or mutexes is managed by the Go runtime.
    - If the goroutine is blocked because of a synchronous I/O operation, then the thread running that goroutine will also be blocked (this is managed by the operating system).
      - In this case, the Go runtime starts a new thread or uses one already available and continues operation.
      - When the OS thread unblocks (that is, the I/O operation ends), the thread is put back into use or returned to the thread pool.
    - The Go runtime limits the number of active OS threads running user goroutines with the GOMAXPROCS variable.
      - However, there is no limit on the number of OS threads waiting for I/O operations.
      - So, the actual OS thread count a Go program uses can be much higher than GOMAXPROCS.
      - However, only GOMAXPROCS of those threads would be executing user goroutines.
    - A similar process happens for asynchronous I/O operations, such as network operations and some file operations on certain platforms.
      - However, instead of blocking a thread for a system call, the goroutine is blocked, and a netpoller thread is used to wait for asynchronous events.
      - When the netpoller receives events, it wakes up the relevant goroutines.
- In Go, most I/O operations
  - Are exposed as synchronous blocking calls from the goroutine's point of view. For example, file reads, network reads, HTTP calls, and database queries block the current goroutine until the operation completes.
  - Asynchronous I/O is usually modeled by running blocking operations in separate goroutines and communicating results through channels, WaitGroup, errgroup, or context cancellation.
  - For network I/O, Go runtime uses a netpoller internally, so blocking network calls usually park the goroutine instead of blocking the underlying OS thread.
  - File I/O is generally closer to real blocking I/O.


## Definition
A goroutine is a lightweight concurrent execution unit managed by the Go runtime.
It is started with the `go` keyword and runs independently from the caller.

## Cost
Goroutines are much cheaper than OS threads because they start with a small stack that can grow dynamically. 
However, they are not free: too many goroutines can still consume memory, scheduler time, and other resources.

## Lifecycle
A goroutine starts when the `go` statement is executed and stops when its function returns.
If the `main` goroutine exits, the whole program exits, even if other goroutines are still running.

## Leaks
A goroutine leak happens when a goroutine never finishes, usually because it is blocked on a channel, I/O operation, lock, or missing cancellation signal.
Leaked goroutines keep memory and resources alive, which can slowly degrade the application.

## WaitGroup
`sync.WaitGroup` is used to wait for a group of goroutines to finish. It coordinates completion, but it does not cancel goroutines or handle errors by itself.

## Context
`context.Context` is used to propagate cancellation, deadlines, timeouts, and request-scoped data across goroutines and API boundaries.
It is commonly used to stop goroutines gracefully when their work is no longer needed.

## Communication
Goroutines usually communicate through channels or shared memory protected by synchronization primitives such as mutexes.
In Go, the preferred style is often “share memory by communicating”, but shared memory is still common and must be synchronized correctly.


# Channels

- The operating system does not know about channel operations or mutexes, which are all managed in the user space by the Go runtime.
- Channels allow goroutines to share memory by communicating, as opposed to communicating by sharing memory.
- When you are working with channels, you have to keep in mind that channels are two things combined together:
  - they are synchronization tools.
  - they are conduits for data. A channel is a first-in, first-out (FIFO) conduit.
- A send operation to a channel will block until the channel is ready to accept a value.
- A receive operation from a channel will block until the channel is ready to provide a value.
- A channel is actually a pointer to a data structure that contains its internal state, so the zero-value of a channel variable is nil.
  - Reading from or writing to a nil channel will block indefinitely.
- The Go garbage collector will collect channels that are no longer in use.
  - If there are no goroutines that directly or indirectly reference a channel, the channel will be garbage collected even if its buffer has elements in it.
  - You do not need to close channels to make them eligible for garbage collection.
  - In fact, closing a channel has more significance than just cleaning up resources.
- Senders/Writers and Receivers/Readers
  - A channel can have multiple receivers, but each sent value is received by only one receiver.
  - This makes a single channel suitable for work distribution, such as worker pools.
  - A normal send is not a broadcast: one value is not delivered to all receivers.
  - Closing a channel is a broadcast signal: all current and future receivers are unblocked.
  - `context.Context` uses this pattern: cancellation closes the `Done()` channel, notifying all goroutines waiting on it.
- Closed channel
  - A receive operation from a closed channel will always succeed with the zero value of the channel type.
  - Writing to a closed channel will always panic.
  - for a receiver, it is usually important to know whether the channel was closed when the read happened
- Unbuffered channel is a channel created without a buffer
  - Unbuffered channel acts as a synchronization point between two goroutines
  - A send operation will block until another goroutine receives from it.
  - A receive operation will block until another goroutine sends to it.
  - A way to transfer data between goroutines atomically.
- A channel can be declared with a direction
  - `var receiveOnly <-chan int   // Can receive, cannot write or close`
  - `var sendOnly chan<- int      // Can send, cannot read or close`
- When multiple goroutines attempt to send to a channel or when multiple goroutines attempt to read from a channel, they are scheduled randomly.
- Worker pool patterns
  - You can create many worker goroutines, all receiving from a channel.
  - Another goroutine sends work items to the channel, and each work item will be picked up by an available worker goroutine and processed.
  - Then, you can have one goroutine reading from a channel that is written by many worker goroutines. 
  - The reading goroutine will collect the results of computations performed by those goroutines.
- A select statement chooses which of a set of possible send or receive operations proceed.
  - At a high level, the select statement chooses one of the send or receive operations that can proceed and then runs the block corresponding to the chosen operation.
  - If there are multiple send or receive operations that can proceed, the select statement chooses one randomly.
  - If there are none, the select statement chooses the default option.
  - If a default option does not exist, the select statement blocks until one of the channel operations becomes available.
  - Using the default option in a select statement is useful for non-blocking sends and receives.
- Channels can be used to gracefully terminate a program based on a signal.

## Unbuffered
An unbuffered channel has no capacity, so send and receive operations must meet at the same time.
A send blocks until another goroutine receives the value, and a receive blocks until another goroutine sends a value.

## Buffered
A buffered channel has a fixed capacity and allows sends to proceed until the buffer is full.
Sends block when the buffer is full, and receives block when the buffer is empty.

## Close
Closing a channel signals that no more values will be sent on it.
Receivers can still read remaining buffered values, and then receive the zero value with `ok == false`.

## nil channel
A nil channel blocks forever on both send and receive operations.
In `select`, a nil channel case is disabled, which can be useful for dynamically turning cases on or off.

## Directional channels
Directional channels restrict how a channel can be used: send-only `chan<- T` or receive-only `<-chan T`.
They are often used in function signatures to make ownership and intent clearer.

## Common mistakes
Common mistakes include:
- sending on a closed channel
- forgetting to close a channel when receivers depend on it
- closing a channel from the receiver side
- leaking goroutines blocked on channels
- using `time.Sleep` instead of proper synchronization.
Another frequent mistake is assuming that buffered channels remove the need for synchronization; they only change blocking behavior.

## Ownership
Channel ownership means being responsible for creating the channel, sending values, and closing it.
In Go, the sender usually owns the channel closure, because only the sender knows when no more values will be sent.

## Common deadlocks
A deadlock happens when goroutines are blocked forever and no goroutine can make progress.
Common cases include:
- sending to an unbuffered channel without a receiver
- receiving from a channel without a sender
- waiting on a `WaitGroup` that is never completed
- ranging over a channel that is never closed.

## Select
`select` waits on multiple channel operations and runs one case that becomes ready.
If multiple cases are ready, one is chosen pseudo-randomly; if none are ready, `default` runs if present, otherwise `select` blocks.

### Multiplexing
`select` allows a goroutine to wait on multiple channel operations at the same time.
It is commonly used to combine results from several channels or to react to whichever operation becomes ready first.

### Timeout
Timeouts are usually implemented with `select` and `time.After` or a context with a deadline.
This prevents a goroutine from waiting forever for a channel operation.

### Cancellation
Cancellation is commonly handled by listening to `ctx.Done()` inside a `select`.
This allows a goroutine to stop its work when the request is cancelled, a timeout expires, or the parent operation is no longer needed.

### Default case risks
A `default` case makes `select` non-blocking, which can be useful but also dangerous.
If used carelessly, it can create busy loops, high CPU usage, missed backpressure, or logic that silently skips work instead of waiting properly.


# Worker Pool

Context
- cancellation;
- timeout;
- deadline;
- values: осторожно.

## Pattern
## When to use
## When not to use
## Backpressure
## Cancellation
## Error handling
## Testing

# Race Detector / Memory Model / Scheduler

## Data race
## Happens-before
## Mutex
## Atomic
## GOMAXPROCS
## CPU-bound vs I/O-bound

# Interview

```text
1. Что такое goroutine и чем она отличается от OS thread?
2. Как остановить goroutine?
3. Что такое goroutine leak?
4. Чем buffered channel отличается от unbuffered?
5. Кто должен закрывать channel?
6. Что произойдёт при чтении из закрытого channel?
7. Что произойдёт при записи в закрытый channel?
8. Как работает select?
9. Для чего нужен default в select и чем он опасен?
10. Как реализовать timeout?
11. Как реализовать cancellation?
12. Чем context cancellation отличается от closing done channel?
13. Что такое worker pool?
14. Что такое fan-out/fan-in?
15. Когда использовать mutex вместо channel?
16. Когда использовать atomic?
17. Что находит race detector?
18. Что race detector не гарантирует?
19. Как сделать graceful shutdown HTTP-сервиса?
20. Как диагностировать блокировки и latency в concurrent Go-приложении?
```

```text
1. Что такое goroutine?
2. Как остановить goroutine?
3. Что такое goroutine leak?
4. Buffered vs unbuffered channel?
5. Кто закрывает channel?
6. Что будет при чтении из закрытого channel?
7. Что будет при записи в закрытый channel?
8. Как работает select?
9. Как сделать timeout?
10. Как работает context cancellation?
11. Что такое worker pool?
12. Что такое fan-out/fan-in?
13. Когда использовать mutex вместо channel?
14. Что проверяет race detector?
15. Как избежать unbounded concurrency?
```