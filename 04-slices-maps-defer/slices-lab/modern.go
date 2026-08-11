package sliceslab

import "slices"

// ClearSlice sets every element of s to its zero value and keeps its length and capacity.
func ClearSlice[T any](s []T) []T {
	clear(s)
	return s
}

// CloneSliceStdlib makes an independent shallow clone using the standard library.
func CloneSliceStdlib[S ~[]E, E any](s S) S {
	return slices.Clone(s)
}

// DeleteSliceStdlib removes s[i:j] using the standard library and clears the unused tail.
func DeleteSliceStdlib[S ~[]E, E any](s S, i, j int) S {
	return slices.Delete(s, i, j)
}
