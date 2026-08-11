package main

import (
	mapslab "github.com/LeeDark/go-core-concurrency-lab/05-slices-maps/maps-lab"
	sliceslab "github.com/LeeDark/go-core-concurrency-lab/05-slices-maps/slices-lab"
)

func main() {
	sliceslab.ArraySlice()
	sliceslab.SliceLenCap()
	sliceslab.SliceAppend()
	sliceslab.AliasingBug()
	sliceslab.FullSliceExpression()
	sliceslab.AppendIntoFunc()

	sliceslab.CopyAppend()
	sliceslab.SliceNilEmpty()
	sliceslab.RangeOverSlice()

	mapslab.MapBasics()
	mapslab.MapState()
	mapslab.MapPatterns()
	mapslab.MapConcurrency()
}
