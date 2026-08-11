package mapslab

import "fmt"

// MapBasics demonstrates creating, reading, inserting, updating, and deleting map entries.
func MapBasics() {
	fmt.Println("\nMap basics")

	inventory := NewInventoryLiteral()
	fmt.Println("literal, apples:", Quantity(inventory, "apples"))
	fmt.Println("literal length:", len(inventory))

	emptyInventory := NewInventory()
	fmt.Println("make length:", len(emptyInventory))

	SetQuantity(inventory, "oranges", 7)
	fmt.Println("after insert, oranges:", Quantity(inventory, "oranges"))
	fmt.Println("after insert length:", len(inventory))

	SetQuantity(inventory, "apples", 12)
	fmt.Println("after update, apples:", Quantity(inventory, "apples"))
	fmt.Println("after update length:", len(inventory))

	fmt.Println("missing pears:", Quantity(inventory, "pears"))

	RemoveItem(inventory, "bananas")
	fmt.Println("after delete length:", len(inventory))
}
