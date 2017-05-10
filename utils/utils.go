package utils

import "encoding/json"

// PrettyPrint pretty prints an object
func PrettyPrint(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
