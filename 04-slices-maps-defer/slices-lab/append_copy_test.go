package sliceslab

import (
	"reflect"
	"testing"
)

func TestCloneSliceCopyPreservesNil(t *testing.T) {
	got := CloneSliceCopy[int](nil)

	if got != nil {
		t.Fatalf("CloneSliceCopy(nil) = %v, want nil", got)
	}
}

func TestCloneSliceCopyIsIndependent(t *testing.T) {
	source := []int{1, 2, 3}
	clone := CloneSliceCopy(source)
	clone[0] = 99

	if !reflect.DeepEqual(source, []int{1, 2, 3}) {
		t.Fatalf("source = %v, want unchanged values", source)
	}
	if !reflect.DeepEqual(clone, []int{99, 2, 3}) {
		t.Fatalf("clone = %v, want [99 2 3]", clone)
	}
}

func TestCloneSliceAppendConvertsEmptyToNil(t *testing.T) {
	got := CloneSliceAppend([]int{})

	if got != nil {
		t.Fatalf("CloneSliceAppend(empty) = %v, want nil", got)
	}
}

func TestAppendSafeDoesNotModifyInput(t *testing.T) {
	input := make([]int, 2, 4)
	input[0], input[1] = 1, 2

	got := AppendSafe(input, 3)
	got[0] = 99

	if !reflect.DeepEqual(input, []int{1, 2}) {
		t.Fatalf("input = %v, want [1 2]", input)
	}
	if !reflect.DeepEqual(got, []int{99, 2, 3}) {
		t.Fatalf("result = %v, want [99 2 3]", got)
	}
}

func TestDeleteAtModifiesInputAndClearsTail(t *testing.T) {
	input := []int{10, 20, 30, 40}

	got := DeleteAt(input, 1)

	if !reflect.DeepEqual(got, []int{10, 30, 40}) {
		t.Fatalf("result = %v, want [10 30 40]", got)
	}
	if !reflect.DeepEqual(input, []int{10, 30, 40, 0}) {
		t.Fatalf("input backing array = %v, want [10 30 40 0]", input)
	}
}

func TestDeleteAtNewDoesNotModifyInput(t *testing.T) {
	input := []int{10, 20, 30, 40}

	got := DeleteAtNew(input, 1)
	got[0] = 99

	if !reflect.DeepEqual(input, []int{10, 20, 30, 40}) {
		t.Fatalf("input = %v, want unchanged values", input)
	}
	if !reflect.DeepEqual(got, []int{99, 30, 40}) {
		t.Fatalf("result = %v, want [99 30 40]", got)
	}
}

func TestFilterInPlaceKeepsValuesAndClearsTail(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6}

	got := FilterInPlace(input, func(n int) bool {
		return n%2 == 0
	})

	if !reflect.DeepEqual(got, []int{2, 4, 6}) {
		t.Fatalf("result = %v, want [2 4 6]", got)
	}
	if !reflect.DeepEqual(input, []int{2, 4, 6, 0, 0, 0}) {
		t.Fatalf("input backing array = %v, want [2 4 6 0 0 0]", input)
	}
}
