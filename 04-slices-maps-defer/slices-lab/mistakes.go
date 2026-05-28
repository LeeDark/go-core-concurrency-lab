package sliceslab

import (
	"encoding/json"
	"fmt"
	"runtime"
)

// Mistake20 demonstrates slice creation, append operations, and the impact on length and capacity during subslicing and appending.
func Mistake20() {
	// slice, len, cap
	// Длина — это количество элементов, содержащихся в срезе, тогда как
	// емкость — это количество элементов в резервном массиве.

	s := make([]int, 3, 6)
	printlnSlice("s", s)
	s = append(s, 2)
	printlnSlice("s", s)
	s = append(s, 3, 4, 5)
	printlnSlice("s!", s)

	s1 := make([]int, 3, 6)
	printlnSlice("s1", s1)

	s2 := s1[1:3]
	printlnSlice("s2", s2)
	s2 = append(s2, 2)
	printlnSlice("s2", s2)
	s2 = append(s2, 3)
	printlnSlice("s2", s2)
	s2 = append(s2, 4)
	printlnSlice("s2", s2)
	s2 = append(s2, 5)
	printlnSlice("s2!", s2)
}

func convertEmptySlice(foos []int) []float32 {
	bars := make([]float32, 0)
	for _, foo := range foos {
		bars = append(bars, float32(foo))
	}
	return bars
}

func convertGivenCapacity(foos []int) []float32 {
	n := len(foos)
	bars := make([]float32, 0, n)
	for _, foo := range foos {
		bars = append(bars, float32(foo))
	}
	return bars
}

func convertGivenLength(foos []int) []float32 {
	n := len(foos)
	bars := make([]float32, n)
	for i, foo := range foos {
		bars[i] = float32(foo)
	}
	return bars
}

// Mistake21 demonstrates the creation and manipulation of slices using different strategies for initialization.
func Mistake21() {
	foos := []int{1, 2, 3, 4}
	printlnSlice("foos", foos)

	bars1 := convertEmptySlice(foos)
	printlnSlice("bars1", bars1)

	bars2 := convertGivenCapacity(foos)
	printlnSlice("bars2", bars2)

	bars3 := convertGivenLength(foos)
	printlnSlice("bars3", bars3)
}

func log(i int, s []string) {
	fmt.Printf("%d: empty=%t\tnil=%t\n", i, len(s) == 0, s == nil)
}

type customer struct {
	ID         string    `json:"id"`
	Operations []float32 `json:"operations"`
}

// Mistake22 demonstrates the behavior of nil vs empty slices in Go and their JSON serialization differences.
func Mistake22() {
	var s []string
	log(1, s)

	s = []string(nil)
	log(2, s)

	s = []string{}
	log(3, s)

	s = make([]string, 0)
	log(4, s)

	//
	var s1 []float32
	customer1 := customer{
		ID:         "foo",
		Operations: s1,
	}
	b, _ := json.Marshal(customer1)
	fmt.Println(string(b))

	s2 := make([]float32, 0)
	customer2 := customer{
		ID:         "bar",
		Operations: s2,
	}
	b, _ = json.Marshal(customer2)
	fmt.Println(string(b))
}

func handle(operations []float32) {
	printlnSlice("operations", operations)
}

func handleOperations(id string) {
	operations := getOperations(id)

	// incorrect
	//if operations != nil {
	//	handle(operations)
	//}

	// correct
	if len(operations) != 0 {
		handle(operations)
	}
}

func getOperations(id string) []float32 {
	operations := make([]float32, 0)
	if id == "" {
		return operations
	}
	return operations
}

// Mistake23 demonstrates proper handling of slices by using len to check emptiness instead of comparing to nil.
func Mistake23() {
	handleOperations("")
}

// Mistake24 demonstrates the correct and incorrect ways to copy slices in Go and the caveats of using uninitialized slices.
func Mistake24() {
	// incorrect
	src1 := []int{0, 1, 2}
	var dst1 []int
	copy(dst1, src1)
	printlnSlice("src1", src1)
	printlnSlice("dst1", dst1)

	// correct
	src2 := []int{0, 1, 2}
	dst2 := make([]int, len(src2))
	copy(dst2, src2)
	printlnSlice("src2", src2)
	printlnSlice("dst2", dst2)

	// correct with append
	src3 := []int{0, 1, 2}
	dst3 := append([]int(nil), src3...)
	printlnSlice("src3", src3)
	printlnSlice("dst3", dst3)
}

func f1(s []int) {
	_ = append(s, 10)
}

// Mistake25 demonstrates common pitfalls and behaviors when working with slices in Go, including aliasing and capacity sharing issues.
func Mistake25() {
	s1 := []int{1, 2, 3}
	printlnSlice("s1", s1)

	s2 := s1[1:2]
	printlnSlice("s2", s2)
	printlnSlice("s1", s1)

	s3 := append(s2, 10)
	printlnSlice("s3", s3)
	printlnSlice("s1", s1)
	printlnSlice("s2", s2)

	//
	s4 := []int{1, 2, 3}
	printlnSlice("s4", s4)
	f1(s4[:2])
	printlnSlice("s4", s4)

	// first option
	s5 := []int{1, 2, 3}
	s5Copy := make([]int, 2)
	copy(s5Copy, s5)
	f1(s5Copy)
	printlnSlice("s5", s5)
	printlnSlice("s5", s5)

	// second option - full slice expression
	s6 := []int{1, 2, 3}
	printlnSlice("s6", s6)
	f1(s6[:2:2])
	printlnSlice("s6", s6)
}

type Foo struct {
	v []byte
}

func printAlloc() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%d KB\n", m.Alloc/1024)
}

func keepFirstTwoElementsOnly(foos []Foo) []Foo {
	// not correct
	//return foos[:2]

	// correct
	res := make([]Foo, 2)
	copy(res, foos)
	return res
}

// Mistake26 demonstrates memory allocation and deallocation issues when working with slices in Go.
func Mistake26() {
	//
	foos := make([]Foo, 1_000)
	printAlloc()

	for i := 0; i < len(foos); i++ {
		foos[i] = Foo{
			v: make([]byte, 1024*1024),
		}
	}
	printAlloc()

	two := keepFirstTwoElementsOnly(foos)
	runtime.GC()
	printAlloc()
	runtime.KeepAlive(two)
}
