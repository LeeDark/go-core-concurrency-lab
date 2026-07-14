package sliceslab

// CloneSliceCopy creates a copy of the input slice using the copy function, preserving the original slice's order and values.
func CloneSliceCopy[T any](s []T) []T {
	if s == nil {
		return nil
	}

	out := make([]T, len(s))
	copy(out, s)

	return out
}

// CloneSliceAppend creates a new slice by cloning the input slice and returns it. It avoids modifying the original slice.
func CloneSliceAppend[T any](s []T) []T {
	return append([]T(nil), s...)
}

// AppendSafe returns a new slice containing s followed by values without modifying s.
func AppendSafe[T any](s []T, values ...T) []T {
	out := make([]T, 0, len(s)+len(values))
	out = append(out, s...)
	out = append(out, values...)

	return out
}

// DeleteAt removes the element at the specified index i from the slice s and returns the resulting slice. Panics if i is out of range.
func DeleteAt[T any](s []T, i int) []T {
	if i < 0 || i >= len(s) {
		panic("index out of range")
	}

	copy(s[i:], s[i+1:])

	var zero T
	s[len(s)-1] = zero

	return s[:len(s)-1]
}

// DeleteAtNew removes the element at the specified index i from the slice and returns a new slice without modifying the original.
func DeleteAtNew[T any](s []T, i int) []T {
	if i < 0 || i >= len(s) {
		panic("index out of range")
	}

	out := make([]T, 0, len(s)-1)
	out = append(out, s[:i]...)
	out = append(out, s[i+1:]...)

	return out
}

// FilterInPlace filters a slice in place by keeping elements that satisfy the given predicate function and returns the result slice.
func FilterInPlace[T any](s []T, keep func(T) bool) []T {
	out := s[:0]

	for _, v := range s {
		if keep(v) {
			out = append(out, v)
		}
	}

	var zero T
	for i := len(out); i < len(s); i++ {
		s[i] = zero
	}

	return out
}

// FilterNewSlice filters elements in a slice based on the provided predicate function and returns a new filtered slice.
func FilterNewSlice[T any](s []T, keep func(T) bool) []T {
	out := make([]T, 0, len(s))

	for _, v := range s {
		if keep(v) {
			out = append(out, v)
		}
	}

	return out
}
