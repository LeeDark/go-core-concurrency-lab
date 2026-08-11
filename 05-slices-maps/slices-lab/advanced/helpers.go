package slicesadvanced

import "fmt"

func printlnSlice[T any](name string, values []T) {
	fmt.Printf("%s len=%d cap=%d %v\n", name, len(values), cap(values), values)
}
