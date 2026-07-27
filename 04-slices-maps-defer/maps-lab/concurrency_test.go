package mapslab

import (
	"fmt"
	"sync"
	"testing"
)

func TestSafeInventoryBasicOperations(t *testing.T) {
	inventory := NewSafeInventory()
	inventory.Set("apples", 10)

	if quantity, ok := inventory.Lookup("apples"); quantity != 10 || !ok {
		t.Fatalf("Lookup(apples) = (%d, %t), want (10, true)", quantity, ok)
	}

	inventory.Delete("apples")
	if quantity, ok := inventory.Lookup("apples"); quantity != 0 || ok {
		t.Fatalf("Lookup(apples) after delete = (%d, %t), want (0, false)", quantity, ok)
	}
}

func TestSafeInventorySupportsConcurrentReadsAndWrites(t *testing.T) {
	const (
		readerCount = 20
		writerCount = 40
		itemCount   = 10
	)

	inventory := NewSafeInventory()
	for i := 0; i < itemCount; i++ {
		inventory.Set(fmt.Sprintf("readable-%d", i), i)
	}

	start := make(chan struct{})
	var workers sync.WaitGroup

	for i := 0; i < writerCount; i++ {
		workers.Add(1)
		go func(i int) {
			defer workers.Done()
			<-start
			inventory.Set(fmt.Sprintf("written-%d", i), i)
		}(i)
	}

	for reader := 0; reader < readerCount; reader++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			<-start
			for i := 0; i < itemCount; i++ {
				quantity, ok := inventory.Lookup(fmt.Sprintf("readable-%d", i))
				if quantity != i || !ok {
					t.Errorf("Lookup(readable-%d) = (%d, %t), want (%d, true)", i, quantity, ok, i)
				}
			}
		}()
	}

	close(start)
	workers.Wait()

	if got, want := inventory.Len(), itemCount+writerCount; got != want {
		t.Fatalf("Len() = %d, want %d", got, want)
	}
}
