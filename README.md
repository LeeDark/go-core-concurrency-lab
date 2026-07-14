# Go Core and Concurrency Lab

This repository is a learning lab for Go core language topics and concurrency patterns.

It is not a production library. The goal is to study Go concepts in small, reviewable steps, write
focused examples, keep notes, and build toward practical concurrency patterns such as worker pools
and pipelines.

Module:

```text
github.com/LeeDark/go-core-concurrency-lab
```

Go version:

```text
1.25.6
```

## Learning Sequence

The main study path is:

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

The private study plan pairs core topics with concurrency labs:

| Core topic                            | Concurrency topic      |
|---------------------------------------|------------------------|
| Step 1: project structure             | Step 7: Worker Pool v1 |
| Step 2: Go Core                       | Step 8: Worker Pool v2 |
| Step 3: errors, context               | Step 9: Pipeline v1    |
| Step 4: structs, methods, interfaces  | Step 10: Pipeline v2   |
| Step 5: slices, maps, defer           | Step 11: shared state  |
| Step 6: generics, tooling, workspaces | Step 12: race detector |

The current priorities group that path into phases. Phase 1, Slices, is finished and needs review.
Phase 2, Worker Pool v1, is finished and closed. The current work is Phase 3, Maps. See
[`PLAN.md`](PLAN.md) for the full roadmap.

| Phase | Topic                                           | Status                   |
|-------|-------------------------------------------------|--------------------------|
| 1     | Slices                                          | Finished, needs review   |
| 2     | Worker Pool v1                                  | Finished, closed         |
| 3     | Maps                                            | Current                  |
| 4     | Defer, errors, context, and Worker Pool v2      | Planned                  |
| 5     | Race detector, memory model, and runtime        | Planned                  |
| 6     | Types, interfaces, and Pipeline v1              | Planned                  |
| 7     | Structure and Pipeline v2                       | Planned                  |
| 8     | Generics, tooling, workspaces, and shared state | Planned                  |
| 9     | Go Core                                         | Planned, lowest priority |

## Repository Layout

```text
01-project-structure/       Notes for modules, packages, visibility, internal packages.
04-slices-maps-defer/       Slice/map/defer notes and focused slice examples.
06-worker-pool-v1/          Minimal worker pool with channels and WaitGroup.
07-worker-pool-v2/          Planned lifecycle-focused worker pool notes.
coding/                     Coding-practice exercises.
docs/                       Core and concurrency cheatsheets.
docs/ai/project-context.md  AI-assistant project context and learning boundaries.
go-release-history.md       Go release notes summary.
PLAN.md                     Roadmap for core topics and concurrency labs.
```

## Current Labs

### Project Structure

[`01-project-structure`](01-project-structure/README.md) covers module layout, package naming,
exported and unexported identifiers, `internal` packages, command packages, and documentation
comments.

### Slices, Maps, Defer

[`04-slices-maps-defer`](04-slices-maps-defer/README.md) focuses on slice internals: backing arrays,
`len`, `cap`, `append`, `copy`, nil vs empty slices, aliasing, and common mistakes.

The `slices-lab` package contains small examples and targeted tests for slice behavior.

### Worker Pool v1

Finished and closed. The minimal worker-pool implementation remains as a reference for later
concurrency labs.

[`06-worker-pool-v1`](06-worker-pool-v1/README.md) builds the first minimal worker pool:

```text
producer -> jobs channel -> N workers -> results channel -> consumer
```

Worker Pool v1 intentionally stays small:

- fixed worker count;
- `jobs` channel;
- `results` channel;
- `sync.WaitGroup`;
- clear channel ownership;
- `results` closes only after all workers finish.

It does not include context cancellation, timeouts, graceful shutdown, leak checks, or advanced
error policy. Those belong to Worker Pool v2.

### Worker Pool v2

[`07-worker-pool-v2`](07-worker-pool-v2/README.md) extends the v1 mental model with lifecycle
control:

- `context.Context`;
- cancellation while waiting for jobs;
- cancellation while sending results;
- timeout policy;
- error policy;
- goroutine leak reasoning;
- graceful stop semantics.

## Notes And Cheatsheets

- [`docs/cheetsheet-core.md`](docs/cheetsheet-core.md) contains core Go notes.
- [`docs/cheetsheet-concurrency.md`](docs/cheetsheet-concurrency.md) contains concurrency,
  goroutine, channel, worker-pool, scheduler, and interview notes.
- [`go-release-history.md`](go-release-history.md) tracks selected Go release changes.

## Running Focused Checks

Prefer targeted commands for the lab you are working on.

Examples:

```bash
go test ./04-slices-maps-defer/slices-lab
go test ./06-worker-pool-v1/workerpool
go test -race ./06-worker-pool-v1/workerpool
```

Avoid broad test runs such as:

```bash
go test ./...
```

Run broad tests only when explicitly needed.

## Working Style

This repo is optimized for learning:

- read notes and small examples first;
- make focused lab changes;
- keep changes reviewable;
- avoid turning one lab into a broad rewrite;
- keep Worker Pool v1 minimal before moving lifecycle concerns into Worker Pool v2;
- use targeted tests or manual reasoning depending on the lab.

The main question for each lab is not only "does it run?", but also "can the behavior be explained
clearly in an interview or code review?".
