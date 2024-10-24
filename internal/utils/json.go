package utils

import "encoding/json"

func JsonString(value any) string {
	return string(Must(json.Marshal(value)))
}
