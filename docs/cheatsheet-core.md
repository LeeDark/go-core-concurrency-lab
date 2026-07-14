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
