package smapper

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestGetTypeSignature(t *testing.T) {
	var test = struct {
		i              int
		f              float32
		slicePtrString []*string
		arrString      [1]string
	}{}

	assert.Equal(t, true, reflect.DeepEqual(
		GetTypeSignature(reflect.TypeOf(test)),
		TypeSignature{reflect.Struct},
	))

	assert.Equal(t, true, reflect.DeepEqual(
		GetTypeSignature(reflect.TypeOf(test.i)),
		TypeSignature{reflect.Int},
	))

	assert.Equal(t, true, reflect.DeepEqual(
		GetTypeSignature(reflect.TypeOf(test.f)),
		TypeSignature{reflect.Float32},
	))

	assert.Equal(t, true, reflect.DeepEqual(
		GetTypeSignature(reflect.TypeOf(test.slicePtrString)),
		TypeSignature{reflect.Slice, reflect.Pointer, reflect.String},
	))

	assert.Equal(t, true, reflect.DeepEqual(
		GetTypeSignature(reflect.TypeOf(test.arrString)),
		TypeSignature{reflect.Array, reflect.String},
	))
}
