package sliceslab

import "fmt"

func printSlice[T any](s []T) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func printlnSlice[T any](s string, x []T) {
	fmt.Printf("%s len=%d cap=%d %v\n", s, len(x), cap(x), x)
}

func printSliceFull[T any](s []T) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
	fmt.Println("slice == nil:", s == nil)
}
