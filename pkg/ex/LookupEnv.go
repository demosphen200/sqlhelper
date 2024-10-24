package ex

import "os"

func LookupEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	} else {
		return ""
	}
}
