package main

import (
	mapslab "github.com/LeeDark/go-core-concurrency-lab/04-slices-maps-defer/maps-lab"
	sliceslab "github.com/LeeDark/go-core-concurrency-lab/04-slices-maps-defer/slices-lab"
)

func main() {
	sliceslab.ArraySlice()
	sliceslab.SliceLenCap()
	sliceslab.SliceAppend()
	sliceslab.AliasingBug()
	sliceslab.FullSliceExpression()
	sliceslab.AppendIntoFunc()

	sliceslab.CopyAppend()
	sliceslab.MemoryLeakSubslice()
	sliceslab.SliceNilEmpty()
	sliceslab.RangeOverSlice()

	sliceslab.Mistake20()
	sliceslab.Mistake21()
	sliceslab.Mistake22()
	sliceslab.Mistake23()
	sliceslab.Mistake24()
	sliceslab.Mistake25()
	sliceslab.Mistake26()

	mapslab.MapBasics()
	mapslab.MapState()
}
