package utils

import "encoding/json"

func InSlice[T comparable](haystack []T, needle T) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

func PrettyPrint(data interface{}) string {
	result, _ := json.MarshalIndent(data, "", "\t")
	return string(result)
}
