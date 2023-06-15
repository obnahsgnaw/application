package utils

import (
	"encoding/json"
)

func ToJson(v interface{}) string {
	if v1, ok := v.(string); ok {
		return v1
	}
	b, _ := json.Marshal(v)
	return string(b)
}

func ParseJson(d []byte, v interface{}) bool {
	err := json.Unmarshal(d, &v)

	return err == nil
}
