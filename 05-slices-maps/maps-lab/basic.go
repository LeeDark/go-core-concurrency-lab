// Package mapslab contains focused examples for Go maps.
package mapslab

// NewInventoryLiteral creates an inventory using a map literal.
func NewInventoryLiteral() map[string]int {
	return map[string]int{
		"apples":  10,
		"bananas": 5,
	}
}

// NewInventory creates an empty inventory that is ready for writes.
func NewInventory() map[string]int {
	return make(map[string]int)
}

// Quantity returns the quantity stored under item. A missing item returns the zero value for int.
func Quantity(inventory map[string]int, item string) int {
	return inventory[item]
}

// SetQuantity inserts item when it is absent and updates it when it is already present.
func SetQuantity(inventory map[string]int, item string, quantity int) {
	inventory[item] = quantity
}

// RemoveItem deletes item. Deleting a missing item has no effect.
func RemoveItem(inventory map[string]int, item string) {
	delete(inventory, item)
}
