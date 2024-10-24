package ex

import "strconv"

func Default[T comparable](value T, def T) T {
	var initValue T
	if value == initValue {
		return def
	} else {
		return value
	}
}

func DefaultIntFromString(str string, def int) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		return def
	} else {
		return value
	}
}
