package mapslab

import "fmt"

// LookupQuantity returns the quantity and whether item is present in inventory.
func LookupQuantity(inventory map[string]int, item string) (int, bool) {
	quantity, ok := inventory[item]
	return quantity, ok
}

// NilInventory returns a nil map. It is safe to read from but not to write to.
func NilInventory() map[string]int {
	var inventory map[string]int
	return inventory
}

// MapState demonstrates comma-ok lookups and the difference between nil and empty maps.
func MapState() {
	fmt.Println("\nMap lookups and state")

	inventory := NewInventory()
	SetQuantity(inventory, "apples", 0)

	quantity, ok := LookupQuantity(inventory, "apples")
	fmt.Println("stored zero, apples:", quantity, "present:", ok)

	quantity, ok = LookupQuantity(inventory, "pears")
	fmt.Println("missing pears:", quantity, "present:", ok)

	nilInventory := NilInventory()
	emptyInventory := NewInventory()
	fmt.Println("nil inventory:", nilInventory == nil, "length:", len(nilInventory))
	fmt.Println("empty inventory:", emptyInventory == nil, "length:", len(emptyInventory))

	quantity, ok = LookupQuantity(nilInventory, "apples")
	fmt.Println("nil lookup, apples:", quantity, "present:", ok)
}
