package siftr

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Expand takes a map and a key (prefix) and expands that value into
// a more complex structure. This is the reverse of the Flatten operation.
func Expand(m map[string]string, key string) interface{} {
	// If the key is exactly a key in the map, just return it
	if v, ok := m[key]; ok {
		if v == "true" {
			return true
		} else if v == "false" {
			return false
		}

		return v
	}

	// Check if the key is an array, and if so, expand the array
	if _, ok := m[key+".#"]; ok {
		return expandArray(m, key)
	}

	// Check if this is a prefix in the map
	prefix := key + "."
	for k := range m {
		if strings.HasPrefix(k, prefix) {
			return expandMap(m, prefix)
		}
	}

	return nil
}

func expandArray(m map[string]string, prefix string) []interface{} {
	//start := time.Now()
	num, err := strconv.ParseInt(m[prefix+".#"], 0, 0)
	if err != nil {
		panic(err)
	}

	// If the number of elements in this array is 0, then return an
	// empty slice as there is nothing to expand. Trying to expand it
	// anyway could lead to crashes as any child maps, arrays or sets
	// that no longer exist are still shown as empty with a count of 0.
	if num == 0 {
		return []interface{}{}
	}

	// NOTE: "num" is not necessarily accurate, e.g. if a user tampers
	// with state, so the following code should not crash when given a
	// number of items more or less than what's given in num. The
	// num key is mainly just a hint that this is a list or set.

	// The Schema "Set" type stores its values in an array format, but
	// using numeric hash values instead of ordinal keys. Take the set
	// of keys regardless of value, and expand them in numeric order.
	// See GH-11042 for more details.
	keySet := map[int]bool{}
	for k := range m {
		if !strings.HasPrefix(k, prefix+".") {
			continue
		}

		key := k[len(prefix)+1:]
		idx := strings.Index(key, ".")
		if idx != -1 {
			key = key[:idx]
		}

		// skip the count value
		if key == "#" {
			continue
		}

		noBrackets := strings.Replace(strings.Replace(key, "[", "", 1), "]", "", 1)
		k, err := strconv.Atoi(noBrackets)
		if err != nil {
			panic(err)
		}
		keySet[k] = true
	}

	keysList := make([]int, 0, num)
	for key := range keySet {
		keysList = append(keysList, key)
	}
	sort.Ints(keysList)

	result := make([]interface{}, len(keysList))
	for i, key := range keysList {
		keyString := strconv.Itoa(key)
		result[i] = Expand(m, fmt.Sprintf("%s.[%s]", prefix, keyString))
	}

	deletePrefix(m, prefix)
	return result
}

func expandMap(m map[string]string, prefix string) map[string]interface{} {
	result := make(map[string]interface{})
	for k := range m {
		if !strings.HasPrefix(k, prefix) {
			continue
		}

		key := k[len(prefix):]
		idx := strings.Index(key, ".")
		if idx != -1 {
			key = key[:idx]
		}
		if _, ok := result[key]; ok {
			continue
		}

		result[key] = Expand(m, k[:len(prefix)+len(key)])
		delete(m, k)
	}

	return result
}
