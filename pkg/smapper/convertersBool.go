package smapper

import (
	"reflect"
)

var BoolConverters = []ValueConverterDesc{
	NewConverterDesc[bool, bool](SameTypeConverter),
	NewConverterDesc[bool, int](boolToIntConverter),
	NewConverterDesc[bool, int8](boolToIntConverter),
	NewConverterDesc[bool, int16](boolToIntConverter),
	NewConverterDesc[bool, int32](boolToIntConverter),
	NewConverterDesc[bool, int64](boolToIntConverter),
	NewConverterDesc[bool, uint](boolToIntConverter),
	NewConverterDesc[bool, uint8](boolToIntConverter),
	NewConverterDesc[bool, uint16](boolToIntConverter),
	NewConverterDesc[bool, uint32](boolToIntConverter),
	NewConverterDesc[bool, uint64](boolToIntConverter),
	NewConverterDesc[bool, uintptr](boolToIntConverter),
	NewConverterDesc[bool, string](boolToStringConverter),
}

func boolToIntConverter(value *reflect.Value, _, dst reflect.Type) (*reflect.Value, error) {
	var resultValue int64 = 0
	if value.Bool() {
		resultValue = 1
	}
	return int64ToTypedInt(resultValue, dst.Kind()), nil
}

func boolToStringConverter(value *reflect.Value, _, _ reflect.Type) (*reflect.Value, error) {
	var str = "false"
	if value.Bool() {
		str = "true"
	}
	result := reflect.ValueOf(str)
	return &result, nil
}
