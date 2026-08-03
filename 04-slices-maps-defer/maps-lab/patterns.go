package mapslab

import (
	"fmt"
	"sort"
)

// User is a minimal record used to demonstrate indexing a slice by ID.
type User struct {
	ID   int
	Name string
}

// Keys returns the keys collected with range. Their order is not stable.
func Keys(values map[string]int) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}

// SortedKeys returns map keys in ascending order for deterministic output.
func SortedKeys(values map[string]int) []string {
	keys := Keys(values)
	sort.Strings(keys)
	return keys
}

// CountWords counts how often each word occurs.
func CountWords(words []string) map[string]int {
	counts := make(map[string]int)
	for _, word := range words {
		counts[word]++
	}
	return counts
}

// GroupByLength groups words by their byte length while preserving input order within each group.
func GroupByLength(words []string) map[int][]string {
	groups := make(map[int][]string)
	for _, word := range words {
		length := len(word)
		groups[length] = append(groups[length], word)
	}
	return groups
}

// IndexByID indexes users by ID. When IDs repeat, the later user replaces the earlier one.
func IndexByID(users []User) map[int]User {
	index := make(map[int]User, len(users))
	for _, user := range users {
		index[user.ID] = user
	}
	return index
}

// NewStringSet builds a set from values. struct{} uses no storage per set value.
func NewStringSet(values []string) map[string]struct{} {
	set := make(map[string]struct{}, len(values))
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

// HasString reports whether value belongs to set.
func HasString(set map[string]struct{}, value string) bool {
	_, ok := set[value]
	return ok
}

// MapPatterns demonstrates iteration order and common map patterns.
func MapPatterns() {
	fmt.Println("\nMap iteration and patterns")

	counts := CountWords([]string{"go", "map", "go", "slice"})
	fmt.Println("keys from range (order is not stable):", Keys(counts))
	for _, word := range SortedKeys(counts) {
		fmt.Println("count", word+":", counts[word])
	}

	groups := GroupByLength([]string{"go", "map", "to", "slice"})
	fmt.Println("words with length 2:", groups[2])
	fmt.Println("words with length 3:", groups[3])

	users := IndexByID([]User{{ID: 10, Name: "Ana"}, {ID: 20, Name: "Bohdan"}})
	fmt.Println("user with ID 20:", users[20].Name)

	set := NewStringSet([]string{"go", "map", "go"})
	fmt.Println("set has map:", HasString(set, "map"))
	fmt.Println("set size:", len(set))
}
