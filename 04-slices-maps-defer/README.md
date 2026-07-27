# Slices, Maps, Defer

This lab covers core Go collection and cleanup topics from the roadmap.

- Phase 1: Slices - finished, closed.
- Phase 3: Maps - current.
- Phase 4: Defer - planned.

## Links

- [Root README](../README.md)
- [Roadmap](../PLAN.md)
- [Core cheatsheet](../docs/cheatsheet-core.md)

## Current Status

Slices are implemented as focused examples and tests under `slices-lab`. Phase 1 is closed.

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

Phase 3, groups 1–4: basic operations, reliable lookups, map state, iteration, patterns, and safe
concurrent access.

Learn to:

- create maps with literals and `make`;
- read values by key;
- insert a new key or update an existing key with `m[key] = value`;
- delete a key with `delete`;
- observe `len` and the zero value returned for a missing key.
- use `value, ok := m[key]` to distinguish a missing key from a stored zero value;
- distinguish a nil map from an empty map created with `make`.
- avoid relying on the order produced by `range` over a map;
- use maps for counting, grouping, indexing, and set-like membership checks.
- protect shared maps with synchronization;
- explain the public map contract separately from runtime internals.

Phase 3 is complete. Mutexes, channels, atomics, and `sync.Map` will be compared in depth in Phase 8.

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
  maps-lab/
    basic.go
    basic_test.go
    concurrency.go
    concurrency_test.go
    maps.go
    patterns.go
    patterns_test.go
    state.go
    state_test.go
  slices-lab/
    append_copy.go
    append_copy_test.go
    helpers.go
    memory_leak.go
    mistakes.go
    mistakes_test.go
    slices.go
```

`slices-lab` contains focused slice examples and tests. Keep examples small enough to explain line by line.

`maps-lab` contains the first focused map examples. `MapBasics` demonstrates each basic operation;
`MapState` demonstrates comma-ok lookups and the difference between nil and empty maps. The tests
verify their observable behavior without relying on iteration order. `MapPatterns` demonstrates how
to make output deterministic and how to use maps for counting, grouping, indexing, and sets.
`MapConcurrency` demonstrates a map protected by `sync.RWMutex`.

### Focused Unit Tests

`append_copy_test.go` contains seven focused tests for slice ownership and mutation behavior:

- preserving `nil` and creating an independent clone;
- the empty-to-`nil` behavior of `CloneSliceAppend`;
- returning an appended slice without modifying the input;
- in-place versus independent deletion;
- clearing the unused tail after in-place deletion and filtering.

## Targeted Checks

Run only the focused package tests for this lab:

```bash
go test ./04-slices-maps-defer/slices-lab
go test ./04-slices-maps-defer/maps-lab
go test -race ./04-slices-maps-defer/maps-lab
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
2. Which one assignment syntax both inserts a key and updates an existing key?
3. What does `delete` do when the key is already absent?
4. What does `len` measure for a map?
5. Why can ordinary lookup not distinguish an absent key from a key with value `0`?
6. How does `value, ok := m[key]` solve this ambiguity?
7. Which operations are safe on a nil map, and which operation panics?
8. How does an empty map created by `make` differ from a nil map?
9. Why must code not depend on the order of a map `range`?
10. How can a program produce deterministic output from map keys?
11. Which map value types fit counting, grouping, indexing, and a set?
12. Why are concurrent map reads and writes unsafe without synchronization?
13. When is a mutex-protected map a good first choice?
14. What does the race detector check, and why is a runtime concurrent-map failure not enough?
15. Which parts of map internals may application code rely on?
16. What changed in the Go 1.24 map implementation?

Defer:

1. When are deferred function arguments evaluated?
2. In what order do multiple deferred calls run?
3. What cleanup tasks are good candidates for `defer`?
4. When can `defer` make control flow harder to read?
