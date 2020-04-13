package siftr

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Policy is defined structure that drives the JSON filtering in Siftr.
type Policy struct {
	Whitelist map[string]string
	Sibling   map[string][]string
}

// Main function of the package that will process and filter the given data
// based on the given policy.
//
// Returns a map[string]interface{} if all successful.
func Sift(data []byte, policy *Policy) (map[string]interface{}, error) {
	if policy == nil {
		return nil, fmt.Errorf("error no policy given")
	}
	if data == nil {
		return nil, fmt.Errorf("error no data give to sift")
	}

	var mapData map[string]interface{}
	err := json.Unmarshal(data, &mapData)
	if err != nil {
		return nil, fmt.Errorf("error reading data as a map")
	}

	flattened := Flatten(mapData)

	if policy.Sibling != nil {
		flattened = sibling(flattened, policy.Sibling)
	}

	if policy.Whitelist != nil {
		flattened = whitelist(flattened, policy.Whitelist)
	}

	return sortAndExpand(flattened), nil

}

// A quicksort of the keys actually helps the performance a little.
// Every little bit counts when doing arbitrary json filtering.
func sortAndExpand(data Map) map[string]interface{} {
	keyList := make([]string, len(data))
	for v := range data {
		keyList = append(keyList, v)
	}

	newMap := make(Map, len(data))
	sort.Strings(keyList)
	for _, k := range keyList {
		newMap[k] = data[k]
	}
	keys := newMap.Keys()
	completed := make(map[string]interface{}, len(keys))
	for _, key := range keys {
		if strings.TrimSpace(key) != "" {
			completed[key] = Expand(newMap, key)
		}
	}

	return completed
}

// Check the whitelist and remove any fields from the data that aren't in it.
// This will also check if the value of the given key should be masked, and
// if so, then go ahead and mask it.
func whitelist(data Map, policy map[string]string) Map {
	re := regexp.MustCompile("\\d+")
	for k, v := range data {
		tmp := re.ReplaceAllString(k, "*")
		if action, ok := policy[tmp]; !ok && !strings.Contains(k, "#") {
			delete(data, k)
		} else if action == "mask" {
			data[k] = mask(v)
		}
	}

	return data
}

// Check the sibling field and values to know if an entire object should be
// removed from the data. If the field does not have a value matching one
// from the list of allowable values, then the whole object should be removed.
func sibling(data Map, policy map[string][]string) Map {
	re := regexp.MustCompile("\\d+")
	for typeField, types := range policy {
		for k, v := range data {
			tmp := re.ReplaceAllString(k, "*")
			if typeField == tmp && !contains(types, v) {
				split := strings.Split(k, ".")
				prefix := strings.Join(split[:len(split)-1], ".")
				data.Delete(prefix)
			}
		}
	}

	return data
}

// Just mask the characters of the string except for the last 4 characters.
func mask(s string) string {
	rs := []rune(s)
	for i := 0; i < len(rs)-4; i++ {
		rs[i] = '*'
	}
	return string(rs)
}

// Check if the string slice contains a certain string.
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
