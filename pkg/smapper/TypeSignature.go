package smapper

import "reflect"

type TypeSignature []reflect.Kind

func (signature TypeSignature) Equal(other TypeSignature) bool {
	if len(signature) != len(other) {
		return false
	}
	for index, kind := range signature {
		if other[index] != kind {
			return false
		}
	}
	return true
}

func GetTypeSignature(fieldType reflect.Type) TypeSignature {
	signature := make(TypeSignature, 1)
	kind := fieldType.Kind()
	signature[0] = kind
	switch kind {
	case reflect.Array, reflect.Chan, reflect.Pointer, reflect.Slice:
		signature = append(signature, GetTypeSignature(fieldType.Elem())...)
	}
	return signature
}
