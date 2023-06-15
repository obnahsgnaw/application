package utils

import (
	"bytes"
	"encoding/json"
)

func Struct2Map(s interface{}) (m map[string]interface{}) {
	if marshalContent, err := json.Marshal(s); err != nil {
		return nil
	} else {
		d := json.NewDecoder(bytes.NewReader(marshalContent))
		d.UseNumber()
		if err = d.Decode(&m); err != nil {
			return nil
		} else {
			return
		}
	}
}
