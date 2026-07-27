# Go Core Cheatsheet

Translations: [Russian](cheatsheet-core.ru.md) · [Ukrainian](cheatsheet-core.ua.md).

## Slices

### Definition

A slice is a small descriptor for a contiguous part of an underlying array. Conceptually, it
contains:

- a pointer to the first visible element of the backing array;
- a length, returned by `len`;
- a capacity, returned by `cap`.

A slice does not store its elements itself. Copying a slice value copies only this descriptor, so
two
slices can refer to the same backing array.

```go
array := [4]int{10, 20, 30, 40}
slice := array[1:3] // [20 30]

slice[0] = 99
fmt.Println(array) // [10 99 30 40]
```

### `len` and `cap`

`len(s)` is the number of elements currently visible through `s`.

`cap(s)` is the number of elements available from the start of `s` to the end of its backing array.
It is not necessarily the capacity of the original slice.

```go
s := []int{2, 3, 5, 7, 11, 13}
fmt.Println(len(s), cap(s)) // 6 6

s = s[:4]
fmt.Println(len(s), cap(s)) // 4 6

s = s[2:]
fmt.Println(len(s), cap(s)) // 2 4
```

Reslicing can change the length without allocating. It can extend a slice only up to its capacity.

### `append`

`append` returns the resulting slice, so always keep its return value.

```go
s := []int{1, 2}
s = append(s, 3)
```

If the slice has spare capacity, `append` may reuse its backing array. Other slices sharing that
array
can then observe the written elements.

```go
a := make([]int, 2, 4)
a[0], a[1] = 10, 20
b := append(a, 30)

b[0] = 99
fmt.Println(a) // [99 20]
fmt.Println(b) // [99 20 30]
```

If capacity is insufficient, `append` allocates a new backing array and copies the old elements.
Do not rely on a particular capacity-growth formula; that is an implementation detail.

```go
a := make([]int, 2, 2)
a[0], a[1] = 10, 20
b := append(a, 30)

b[0] = 99
fmt.Println(a) // [10 20]
fmt.Println(b) // [99 20 30]
```

### Subslice aliasing

A subslice normally shares its backing array with the original slice. Appending to a subslice can
overwrite elements that are visible through another slice.

```go
a := []int{1, 2, 3, 4}
b := a[:2] // len=2, cap=4
b = append(b, 99)

fmt.Println(a) // [1 2 99 4]
```

Use a full slice expression to limit the capacity exposed to a subslice. It forces a later `append`
to allocate, but it does not copy the existing elements.

```go
a := []int{1, 2, 3, 4}
b := a[:2:2] // len=2, cap=2
b = append(b, 99)

fmt.Println(a) // [1 2 3 4]
fmt.Println(b) // [1 2 99]
```

Use a full slice expression when the slice should still be a view, but callers must not grow it into
the remaining backing array. Use a copy when complete independence is required.

### Copying and cloning

`copy(dst, src)` copies up to `min(len(dst), len(src))` elements and returns the number copied. It
does not change the length of `dst`.

```go
src := []int{1, 2, 3, 4}
dst := make([]int, 2)
n := copy(dst, src)

fmt.Println(n, dst) // 2 [1 2]
```

To make an independent clone while preserving the distinction between `nil` and empty slices:

```go
func Clone[T any](s []T) []T {
if s == nil {
return nil
}

out := make([]T, len(s))
copy(out, s)
return out
}
```

`append([]T(nil), s...)` is also a common clone idiom. For an empty non-nil slice, however, it may
return a nil slice, so choose it only when that distinction does not matter.

### `nil` versus empty slices

Both nil and empty slices have length zero, can be ranged over, and can be passed to `append`.
They are different values.

```go
var nilSlice []int
emptySlice := []int{}

fmt.Println(len(nilSlice), nilSlice == nil) // 0 true
fmt.Println(len(emptySlice), emptySlice == nil) // 0 false
```

The difference matters at API boundaries. For example, `encoding/json` encodes a nil slice as `null`
and an empty slice as `[]`.

For an emptiness check, use `len(s) == 0` unless your API specifically distinguishes nil from empty.

### Retaining a large backing array

A small subslice keeps the entire backing array reachable. Returning a 10-byte subslice of a 100 MB
buffer can therefore retain roughly 100 MB of memory.

```go
func firstTenBad() []byte {
big := make([]byte, 100<<20)
return big[:10]
}
```

Copy the small part before returning it when the large buffer is no longer needed.

```go
func firstTenGood() []byte {
big := make([]byte, 100<<20)
small := make([]byte, 10)
copy(small, big[:10])
return small
}
```

### Common operations and ownership

In-place operations reuse the backing array and may modify the input slice's elements. A function
that
returns a new slice allocates separate storage and leaves the input unchanged.

```go
// In place: the input's backing array is modified.
func DeleteAt[T any](s []T, i int) []T {
copy(s[i:], s[i+1:])
var zero T
s[len(s)-1] = zero // release a reference held in the unused tail, if any
return s[:len(s)-1]
}

// New slice: the input is unchanged.
func DeleteAtNew[T any](s []T, i int) []T {
out := make([]T, 0, len(s)-1)
out = append(out, s[:i]...)
return append(out, s[i+1:]...)
}
```

When documenting or designing an API, state clearly whether a function may retain, modify, or reuse
the caller's slice.

### `range` over slices

The value produced by `range` is a copy of the element. Change elements through their index when
needed.

```go
s := []int{1, 2, 3}
for _, v := range s {
v *= 10 // changes only the copy
}

for i := range s {
s[i] *= 10 // changes the slice element
}
```

The number of iterations is determined when the `range` starts. Appending during that loop does not
make it iterate over the appended elements; avoid this pattern unless the exact behavior is
intended.

### Concurrency

Slices are not safe for concurrent mutation. If goroutines share a slice and at least one writes to
it, coordinate access with synchronization or give each goroutine an independent copy.

### Review questions

1. What does a slice descriptor contain?
2. Why can changing one slice change another one?
3. When can `append` reuse the backing array, and when must it allocate?
4. What does `s[low:high:max]` change, and what does it not change?
5. Why does `copy` need a destination with a non-zero length?
6. When does nil differ from an empty slice in practice?
7. How can a small slice retain a large allocation?

### Answers to review questions

1. A slice descriptor conceptually contains a pointer to its backing array, a length, and a
   capacity.
2. Copying or subslicing a slice normally leaves both slices pointing to the same backing array.
   Writing an element through one slice therefore changes that array and can be visible through the
   other.
3. `append` may reuse the backing array when the resulting length fits within the current capacity.
   It must allocate a new array when the capacity is insufficient.
4. `s[low:high:max]` sets the resulting slice's length to `high-low` and its capacity to `max-low`.
   It limits future growth through that slice, but does not copy elements or remove sharing of the
   existing visible elements.
5. `copy` copies only into existing destination elements, up to the smaller source or destination
   length. A nil or zero-length destination has no elements to receive the copy.
6. They differ when an API observes nilness, notably JSON encoding (`null` versus `[]`), or when nil
   itself has semantic meaning. Otherwise, use `len(s) == 0` to test emptiness.
7. A slice descriptor keeps its backing array reachable. A small subslice still points into the
   original large array, so the garbage collector cannot reclaim that array until the subslice is no
   longer reachable. Copy the needed elements into a new slice to release it.

### Related lab

See [`04-slices-maps-defer/slices-lab`](../04-slices-maps-defer/slices-lab) for runnable examples.
Its seven focused unit tests cover clone semantics, input ownership during append and deletion, and
zeroing the unused tail after in-place operations.
Run the focused package check with:

```bash
go test ./04-slices-maps-defer/slices-lab
```

## Maps

### Basic operations

A map associates comparable keys with values: `map[K]V`. Create one with a literal when its initial
entries are known, or with `make` when it starts empty.

```go
prices := map[string]int{"apple": 10}
stock := make(map[string]int)
```

Read a value with `m[key]`. When the key is absent, the expression returns the zero value of the
map's value type. For `int`, that is `0`.

```go
fmt.Println(prices["apple"]) // 10
fmt.Println(prices["pear"])  // 0
```

Use the same assignment syntax for insertion and update. An assignment to a new key inserts it; an
assignment to an existing key replaces its value.

```go
stock["apple"] = 3 // insert
stock["apple"] = 5 // update
```

`delete(m, key)` removes a key. It is safe to delete a key that is already absent. `len(m)` returns
the number of entries currently stored in the map.

```go
delete(stock, "apple")
fmt.Println(len(stock))
```

The zero value returned by a missing lookup cannot distinguish an absent key from a key explicitly
stored with the zero value. The comma-ok lookup that solves this belongs to the next map group.

### Review questions

1. How can you create a map with initial entries and an empty writable map?
2. What happens when `m[key]` reads an absent key?
3. Which operation inserts a new key and updates an existing key?
4. What happens when `delete` receives an absent key?
5. What does `len(m)` return?

### Answers to review questions

1. Use a literal such as `map[string]int{"apple": 10}` for initial entries, and `make(map[string]int)` for an empty writable map.
2. It returns the zero value of the map's value type.
3. `m[key] = value` inserts if the key is absent and replaces the value if it exists.
4. Nothing; deleting an absent key is safe.
5. The current number of key-value entries.

### Reliable lookups and map state

An ordinary lookup cannot tell a missing key apart from a key that stores the zero value. Use the
comma-ok form when that distinction matters.

```go
stock := map[string]int{"apple": 0}

quantity, present := stock["apple"] // 0, true
missing, present := stock["pear"]   // 0, false
```

A nil map and an empty initialized map both have length zero and can be read safely. A lookup from
either returns the value type's zero value; comma-ok returns `false` for a missing key.

```go
var nilStock map[string]int
emptyStock := make(map[string]int)

fmt.Println(len(nilStock), nilStock == nil)     // 0 true
fmt.Println(len(emptyStock), emptyStock == nil) // 0 false
```

The important difference is writing: assigning to a nil map panics, while an empty map made with
`make` is ready for inserts and updates. Prefer `len(m) == 0` when you only need to know whether a
map has entries; test `m == nil` only when nilness is meaningful to the API.

### Review questions: reliable lookups and state

1. Why is `value := m[key]` ambiguous for maps with zero-valued entries?
2. What do the two results of `value, ok := m[key]` mean?
3. Which basic operations are safe on a nil map?
4. What happens when assigning to a nil map?
5. When should code use `len(m) == 0` instead of `m == nil`?

### Answers to review questions: reliable lookups and state

1. A missing key and a present key whose value is the type's zero value both produce that zero value.
2. `value` is the stored value or the zero value; `ok` reports whether the key is present.
3. Reading, comma-ok lookup, `len`, and `delete` are safe. Writing is not.
4. It panics at runtime.
5. Use `len(m) == 0` for an emptiness check; use `m == nil` only if the API gives nil a separate meaning.

### Iteration order

The order produced by `range` over a map is not specified and must not be used as program logic or
test expectation. Collect keys and sort them when output needs a stable order.

```go
keys := make([]string, 0, len(counts))
for key := range counts {
	keys = append(keys, key)
}
sort.Strings(keys)
```

Tests for an unordered result should compare membership and length, or compare a normalized sorted
result. Do not make a test pass by assuming one observed map order is permanent.

### Common patterns

Counting uses the zero value of an integer, so incrementing a new key works without a preliminary
lookup.

```go
counts := make(map[string]int)
for _, word := range words {
	counts[word]++
}
```

Grouping uses a slice as the map value. `append` works on a nil slice, so the first value in a group
needs no special case. This example groups by byte length.

```go
groups := make(map[int][]string)
for _, word := range words {
	groups[len(word)] = append(groups[len(word)], word)
}
```

Indexing converts a slice into a lookup table. Define how duplicate keys behave; ordinary assignment
makes the later entry replace the earlier one.

```go
index := make(map[int]User)
for _, user := range users {
	index[user.ID] = user
}
```

A set is commonly represented by `map[T]struct{}`. The empty struct signals membership without a
separate value payload.

```go
set := make(map[string]struct{})
set["go"] = struct{}{}
_, hasGo := set["go"]
```

### Review questions: iteration and patterns

1. Why is map iteration order unsuitable for program logic?
2. How can map keys be displayed in a stable order?
3. Why can `counts[key]++` count a new key directly?
4. Why does `groups[key] = append(groups[key], value)` work for a new group?
5. What happens to duplicate keys in an index built with ordinary assignment?
6. Why is `map[T]struct{}` a suitable set representation?

### Answers to review questions: iteration and patterns

1. The language does not specify the order, so it can vary between iterations and runs.
2. Collect the keys into a slice and sort the slice before displaying it.
3. A missing `int` map entry reads as zero, and incrementing zero produces the first count.
4. A missing slice entry is nil, and `append` accepts a nil slice.
5. The later assignment replaces the earlier value for that key.
6. Keys represent members and `struct{}` is an empty value, so the map stores membership without a payload.

### Concurrency boundary

An ordinary map is unsafe when goroutines access it concurrently and at least one goroutine writes.
That is a data race. It may produce a runtime error such as concurrent map read and map write, but
the absence of that error does not prove the code is safe.

For a small shared map with simple operations, a mutex-protected wrapper is usually the clearest
first choice. Use `sync.Mutex` when every operation writes; use `sync.RWMutex` when independent
reads are common and profiling shows that the extra complexity is worthwhile.

```go
type SafeInventory struct {
	mu sync.RWMutex
	m  map[string]int
}

func (s *SafeInventory) Lookup(key string) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.m[key]
	return value, ok
}
```

An alternative is channel ownership: one goroutine owns the map and other goroutines send it
requests. This is useful when map operations are part of a larger coordination protocol. `sync.Map`
is a specialised concurrent type for particular access patterns; do not choose it by default instead
of an ordinary map plus a mutex.

Use the race detector as a separate check:

```bash
go test -race ./04-slices-maps-defer/maps-lab
```

### Internals for interviews

At the language level, `map` is a built-in associative data type. Its runtime representation is not
part of the Go specification, so application code must not depend on buckets, addresses, growth
thresholds, or iteration layout. Lookup, insertion, update, and deletion have expected O(1) cost;
collisions and resizing mean this is not a per-operation worst-case guarantee.

Before Go 1.24, the usual description of the runtime implementation was buckets with overflow
buckets. Go 1.24 replaced it with an implementation based on Swiss Tables. It uses open addressing,
logical groups of eight slots, and control metadata containing slot state and a fragment of each
key's hash. This lets lookups discard most nonmatching slots efficiently.

Go retains incremental-growth behaviour for latency-sensitive programs: a map is split into smaller
tables so one insertion does not need to copy an arbitrarily large map. These details are useful for
an interview answer, but the observable rules are unchanged: iteration order is unspecified and
unsynchronised concurrent access with a writer is unsafe.

Read the official [Go 1.24 release notes](https://go.dev/doc/go1.24) and the Go team's article
[Faster Go maps with Swiss Tables](https://go.dev/blog/swisstable) for the implementation details.

### Review questions: concurrency and internals

1. When is a normal map unsafe to share between goroutines?
2. Why is a runtime concurrent-map failure not a substitute for `go test -race`?
3. When is a mutex-protected map a clear first choice?
4. What is the channel-ownership alternative?
5. Why should `sync.Map` not be the default concurrent map?
6. Which map implementation details are guaranteed by the language specification?
7. What high-level implementation change arrived in Go 1.24?

### Answers to review questions: concurrency and internals

1. When goroutines access it concurrently and at least one writes, unless access is synchronised.
2. It may fail to occur; the race detector instruments execution to report races exercised by tests.
3. For a small shared map with straightforward reads and writes.
4. One goroutine owns the map and other goroutines request operations through channels.
5. It is designed for particular concurrent access patterns; an ordinary map plus a mutex is often simpler.
6. None of its runtime layout details; only the language-level map semantics are guaranteed.
7. The runtime moved from the earlier bucket-based design to a Swiss Table-based implementation.

### Related lab

See [`04-slices-maps-defer/maps-lab`](../04-slices-maps-defer/maps-lab) for runnable basic map
examples, reliable lookup examples, iteration, common patterns, safe concurrency, and focused unit tests.

```bash
go test ./04-slices-maps-defer/maps-lab
```
