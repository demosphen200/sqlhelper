package ex

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestMapStringToInt(t *testing.T) {
	src := []string{"123", "456"}
	dst := Map(src, func(stringValue string) int {
		intValue, _ := strconv.Atoi(stringValue)
		return intValue
	})

	assert.Equal(t, len(src), len(dst))
	assert.Equal(t, 123, dst[0])
	assert.Equal(t, 456, dst[1])
}

func TestMapIntToString(t *testing.T) {
	src := []int{123, 456}
	dst := Map(src, func(intValue int) string {
		stringValue := strconv.Itoa(intValue)
		return stringValue
	})

	assert.Equal(t, len(src), len(dst))
	assert.Equal(t, "123", dst[0])
	assert.Equal(t, "456", dst[1])
}

func TestMapOnEmptySlice(t *testing.T) {
	src := make([]int, 0)
	dst := Map(src, func(intValue int) string {
		stringValue := strconv.Itoa(intValue)
		return stringValue
	})

	assert.Equal(t, 0, len(dst))
}

func TestFiltered(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6, 7}
	dst := Filtered(src, func(value int) bool {
		return value > 4
	})

	assert.Equal(t, 7, len(src))
	assert.Equal(t, 1, src[0])
	assert.Equal(t, 2, src[1])
	assert.Equal(t, 3, src[2])
	assert.Equal(t, 4, src[3])
	assert.Equal(t, 5, src[4])
	assert.Equal(t, 6, src[5])
	assert.Equal(t, 7, src[6])

	assert.Equal(t, 3, len(dst))
	assert.Equal(t, 5, dst[0])
	assert.Equal(t, 6, dst[1])
	assert.Equal(t, 7, dst[2])
}

func TestFilteredOnEmptySlice(t *testing.T) {
	src := make([]int, 0)
	dst := Filtered(src, func(value int) bool {
		return value > 4
	})
	assert.Equal(t, 0, len(dst))
}

func TestFilteredOnEmptyResult(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6, 7}
	dst := Filtered(src, func(value int) bool {
		return value > 44
	})
	assert.Equal(t, 7, len(src))
	assert.Equal(t, 0, len(dst))
}
