package util

import "encoding/json"

func Json(value interface{}) string {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return ""
	}
	return BytesToString(data)
}
