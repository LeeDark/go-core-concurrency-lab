package sliceslab

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ArraySlice array, slice
func ArraySlice() {
	fmt.Printf("Array and Slice basics\n\n")

	// Array
	fmt.Println("Array")
	// An array has a fixed size, and it's length is part of its type, so arrays cannot be resized.
	var a [10]int = [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println(a)

	// Slice
	fmt.Println("Slice")
	// A slice is a dynamically sized, flexible view into the elements of an array.
	b := a[3:8]
	fmt.Println(b)

	// Slice with underlying array
	// A slice does not store any data, it just describes a section of an underlying array.
	fmt.Println("Slice with underlying array")
	c := []int{2, 4, 6, 8, 10, 12, 14}
	fmt.Println(c)

	// When slicing, you may omit the high or low bounds to use their defaults instead.
	fmt.Println("Slicing")
	d := c[3:] // {8, 10, 12, 14}
	fmt.Println(d)

	// Changing the elements of a slice modifies the corresponding elements of its underlying array.
	fmt.Println("Slice changing")
	d[2] = 100
	fmt.Println(d) // {8, 10, 100, 14}
	printSlice(d)
	fmt.Println(c) // {2, 4, 6, 8, 10, 100, 14}
	printSlice(c)
}

// SliceLenCap len, cap
func SliceLenCap() {
	fmt.Printf("\nSlice length and capacity\n\n")

	s := []int{2, 3, 5, 7, 11, 13}
	printSlice(s)

	// Slice the slice to give it zero length.
	s = s[:0]
	printSlice(s)

	// Extend its length.
	s = s[:4]
	printSlice(s)

	// Drop its first two values.
	s = s[2:]
	printSlice(s)

	a := make([]int, 5)
	printlnSlice("a", a)

	b := make([]int, 0, 5)
	printlnSlice("b", b)

	c := b[:2]
	printlnSlice("c", c)

	d := c[2:5]
	printlnSlice("d", d)

	// Create a tic-tac-toe board.
	board := [][]string{
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
	}

	// The players take turns.
	board[0][0] = "X"
	board[2][2] = "O"
	board[1][2] = "X"
	board[1][0] = "O"
	board[0][2] = "X"

	for i := 0; i < len(board); i++ {
		fmt.Printf("%s\n", strings.Join(board[i], " "))
	}
}

// SliceAppend 1: append
func SliceAppend() {
	fmt.Printf("\nAppending to a Slice\n\n")

	var s []int
	printSlice(s)

	// append works on nil slices.
	s = append(s, 0)
	printSlice(s)

	// The slice grows as needed.
	s = append(s, 1)
	printSlice(s)

	// We can add more than one element at a time.
	s = append(s, 2, 3, 4)
	printSlice(s)

	fmt.Printf("\nappend changes same \"backing array\"\n\n")
	// append change same "backing array"
	a := make([]int, 2, 4)
	a[0] = 10
	a[1] = 20

	fmt.Println("before:")
	printlnSlice("a", a)

	b := append(a, 30)

	fmt.Println("after append:")
	printlnSlice("a", a)
	printlnSlice("b", b)

	b[0] = 999

	fmt.Println("after b[0] = 999:")
	fmt.Println("a:", a)
	fmt.Println("b:", b)

	fmt.Printf("\nappend creates new \"backing array\"\n\n")
	// append creates new "backing array"
	c := make([]int, 2, 2)
	c[0] = 10
	c[1] = 20

	fmt.Println("before:")
	printlnSlice("c", c)

	d := append(c, 30)

	fmt.Println("after append:")
	printlnSlice("c", c)
	printlnSlice("d", d)

	d[0] = 999

	fmt.Println("after d[0] = 999:")
	fmt.Println("c:", c)
	fmt.Println("d:", d)

	fmt.Printf("\nappend changes another Slice\n\n")
	// append changes another Slice
	e := []int{1, 2, 3, 4}

	f := e[:2] // [1 2], cap = 4
	g := e[2:] // [3 4], cap = 2

	fmt.Println("before:")
	printlnSlice("e", e)
	printlnSlice("f", f)
	printlnSlice("g", g)

	f = append(f, 99)

	fmt.Println("after append(f, 99):")
	printlnSlice("e", e)
	printlnSlice("f", f)
	printlnSlice("g", g)
}

// AliasingBug 1: subslice aliasing bug
func AliasingBug() {
	fmt.Printf("\nSubslice aliasing bug\n\n")

	a := []int{1, 2, 3, 4}
	printlnSlice("a", a)

	b := a[:2]
	printlnSlice("b", b)

	b = append(b, 99)
	printlnSlice("a after append", a)
	printlnSlice("b after append", b)
}

// FullSliceExpression 1: full slice expression
func FullSliceExpression() {
	fmt.Printf("\nFull Slice expression\n\n")

	a := []int{1, 2, 3, 4}

	b := a[:2:2] // len = 2, cap = 2
	c := a[2:]

	fmt.Println("before:")
	printlnSlice("a", a)
	printlnSlice("b", b)
	printlnSlice("c", c)

	b = append(b, 99)

	fmt.Println("after append(b, 99):")
	printlnSlice("a", a)
	printlnSlice("b", b)
	printlnSlice("c", c)
}

// AppendIntoFunc 1: append into func
func AppendIntoFunc() {
	addOne := func(s []int) []int {
		s = append(s, 100)
		s[0] = 999
		return s
	}

	fmt.Printf("\nappend into func\n\n")

	a := make([]int, 2, 4)
	a[0] = 1
	a[1] = 2

	fmt.Println("before:")
	printlnSlice("a", a)

	b := addOne(a)

	fmt.Println("after:")
	printlnSlice("a", a)
	printlnSlice("b", b)
}

// CopyAppend 2: copy vs append
func CopyAppend() {
	fmt.Printf("\ncopy\n\n")
	src := []int{1, 2, 3, 4}
	dst := make([]int, 2)
	num := copy(dst, src)

	fmt.Println("num:", num)
	printlnSlice("src", src)
	printlnSlice("dst", dst)

	fmt.Printf("\nCloneSlice with copy\n\n")
	a := []int{1, 2, 3}
	b := CloneSliceCopy(a)
	b[0] = 999
	printlnSlice("a", a)
	printlnSlice("b", b)

	fmt.Printf("\nCloneSlice with append\n\n")
	c := []int{1, 2, 3}
	d := append([]int(nil), c...)
	d[0] = 777
	printlnSlice("c", c)
	printlnSlice("d", d)

	fmt.Printf("\nCloneSlice with append - empty to nil Slice!\n\n")
	empty := []int{}
	cloned := append([]int(nil), empty...)
	printlnSlice("empty", empty)
	printlnSlice("cloned", cloned)
	fmt.Println(empty == nil)
	fmt.Println(cloned == nil)

	fmt.Printf("\nAppendSafe with append\n\n")
	e := make([]int, 2, 4)
	e[0] = 1
	e[1] = 2
	f := AppendSafe(e, 3)
	f[0] = 999
	printlnSlice("e", e)
	printlnSlice("f", f)

	fmt.Printf("\nDeleteAt in-place, but change original array\n\n")
	g := []int{10, 20, 30, 40}
	h := DeleteAt(g, 1)
	printlnSlice("g", g)
	printlnSlice("h", h)

	fmt.Printf("\nDeleteAtNew with new Slice\n\n")
	i := []int{10, 20, 30, 40}
	j := DeleteAtNew(i, 1)
	j[0] = 999
	printlnSlice("i", i)
	printlnSlice("j", j)

	fmt.Printf("\nFilterInPlace in-place, but change original array\n\n")
	k := []int{1, 2, 3, 4, 5, 6}
	l := FilterInPlace(k, func(x int) bool {
		return x%2 == 0
	})
	printlnSlice("k", k)
	printlnSlice("l", l)

	fmt.Printf("\nFilterNewSlice with new Slice\n\n")
	m := []int{1, 2, 3, 4, 5, 6}
	n := FilterNewSlice(m, func(x int) bool {
		return x%2 == 0
	})
	n[0] = 999
	printlnSlice("m", m)
	printlnSlice("n", n)
}

// MemoryLeakSubslice 2: memory leak with subslice
func MemoryLeakSubslice() {
	fmt.Printf("\nMemory Leak with Subslice\n\n")

	bad := TakeSmallPart()
	goodCopy := TakeSmallPartSafeCopy()
	goodAppend := TakeSmallPartSafeAppend()

	printlnSlice("bad", bad)
	printlnSlice("goodCopy", goodCopy)
	printlnSlice("goodAppend", goodAppend)
}

// SliceNilEmpty 2: nil slice vs empty slice
func SliceNilEmpty() {
	fmt.Printf("\nnil Slice vs Empty Slice\n\n")

	var nilSlice []int
	emptySlice := []int{}
	madeEmptySlice := make([]int, 0)

	fmt.Println("nilSlice:")
	printSliceFull(nilSlice)

	fmt.Println("emptySlice:")
	printSliceFull(emptySlice)

	fmt.Println("madeEmptySlice:")
	printSliceFull(madeEmptySlice)

	fmt.Println("append:")
	nilSlice = append(nilSlice, 1)
	emptySlice = append(emptySlice, 1)

	fmt.Println("nilSlice after append:", nilSlice)
	fmt.Println("emptySlice after append:", emptySlice)

	nilJSON, _ := json.Marshal([]int(nil))
	emptyJSON, _ := json.Marshal([]int{})

	fmt.Println("json nil slice:", string(nilJSON))
	fmt.Println("json empty slice:", string(emptyJSON))
}

type User struct {
	Name string
	Age  int
}

// RangeOverSlice 2: range over slice mistakes
func RangeOverSlice() {
	// 1
	// incorrect
	fmt.Printf("\nIncorrect change with \"element copy\"\n\n")
	s1 := []int{1, 2, 3}
	for _, v := range s1 {
		v *= 10
	}
	fmt.Println(s1)
	// correct
	fmt.Printf("\nCorrect change with \"by index\"\n\n")
	s2 := []int{1, 2, 3}
	for i := range s2 {
		s2[i] *= 10
	}
	fmt.Println(s2)

	// 2
	// incorrect
	fmt.Printf("\nIncorrect change struct with \"element copy\"\n\n")
	users1 := []User{
		{Name: "Bob", Age: 20},
		{Name: "Alice", Age: 25},
	}
	for _, u := range users1 {
		u.Age++
	}
	fmt.Println(users1)
	// correct
	fmt.Printf("\nCorrect change struct with \"by index\"\n\n")
	users2 := []User{
		{Name: "Bob", Age: 20},
		{Name: "Alice", Age: 25},
	}
	for i := range users2 {
		users2[i].Age++
	}
	fmt.Println(users2)

	// 3
	fmt.Printf("\nappend into range\n\n")
	s3 := []int{1, 2, 3}
	for _, v := range s3 {
		fmt.Println("v:", v)
		s3 = append(s3, v*10)
	}
	fmt.Println("result:", s3)

	s4 := make([]int, 3, 6)
	s4[0] = 1
	s4[1] = 2
	s4[2] = 3
	for i, v := range s4 {
		s4 = append(s4, v*10)
		s4[i] = 999
	}
	fmt.Println("result:", s4)

	s5 := make([]int, 3, 3)
	s5[0] = 1
	s5[1] = 2
	s5[2] = 3
	for i, v := range s5 {
		s5 = append(s5, v*10)
		s5[i] = 999
	}
	fmt.Println("result:", s5)
}

// 3: 100 Go Mistakes
