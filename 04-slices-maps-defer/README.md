# Slices, Maps, Defer

This lab covers core Go collection and cleanup topics from the roadmap.

- Phase 1: Slices - finished, needs review.
- Phase 3: Maps - current.
- Phase 4: Defer - planned.

## Links

- [Root README](../README.md)
- [Roadmap](../PLAN.md)
- [Core cheatsheet](../docs/cheatsheet-core.md)

## Current Status

Slices are implemented as focused examples and tests under `slices-lab`.

Maps are the current topic; defer is still planned. Keep future changes small and separate: do not
mix map examples, defer examples, and slice review fixes in one broad rewrite.

## Slices

A slice is a small descriptor over part of an array. It contains a pointer to a backing array, a length, and a capacity.

Key points to review:

- a slice does not store elements directly;
- `len` is the current number of visible elements;
- `cap` is the available capacity from the slice start to the end of the backing array;
- `append` may reuse the existing backing array;
- `append` may allocate a new backing array when capacity is not enough;
- subslices can accidentally share the same backing array;
- subslices can keep a large backing array alive;
- `nil` slices and empty slices behave similarly, but are not identical;
- use `copy` or a full slice expression to avoid unwanted aliasing;
- slices are not safe for concurrent mutation without synchronization.

Review checklist:

- explain `len` vs `cap`;
- explain `append` with and without spare capacity;
- show a subslice aliasing bug;
- fix aliasing with `copy` or a full slice expression;
- explain nil vs empty slice behavior;
- explain the large backing-array retention problem.

Recommended reading:

- [Go Slices: usage and internals](https://go.dev/blog/slices-intro)
- [Core cheatsheet: Slices](../docs/cheatsheet-core.md#slices)

## Maps

Current focus:

- create maps with literals and `make`;
- read, insert, update, and delete keys;
- use the comma-ok form for lookups;
- distinguish nil maps from empty maps;
- understand that map iteration order is not stable;
- use maps for counting, grouping, indexing, and set-like behavior;
- avoid concurrent map reads/writes without synchronization.

## Defer

Planned focus:

- use `defer` for cleanup;
- understand LIFO execution order;
- understand when deferred function arguments are evaluated;
- avoid hiding important control flow;
- use defer with files, locks, recovery examples, and small cleanup cases.

## Lab Files

```text
04-slices-maps-defer/
  README.md
  main.go
  slices-lab/
    append_copy.go
    helpers.go
    memory_leak.go
    mistakes.go
    mistakes_test.go
    slices.go
```

`slices-lab` contains focused slice examples and tests. Keep examples small enough to explain line by line.

## Targeted Checks

Run only the focused package tests for this lab:

```bash
go test ./04-slices-maps-defer/slices-lab
```

Avoid broad test runs such as `go test ./...` unless explicitly requested.

## Review Questions

Slices:

1. What is stored in a slice header?
2. When does `append` reuse the existing backing array?
3. When does `append` allocate a new backing array?
4. How can two slices accidentally modify the same array?
5. How does a full slice expression help prevent aliasing?
6. What is the difference between a nil slice and an empty slice?
7. How can a small subslice keep a large array in memory?

Maps:

1. What happens when reading a missing key?
2. Why is the comma-ok form useful?
3. What is the difference between a nil map and an empty map?
4. Why should code not depend on map iteration order?
5. Why are normal maps unsafe for concurrent read/write access?

Defer:

1. When are deferred function arguments evaluated?
2. In what order do multiple deferred calls run?
3. What cleanup tasks are good candidates for `defer`?
4. When can `defer` make control flow harder to read?
