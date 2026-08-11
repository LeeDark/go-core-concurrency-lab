# Slices and Maps

This lab covers core Go collection topics from the roadmap.

- Phase 1: Slices - finished, closed.
- Phase 3: Maps - finished, closed.

## Links

- [Root README](../README.md)
- [Roadmap](../PLAN.md)
- [Core cheatsheet](../docs/cheatsheet-core.md)

## Current Status

Slices are implemented as focused examples and tests under `slices-lab`. Phase 1 is closed.

Maps are implemented as focused examples and tests under `maps-lab`. Phase 3 is closed. Defer is
covered by Phase 4 and has a separate lab planned there.

The default `main.go` run skips the memory-heavy retention demonstrations in `slices-lab`. Call
`MemoryLeakSubslice` and `Mistake26` explicitly when studying those examples.

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
    modern.go
    modern_test.go
    patterns.go
    patterns_test.go
    state.go
    state_test.go
  slices-lab/
    append_copy.go
    append_copy_test.go
    helpers.go
    memory_leak.go
    modern.go
    modern_test.go
    mistakes.go
    mistakes_test.go
    slices.go
```

`slices-lab` contains focused slice examples and tests. Keep examples small enough to explain line by line.

`maps-lab` contains focused map examples. `MapBasics` demonstrates each basic operation;
`MapState` demonstrates comma-ok lookups and the difference between nil and empty maps. The tests
verify their observable behavior without relying on iteration order. `MapPatterns` demonstrates how
to make output deterministic and how to use maps for counting, grouping, indexing, and sets.
`MapConcurrency` demonstrates a map protected by `sync.RWMutex`. `modern.go` demonstrates the built-in
`clear` and the standard `slices` and `maps` packages.

### Focused Unit Tests

`append_copy_test.go` contains seven focused tests for slice ownership and mutation behavior:

- preserving `nil` and creating an independent clone;
- the empty-to-`nil` behavior of `CloneSliceAppend`;
- returning an appended slice without modifying the input;
- in-place versus independent deletion;
- clearing the unused tail after in-place deletion and filtering.

`modern_test.go` covers `clear`, standard-library cloning, and standard-library deletion.

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
8. What does `clear` change in a slice?
9. What does `slices.Delete` do with the unused tail?
10. What remains shared when a slice is passed to a function?

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
17. Which types can be map keys, and which cannot?
18. What does assigning one map variable to another copy?
19. How do you update a field in a struct stored as a map value?
20. What does the capacity argument to `make` mean for a map?
21. What does `clear` do to a map?
22. Which common operations are provided by the standard `maps` package?
