# Roadmap

This roadmap records the complete study path and its current priorities. The goal is to build
understanding in small steps: read the notes, write focused examples, explain the behavior, then
move to the next lab.

## Study Flow

The complete study path is:

1. Project structure, modules, packages, visibility.
2. Go Core.
3. Context and errors.
4. Structs, methods, receivers, interfaces.
5. Slices, maps, defer.
6. Generics, tooling, workspaces.
7. Worker Pool v1.
8. Worker Pool v2.
9. Pipeline v1.
10. Pipeline v2.
11. Shared state: mutex vs channel vs atomic.
12. Race detector and Go memory model basics.
13. Book track: *Learning Functional Programming in Go* (Sheehan, 2017), with discussion of modern generics.

The study plan pairs core topics with concurrency labs:

| Core topic                            | Concurrency topic      |
|---------------------------------------|------------------------|
| Step 1: project structure             | Step 7: Worker Pool v1 |
| Step 2: Go Core                       | Step 8: Worker Pool v2 |
| Step 3: errors, context               | Step 9: Pipeline v1    |
| Step 4: structs, methods, interfaces  | Step 10: Pipeline v2   |
| Step 5: slices, maps, defer           | Step 11: shared state  |
| Step 6: generics, tooling, workspaces | Step 12: race detector |

The current priorities group that path into phases:

| Phase | Topic                                           | Status                   |
|-------|-------------------------------------------------|--------------------------|
| 1     | Slices                                          | Finished, closed         |
| 2     | Worker Pool v1                                  | Finished, closed         |
| 3     | Maps                                            | Current                  |
| 4     | Defer, errors, context, and Worker Pool v2      | Planned                  |
| 5     | Race detector, memory model, and runtime        | Planned                  |
| 6     | Types, interfaces, and Pipeline v1              | Planned                  |
| 7     | Structure and Pipeline v2                       | Planned                  |
| 8     | Generics, tooling, workspaces, and shared state | Planned                  |
| 9     | Go Core                                         | Planned, lowest priority |

## Phase 1: Slices

Finished and closed.

Focus:

- explain slice headers, backing arrays, `len`, and `cap`;
- understand when `append` reuses the backing array and when it allocates;
- avoid accidental aliasing with subslices;
- copy slices intentionally;
- distinguish nil slices from empty slices;
- understand why subslices can retain large backing arrays;
- remember that slices are not safe for concurrent mutation without synchronization.

Primary docs:

- [`04-slices-maps-defer/README.md`](04-slices-maps-defer/README.md)
- [`docs/cheatsheet-core.md`](docs/cheatsheet-core.md)

Review checklist:

- explain `len` vs `cap`;
- explain `append` with and without spare capacity;
- show a subslice aliasing bug;
- fix aliasing with `copy` or a full slice expression;
- explain nil vs empty slice behavior;
- explain the large backing-array retention problem.

## Phase 2: Worker Pool v1

Finished and closed.

Focus:

- start a fixed number of workers;
- send work through a `jobs` channel;
- collect output through a `results` channel;
- use `sync.WaitGroup` to wait for workers;
- close `results` only after all workers finish;
- explain channel ownership clearly.

Primary docs:

- [`06-worker-pool-v1/README.md`](06-worker-pool-v1/README.md)
- [`docs/cheatsheet-concurrency.md`](docs/cheatsheet-concurrency.md)

Stop line:

- Worker Pool v1 should stay minimal. Do not add context cancellation, timeouts, leak checks,
  graceful shutdown, or advanced error policy here.

## Phase 3: Maps

Current work.

Focus:

- understand map creation, lookup, insertion, update, and deletion;
- use the comma-ok form for lookups;
- distinguish nil maps from empty maps;
- understand that map iteration order is not stable;
- avoid concurrent map reads/writes without synchronization;
- use maps for counting, grouping, indexing, and set-like behavior.

Planned output:

- extend [`04-slices-maps-defer/README.md`](04-slices-maps-defer/README.md);
- add focused map examples and tests under `04-slices-maps-defer/`.

## Phase 4: Defer, Errors, Context, And Worker Pool v2

Core focus:

- use `defer` for cleanup without hiding control flow;
- understand defer argument evaluation and LIFO execution;
- return errors explicitly;
- wrap errors with useful context;
- distinguish sentinel errors, typed errors, and error inspection;
- understand `context.Context` as cancellation and deadline propagation.

Concurrency focus:

- add `context.Context` to the worker-pool API;
- stop workers while waiting for jobs;
- stop workers while sending results;
- define a simple cancellation and timeout policy;
- reason about goroutine leaks.

Primary docs:

- [`07-worker-pool-v2/README.md`](07-worker-pool-v2/README.md)
- [`docs/cheatsheet-concurrency.md`](docs/cheatsheet-concurrency.md)

Stop line:

- Keep v2 focused on lifecycle control. Do not turn it into a full service architecture with
  retries, metrics, tracing, persistent queues, or signal handling.

## Phase 5: Race Detector, Memory Model, And Runtime

Focus:

- explain what the race detector finds;
- explain what the race detector does not prove;
- understand happens-before at a practical level;
- connect synchronization choices to correctness;
- understand goroutine lifecycle;
- distinguish concurrency from parallelism;
- explain the Go scheduler model at a high level;
- reason about CPU-bound versus I/O-bound concurrency.

Primary docs:

- [`docs/cheatsheet-concurrency.md`](docs/cheatsheet-concurrency.md)
- [`go-release-history.md`](go-release-history.md)

Planned output:

- `05-race-detector-memory-model-runtime/README.md`
- targeted examples with `go test -race` where useful.

## Phase 6: Types, Interfaces, And Pipeline v1

Core focus:

- define structs with clear ownership;
- choose pointer or value receivers intentionally;
- keep interfaces small and consumer-owned;
- explain method sets and interface satisfaction.

Concurrency focus:

- build a simple pipeline;
- connect stages with channels;
- close outbound channels correctly;
- understand stage ownership and backpressure.

Planned output:

- `06-types-interfaces-pipeline-v1/README.md`
- small examples for pipeline stages and channel ownership.

## Phase 7: Structure And Pipeline v2

Core focus:

- understand `go.mod` and module paths;
- choose package names deliberately;
- separate exported API from private helpers;
- understand `internal` package boundaries;
- explain `cmd/` packages versus library packages.

Concurrency focus:

- add cancellation to pipelines;
- handle errors across stages;
- introduce fan-out and fan-in;
- avoid goroutine leaks when a downstream stage stops early.

Primary docs:

- [`01-project-structure/README.md`](01-project-structure/README.md)
- [`docs/cheatsheet-concurrency.md`](docs/cheatsheet-concurrency.md)

Planned output:

- `07-structure-pipeline-v2/README.md`
- focused examples for cancellation, errors, fan-out, and fan-in.

## Phase 8: Generics, Tooling, Workspaces, And Shared State

Core focus:

- use generics where they remove duplication without hiding intent;
- understand constraints and type parameters;
- use Go tooling deliberately;
- understand workspace and module boundaries.

Concurrency focus:

- compare mutexes, channels, and atomics;
- identify when shared memory is simpler than channel coordination;
- understand data races and synchronization boundaries.

Planned output:

- `08-generics-tooling-workspaces-shared-state/README.md`
- examples for mutex, channel ownership, and atomic counters.

## Phase 9: Go Core

Planned with the lowest priority. Define the detailed scope before starting this phase.

## Book Track: Functional Programming In Go

Step 13 is separate from the phase priorities and may be studied first.

Read *Learning Functional Programming in Go* by Lex Sheehan (2017) as a dedicated book track. The book predates Go generics, so discuss its `interface{}`-based collection helpers in terms of modern type parameters.

Short reading plan:

1. **Pure functional programming in Go**: imperative versus declarative style, pure functions, recursion, memoization, closures, tests, and benchmarks.
2. **Manipulating collections**: iteration, composition, `map`, `filter`, `reduce`, predicates, and the book's pre-generics collection abstractions.
3. **Higher-order functions**: first-class functions, function composition, currying, generators, and the examples that use goroutines and `WaitGroup`.
4. **SOLID design in Go**: connect functional composition with interfaces, embedding, error handling, and MapReduce.
5. **Decoration and dependency injection**: interface composition, `io.Reader`/`io.Writer`, decorator, strategy, inversion of control, and lifecycle coordination with channels.
6. **Functional ideas at architecture level**: state management, dependency direction, boundaries, layers, observer, and dependency injection.

After each practical block, write a short comparison: what remains idiomatic Go, what modern generics simplify, and where a direct loop or ordinary interface is clearer than a functional abstraction. Do not treat recursion, reflection, monads, or generic collection helpers as defaults; the interview goal is to explain the trade-off.

## Practice Track

The `coding/` directory is separate from the main lab sequence. Use it for interview-style practice
and algorithm exercises without mixing those problems into the core/concurrency roadmap.

Current example:

- [`coding/hackerrank/coding3/README.md`](coding/hackerrank/coding3/README.md)

## Working Rules

- Prefer small, reviewable steps.
- Read notes before writing code.
- Keep each lab focused on one concept boundary.
- Avoid broad rewrites.
- Prefer targeted commands for the current lab.
- Do not run `go test ./...` unless explicitly requested.

Useful targeted checks:

```bash
go test ./04-slices-maps-defer/slices-lab
go test ./06-worker-pool-v1/workerpool
go test -race ./06-worker-pool-v1/workerpool
```
