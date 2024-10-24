package smapper

import (
	"reflect"
	"strconv"
)

var FloatConverters = []ValueConverterDesc{
	NewConverterDesc[float32, float32](SameTypeConverter),
	NewConverterDesc[float64, float64](SameTypeConverter),
	NewConverterDesc[float32, float64](floatToFloatConverter),
	NewConverterDesc[float64, float32](floatToFloatConverter),
	NewConverterDesc[float32, string](CreateFloatToStringConverter('f', 5, 32)),
	NewConverterDesc[float64, string](CreateFloatToStringConverter('f', 5, 64)),
}

func CreateFloatToStringConverter(fmt byte, prec int, bitSize int) ValueConverter {
	return func(value *reflect.Value, _, _ reflect.Type) (*reflect.Value, error) {
		typed := strconv.FormatFloat(value.Float(), fmt, prec, bitSize)
		result := reflect.ValueOf(typed)
		return &result, nil
	}
}

func floatToFloatConverter(value *reflect.Value, _, dst reflect.Type) (*reflect.Value, error) {
	float64Value := value.Float()
	var result reflect.Value
	switch dst.Kind() {
	case reflect.Float32:
		{
			typed := float32(float64Value)
			result = reflect.ValueOf(typed)
		}
	case reflect.Float64:
		{
			typed := float64Value
			result = reflect.ValueOf(typed)

		}
	}
	return &result, nil
}
