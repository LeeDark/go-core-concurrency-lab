package slicesadvanced

import "fmt"

// MemoryLeakSubslice compares retaining a large backing array with copying the needed data.
func MemoryLeakSubslice() {
	fmt.Printf("\nMemory Leak with Subslice\n\n")

	bad := TakeSmallPart()
	goodCopy := TakeSmallPartSafeCopy()
	goodAppend := TakeSmallPartSafeAppend()

	printlnSlice("bad", bad)
	printlnSlice("goodCopy", goodCopy)
	printlnSlice("goodAppend", goodAppend)
}
