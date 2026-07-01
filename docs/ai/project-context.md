# Project Context

This repository is a Go Core and Concurrency learning project.

## Current Focus

Phase 1, Slices, is finished and needs review.

Current work is Phase 2, Worker Pool v1.

Use `PLAN.md` as the source of truth for the roadmap and phase priorities.

## Learning Sequence

Current high-level sequence:

1. Slices.
2. Worker Pool v1.
3. Maps.
4. Defer, errors, context, and Worker Pool v2.
5. Race detector, memory model, and runtime.
6. Types, interfaces, and Pipeline v1.
7. Structure and Pipeline v2.
8. Generics, tooling, workspaces, and shared state.

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

For this project, broad test runs should be opt-in. Prefer review, small examples, and targeted manual reasoning unless explicit test execution is requested.
