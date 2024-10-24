package smapper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapper_GetKeys(t *testing.T) {
	a := struct {
		Ignored string
		Field1  []int `map:"key1"`
		Field2  []int `map:"key2"`
		Field22 []int `map:"key2"`
	}{}
	mapper := MakeMapper()
	assert.Equal(t, []string{"key1", "key2", "key2"}, mapper.GetKeys(&a))
}

func TestMapper_GetTaggedFieldNames(t *testing.T) {
	a := struct {
		Ignored string
		Field1  []int `map:"key"`
		Field2  []int `map:"key2"`
	}{}
	mapper := MakeMapper()
	assert.Equal(t, []string{"Field1", "Field2"}, mapper.GetTaggedFieldNames(&a))
}
