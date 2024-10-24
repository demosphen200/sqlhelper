package smapper

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_CreateTypedNil(t *testing.T) {
	var s = struct{}{}
	var i = 11
	var ptrI = &i
	var ptrS = &s

	setI := func() {
		typedNil := CreateTypedNil(reflect.TypeOf(ptrI))
		reflect.ValueOf(&ptrI).Elem().Set(*typedNil)
	}
	setS := func() {
		typedNil := CreateTypedNil(reflect.TypeOf(ptrS))
		reflect.ValueOf(&ptrS).Elem().Set(*typedNil)
	}

	assert.NotPanics(t, setI)
	assert.Nil(t, ptrI)
	assert.Equal(t, 11, i)

	assert.NotPanics(t, setS)
	assert.Nil(t, ptrS)
	assert.Equal(t, s, struct{}{})

}

func Test_IsPtrToStruct(t *testing.T) {
	var str = struct{}{}
	var iface interface{} = &str
	var an any = &str

	var str2 *struct{}
	var i int
	var s *string

	assert.Equal(t, true, IsPtrToStruct(&str))
	assert.Equal(t, true, IsPtrToStruct(iface))
	assert.Equal(t, true, IsPtrToStruct(an))
	assert.Equal(t, false, IsPtrToStruct(nil))
	assert.Equal(t, false, IsPtrToStruct(&str2))
	assert.Equal(t, false, IsPtrToStruct(str))
	assert.Equal(t, false, IsPtrToStruct(i))
	assert.Equal(t, false, IsPtrToStruct(s))
}
