package mapslab

import (
	"fmt"
	"sync"
)

// SafeInventory protects its map with an RWMutex.
type SafeInventory struct {
	mu         sync.RWMutex
	quantities map[string]int
}

// NewSafeInventory creates an empty inventory ready for concurrent access.
func NewSafeInventory() *SafeInventory {
	return &SafeInventory{quantities: make(map[string]int)}
}

// Set inserts or updates item while holding the write lock.
func (inventory *SafeInventory) Set(item string, quantity int) {
	inventory.mu.Lock()
	defer inventory.mu.Unlock()

	inventory.quantities[item] = quantity
}

// Lookup reads item while holding the read lock.
func (inventory *SafeInventory) Lookup(item string) (int, bool) {
	inventory.mu.RLock()
	defer inventory.mu.RUnlock()

	quantity, ok := inventory.quantities[item]
	return quantity, ok
}

// Delete removes item while holding the write lock.
func (inventory *SafeInventory) Delete(item string) {
	inventory.mu.Lock()
	defer inventory.mu.Unlock()

	delete(inventory.quantities, item)
}

// Len returns the number of items while holding the read lock.
func (inventory *SafeInventory) Len() int {
	inventory.mu.RLock()
	defer inventory.mu.RUnlock()

	return len(inventory.quantities)
}

// MapConcurrency demonstrates safe map access through SafeInventory.
func MapConcurrency() {
	fmt.Println("\nMap concurrency boundary")

	inventory := NewSafeInventory()
	inventory.Set("apples", 10)

	quantity, ok := inventory.Lookup("apples")
	fmt.Println("apples:", quantity, "present:", ok)

	inventory.Delete("apples")
	_, ok = inventory.Lookup("apples")
	fmt.Println("apples after delete, present:", ok)
}
