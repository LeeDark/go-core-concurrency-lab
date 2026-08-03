package mapslab

import "testing"

func TestNewInventoryLiteral(t *testing.T) {
	inventory := NewInventoryLiteral()

	if got := Quantity(inventory, "apples"); got != 10 {
		t.Fatalf("apples = %d, want 10", got)
	}
	if got := Quantity(inventory, "bananas"); got != 5 {
		t.Fatalf("bananas = %d, want 5", got)
	}
}

func TestNewInventoryIsEmptyAndWritable(t *testing.T) {
	inventory := NewInventory()

	if inventory == nil {
		t.Fatal("NewInventory() returned nil, want initialized map")
	}
	if len(inventory) != 0 {
		t.Fatalf("len(NewInventory()) = %d, want 0", len(inventory))
	}

	SetQuantity(inventory, "apples", 10)
	if got := Quantity(inventory, "apples"); got != 10 {
		t.Fatalf("apples after insert = %d, want 10", got)
	}
}

func TestSetQuantityUpdatesExistingItemWithoutChangingLength(t *testing.T) {
	inventory := NewInventoryLiteral()
	before := len(inventory)

	SetQuantity(inventory, "apples", 12)

	if got := Quantity(inventory, "apples"); got != 12 {
		t.Fatalf("apples after update = %d, want 12", got)
	}
	if got := len(inventory); got != before {
		t.Fatalf("len after update = %d, want %d", got, before)
	}
}

func TestQuantityReturnsZeroValueForMissingItem(t *testing.T) {
	inventory := NewInventoryLiteral()

	if got := Quantity(inventory, "pears"); got != 0 {
		t.Fatalf("missing pears = %d, want zero value 0", got)
	}
}

func TestRemoveItem(t *testing.T) {
	inventory := NewInventoryLiteral()

	RemoveItem(inventory, "bananas")
	if got := len(inventory); got != 1 {
		t.Fatalf("len after delete = %d, want 1", got)
	}
	if got := Quantity(inventory, "bananas"); got != 0 {
		t.Fatalf("bananas after delete = %d, want zero value 0", got)
	}

	RemoveItem(inventory, "bananas")
	if got := len(inventory); got != 1 {
		t.Fatalf("len after deleting missing item = %d, want 1", got)
	}
}
