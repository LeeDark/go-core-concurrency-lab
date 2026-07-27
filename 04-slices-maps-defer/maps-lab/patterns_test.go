package mapslab

import (
	"reflect"
	"testing"
)

func TestKeysContainsEveryMapKeyWithoutAssumingOrder(t *testing.T) {
	values := map[string]int{"go": 1, "map": 2, "slice": 3}
	keys := Keys(values)

	if len(keys) != len(values) {
		t.Fatalf("len(Keys(values)) = %d, want %d", len(keys), len(values))
	}
	for _, key := range keys {
		if _, ok := values[key]; !ok {
			t.Fatalf("Keys(values) contains unexpected key %q", key)
		}
	}
}

func TestSortedKeysReturnsDeterministicOrder(t *testing.T) {
	values := map[string]int{"slice": 1, "go": 2, "map": 3}

	got := SortedKeys(values)
	want := []string{"go", "map", "slice"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SortedKeys(values) = %v, want %v", got, want)
	}
}

func TestCountWords(t *testing.T) {
	got := CountWords([]string{"go", "map", "go", "slice", "map", "go"})
	want := map[string]int{"go": 3, "map": 2, "slice": 1}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("CountWords() = %v, want %v", got, want)
	}
}

func TestGroupByLengthPreservesOrderWithinGroups(t *testing.T) {
	got := GroupByLength([]string{"go", "map", "to", "slice"})
	want := map[int][]string{
		2: {"go", "to"},
		3: {"map"},
		5: {"slice"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GroupByLength() = %v, want %v", got, want)
	}
}

func TestIndexByIDKeepsLaterUserForDuplicateID(t *testing.T) {
	got := IndexByID([]User{
		{ID: 10, Name: "Ana"},
		{ID: 20, Name: "Bohdan"},
		{ID: 10, Name: "Andrii"},
	})
	want := map[int]User{
		10: {ID: 10, Name: "Andrii"},
		20: {ID: 20, Name: "Bohdan"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("IndexByID() = %v, want %v", got, want)
	}
}

func TestStringSetRemovesDuplicatesAndChecksMembership(t *testing.T) {
	set := NewStringSet([]string{"go", "map", "go"})

	if got := len(set); got != 2 {
		t.Fatalf("len(set) = %d, want 2", got)
	}
	if !HasString(set, "go") {
		t.Fatal("HasString(set, go) = false, want true")
	}
	if HasString(set, "slice") {
		t.Fatal("HasString(set, slice) = true, want false")
	}
}
