package testhelper

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var currentT *testing.T

func HelperSetT(t *testing.T) {
	currentT = t
}

func NoError0(err error) {
	if !assert.NoError(currentT, err) {
		panic(err.Error())
	}
}

func NoError[T any](value T, err error) T {
	if !assert.NoError(currentT, err) {
		panic(err.Error())
	}
	return value
}

func Success0(err error) {
	if !assert.NoError(currentT, err) {
		panic(err.Error())
	}
}

func Success[T any](value T, err error) T {
	if !assert.NoError(currentT, err) {
		panic(err.Error())
	}
	return value
}

func NoError2[T1 any, T2 any](value1 T1, value2 T2, err error) (T1, T2) {
	if !assert.NoError(currentT, err) {
		panic(err.Error())
	}
	return value1, value2
}

func Success2[T1 any, T2 any](value1 T1, value2 T2, err error) (T1, T2) {
	if !assert.NoError(currentT, err) {
		//panic(err.Error())
		currentT.Fatal(err.Error())
	}
	return value1, value2
}

func Run(t *testing.T, fn func()) {
	currentT = t
	fn()
}

func CalcDuration(fn func()) time.Duration {
	start := time.Now()
	fn()
	return time.Since(start)
}
