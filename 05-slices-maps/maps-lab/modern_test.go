package mapslab

import (
	"reflect"
	"testing"
)

func TestClearInventoryKeepsMapWritable(t *testing.T) {
	inventory := NewInventoryLiteral()
	ClearInventory(inventory)

	if len(inventory) != 0 {
		t.Fatalf("len after ClearInventory = %d, want 0", len(inventory))
	}
	SetQuantity(inventory, "apples", 10)
	if got := Quantity(inventory, "apples"); got != 10 {
		t.Fatalf("quantity after clear and write = %d, want 10", got)
	}
}

func TestClearInventoryIsSafeForNilMap(t *testing.T) {
	var inventory map[string]int

	ClearInventory(inventory)
}

func TestCloneInventoryIsIndependentAndPreservesNil(t *testing.T) {
	var nilInventory map[string]int
	if got := CloneInventory(nilInventory); got != nil {
		t.Fatalf("CloneInventory(nil) = %v, want nil", got)
	}

	original := NewInventoryLiteral()
	clone := CloneInventory(original)
	clone["apples"] = 99

	if reflect.DeepEqual(original, clone) {
		t.Fatal("CloneInventory returned a shared map")
	}
	if got := original["apples"]; got != 10 {
		t.Fatalf("original apples after clone mutation = %d, want 10", got)
	}
}

func TestCopyInventoryOverwritesDuplicateKeys(t *testing.T) {
	destination := map[string]int{"apples": 10, "pears": 2}
	source := map[string]int{"apples": 12, "bananas": 5}

	CopyInventory(destination, source)

	want := map[string]int{"apples": 12, "pears": 2, "bananas": 5}
	if !reflect.DeepEqual(destination, want) {
		t.Fatalf("destination after CopyInventory = %v, want %v", destination, want)
	}
}

func TestUpdateUserNameCopiesMapValueBack(t *testing.T) {
	index := map[int]User{10: {ID: 10, Name: "Ana"}}

	if !UpdateUserName(index, 10, "Anna") {
		t.Fatal("UpdateUserName(existing) = false, want true")
	}
	if got := index[10].Name; got != "Anna" {
		t.Fatalf("updated name = %q, want Anna", got)
	}
	if UpdateUserName(index, 20, "Bohdan") {
		t.Fatal("UpdateUserName(missing) = true, want false")
	}
}
