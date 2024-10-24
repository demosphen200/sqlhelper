package smapper

import (
	"reflect"
	"strconv"
)

var StringConverters = []ValueConverterDesc{
	NewConverterDesc[string, string](SameTypeConverter),
	NewConverterDesc[string, int](stringToIntConverter),
	NewConverterDesc[string, int8](stringToIntConverter),
	NewConverterDesc[string, int16](stringToIntConverter),
	NewConverterDesc[string, int32](stringToIntConverter),
	NewConverterDesc[string, int64](stringToIntConverter),
	NewConverterDesc[string, uint](stringToIntConverter),
	NewConverterDesc[string, uint8](stringToIntConverter),
	NewConverterDesc[string, uint16](stringToIntConverter),
	NewConverterDesc[string, uint32](stringToIntConverter),
	NewConverterDesc[string, uint64](stringToIntConverter),
	NewConverterDesc[string, uintptr](stringToIntConverter),

	NewConverterDesc[string, float32](stringToFloat32Converter),
	NewConverterDesc[string, float64](stringToFloat64Converter),

	NewConverterDesc[string, complex64](stringToComplex64Converter),
	NewConverterDesc[string, complex128](stringToComplex128Converter),

	NewConverterDesc[string, bool](stringToBoolConverter),
}

func stringToFloat64Converter(value *reflect.Value, _, _ reflect.Type) (*reflect.Value, error) {
	var err error
	typed64, err := strconv.ParseFloat(value.String(), 64)
	if err != nil {
		return nil, err
	}
	result := reflect.ValueOf(typed64)
	return &result, nil
}

func stringToFloat32Converter(value *reflect.Value, _, _ reflect.Type) (*reflect.Value, error) {
	var err error
	typed64, err := strconv.ParseFloat(value.String(), 32)
	if err != nil {
		return nil, err
	}
	typed32 := float32(typed64)
	result := reflect.ValueOf(typed32)
	return &result, nil
}

func stringToComplex128Converter(value *reflect.Value, _, _ reflect.Type) (*reflect.Value, error) {
	var err error
	typed128, err := strconv.ParseComplex(value.String(), 128)
	if err != nil {
		return nil, err
	}
	result := reflect.ValueOf(typed128)
	return &result, nil
}

func stringToComplex64Converter(value *reflect.Value, _, _ reflect.Type) (*reflect.Value, error) {
	var err error
	typed128, err := strconv.ParseComplex(value.String(), 64)
	if err != nil {
		return nil, err
	}
	typed64 := complex64(typed128)
	result := reflect.ValueOf(typed64)
	return &result, nil
}

func stringToIntConverter(value *reflect.Value, _, dst reflect.Type) (*reflect.Value, error) {
	var int64Value, err = strconv.ParseInt(value.String(), 10, 64)
	if err != nil {
		return nil, err
	}
	return int64ToTypedInt(int64Value, dst.Kind()), nil
}

func stringToBoolConverter(value *reflect.Value, _, _ reflect.Type) (*reflect.Value, error) {
	boolValue := value.String() == "true"
	result := reflect.ValueOf(boolValue)
	return &result, nil
}
