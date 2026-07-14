# Project Context

This repository is a Go Core and Concurrency learning project.

## Current Focus

Phase 1, Slices, is finished and needs review.

Phase 2, Worker Pool v1, is finished and closed.

The current work is Phase 3, Maps.

Use `PLAN.md` as the source of truth for the roadmap and phase priorities.

## Learning Sequence

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

The current priorities group that path into phases:

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

The working style is:

- read and poke with small examples;
- write focused lab code;
- keep notes in cheatsheets;
- prepare short interview answers;
- avoid turning one lab into a broad rewrite.

## Worker Pool Boundary

Worker Pool v1 should stay minimal:

- fixed number of workers;
- `jobs` channel;
- `results` channel;
- `sync.WaitGroup`;
- clear channel ownership;
- correct `results` closing after all workers finish.

Worker Pool v1 should not include:

- context cancellation;
- timeouts;
- graceful shutdown;
- goroutine leak checks;
- race-detector hardening;
- advanced error policy.

Those belong to Worker Pool v2.

For this project, broad test runs should be opt-in. Prefer review, small examples, and targeted
manual reasoning unless explicit test execution is requested.
