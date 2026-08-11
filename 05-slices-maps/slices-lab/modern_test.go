package sliceslab

import (
	"reflect"
	"testing"
)

func TestClearSliceKeepsShapeAndZeroesElements(t *testing.T) {
	input := make([]string, 2, 4)
	input[0], input[1] = "go", "slice"

	got := ClearSlice(input)

	if len(got) != 2 || cap(got) != 4 {
		t.Fatalf("shape after ClearSlice = len %d cap %d, want len 2 cap 4", len(got), cap(got))
	}
	if !reflect.DeepEqual(got, []string{"", ""}) {
		t.Fatalf("ClearSlice(input) = %v, want zeroed elements", got)
	}
}

func TestCloneSliceStdlibIsIndependentAndPreservesNil(t *testing.T) {
	var nilInput []int
	if got := CloneSliceStdlib(nilInput); got != nil {
		t.Fatalf("CloneSliceStdlib(nil) = %v, want nil", got)
	}

	input := []int{1, 2, 3}
	clone := CloneSliceStdlib(input)
	clone[0] = 99

	if !reflect.DeepEqual(input, []int{1, 2, 3}) {
		t.Fatalf("input after clone mutation = %v, want unchanged", input)
	}
}

func TestDeleteSliceStdlibClearsUnusedTail(t *testing.T) {
	first, second := "first", "second"
	input := []*string{&first, &second}

	got := DeleteSliceStdlib(input, 0, 1)

	if len(got) != 1 || *got[0] != "second" {
		t.Fatalf("DeleteSliceStdlib result = %v, want one remaining element", got)
	}
	if input[1] != nil {
		t.Fatal("DeleteSliceStdlib did not clear the unused tail")
	}
}
