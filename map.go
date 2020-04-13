package siftr

import (
	"strings"
)

// Map is a wrapper around map[string]string that provides some helpers
// above it that assume the map is in the format that flatmap expects
// (the result of Flatten).
//
// All modifying functions such as Delete are done in-place unless
// otherwise noted.
type Map map[string]string

// Contains returns true if the map contains the given key.
func (m Map) Contains(key string) bool {
	for _, k := range m.Keys() {
		if k == key {
			return true
		}
	}

	return false
}

// Delete deletes a key out of the map with the given prefix.
func (m Map) Delete(prefix string) {
	deletePrefix(m, prefix)
}

func deletePrefix(m map[string]string, prefix string) {
	for k := range m {
		match := k == prefix
		if !match {
			if !strings.HasPrefix(k, prefix) {
				continue
			}

			if k[len(prefix):len(prefix)+1] != "." {
				continue
			}
		}

		delete(m, k)
	}
}

// Keys returns all of the top-level keys in this map
func (m Map) Keys() []string {
	ks := make(map[string]struct{})
	for k := range m {
		idx := strings.Index(k, ".")
		if idx == -1 {
			idx = len(k)
		}

		ks[k[:idx]] = struct{}{}
	}

	result := make([]string, 0, len(ks))
	for k := range ks {
		result = append(result, k)
	}

	return result
}
