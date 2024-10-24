package ex

import "fmt"

func ErrorMessage(message string, err error) string {
	return fmt.Sprintf("%s, error:%s", message, err.Error())
}
