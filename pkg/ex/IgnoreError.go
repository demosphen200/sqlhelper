package ex

import "log"

func IgnoreError(err error) {
	log.Printf("ignored error %s", err.Error())
}
