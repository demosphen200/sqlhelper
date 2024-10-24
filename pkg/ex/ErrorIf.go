package ex

import "errors"

func ErrorIf(err *error, setError bool, message string) bool {
	if setError {
		*err = errors.New(message)
		return false
	} else {
		return true
	}
}
