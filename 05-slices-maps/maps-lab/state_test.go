package mapslab

import "testing"

func TestLookupQuantityDistinguishesStoredZeroFromMissingItem(t *testing.T) {
	inventory := NewInventory()
	SetQuantity(inventory, "apples", 0)

	quantity, ok := LookupQuantity(inventory, "apples")
	if quantity != 0 || !ok {
		t.Fatalf("LookupQuantity(apples) = (%d, %t), want (0, true)", quantity, ok)
	}

	quantity, ok = LookupQuantity(inventory, "pears")
	if quantity != 0 || ok {
		t.Fatalf("LookupQuantity(pears) = (%d, %t), want (0, false)", quantity, ok)
	}
}

func TestNilInventoryCanBeRead(t *testing.T) {
	inventory := NilInventory()

	if inventory != nil {
		t.Fatal("NilInventory() is not nil")
	}
	if got := len(inventory); got != 0 {
		t.Fatalf("len(nil inventory) = %d, want 0", got)
	}
	if got := Quantity(inventory, "apples"); got != 0 {
		t.Fatalf("Quantity(nil inventory, apples) = %d, want 0", got)
	}

	quantity, ok := LookupQuantity(inventory, "apples")
	if quantity != 0 || ok {
		t.Fatalf("LookupQuantity(nil inventory, apples) = (%d, %t), want (0, false)", quantity, ok)
	}
}

func TestEmptyInventoryIsWritable(t *testing.T) {
	inventory := NewInventory()

	if inventory == nil {
		t.Fatal("NewInventory() returned nil, want empty initialized map")
	}
	SetQuantity(inventory, "apples", 10)

	if quantity, ok := LookupQuantity(inventory, "apples"); quantity != 10 || !ok {
		t.Fatalf("LookupQuantity(apples) = (%d, %t), want (10, true)", quantity, ok)
	}
}

func TestSetQuantityPanicsForNilInventory(t *testing.T) {
	inventory := NilInventory()

	defer func() {
		if recover() == nil {
			t.Fatal("SetQuantity(nil inventory) did not panic")
		}
	}()

	SetQuantity(inventory, "apples", 10)
}
