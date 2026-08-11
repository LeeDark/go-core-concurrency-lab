package mapslab

import "maps"

// ClearInventory removes all entries while keeping the map writable.
func ClearInventory(inventory map[string]int) {
	clear(inventory)
}

// CloneInventory makes a shallow copy using the standard library.
func CloneInventory(inventory map[string]int) map[string]int {
	return maps.Clone(inventory)
}

// CopyInventory copies entries from source into destination, replacing duplicate keys.
func CopyInventory(destination, source map[string]int) {
	maps.Copy(destination, source)
}

// UpdateUserName demonstrates copy-update-writeback for a struct map value.
func UpdateUserName(index map[int]User, id int, name string) bool {
	user, ok := index[id]
	if !ok {
		return false
	}

	user.Name = name
	index[id] = user
	return true
}
