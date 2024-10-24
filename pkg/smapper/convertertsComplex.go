package smapper

import (
	"reflect"
	"strconv"
)

var ComplexConverters = []ValueConverterDesc{
	NewConverterDesc[complex64, complex64](SameTypeConverter),
	NewConverterDesc[complex128, complex128](SameTypeConverter),

	NewConverterDesc[complex64, string](CreateComplexToStringConverter('f', 5, 64)),
	NewConverterDesc[complex128, string](CreateComplexToStringConverter('f', 5, 128)),

	/*
		{
			Src:       TypeSignature{reflect.Complex64},
			Dst:       TypeSignature{reflect.Int},
			Converter: CreateComplexToStringConverter('f', 5, 64),
		},*/

}

func CreateComplexToStringConverter(fmt byte, prec int, bitSize int) ValueConverter {
	return func(value *reflect.Value, _, _ reflect.Type) (*reflect.Value, error) {
		typed := strconv.FormatComplex(value.Complex(), fmt, prec, bitSize)
		result := reflect.ValueOf(typed)
		return &result, nil
	}
}
