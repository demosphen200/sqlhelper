package smapper

import (
	"errors"
	"reflect"
	"strconv"
)

func buildIntConverters[
	T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr,
]() []ValueConverterDesc {
	return []ValueConverterDesc{
		NewConverterDesc[T, int](anyIntConverter),
		NewConverterDesc[T, int8](anyIntConverter),
		NewConverterDesc[T, int16](anyIntConverter),
		NewConverterDesc[T, int32](anyIntConverter),
		NewConverterDesc[T, int64](anyIntConverter),
		NewConverterDesc[T, uint](anyIntConverter),
		NewConverterDesc[T, uint8](anyIntConverter),
		NewConverterDesc[T, uint16](anyIntConverter),
		NewConverterDesc[T, uint32](anyIntConverter),
		NewConverterDesc[T, uint64](anyIntConverter),
		NewConverterDesc[T, uintptr](anyIntConverter),

		NewConverterDesc[T, float32](anyIntToFloat32Converter),
		NewConverterDesc[T, float64](anyIntToFloat32Converter),

		NewConverterDesc[T, complex64](anyIntToComplexConverter),
		NewConverterDesc[T, complex128](anyIntToComplexConverter),

		NewConverterDesc[T, string](anyIntToStringConverter),
		NewConverterDesc[T, bool](anyIntToBoolConverter),
	}
}

/*
	func buildIntConverters(intKind reflect.Kind) []ValueConverterDesc {
		return []ValueConverterDesc{
			NewConverterDesc(),
			{
				Src:       TypeSignature{intKind},
				Dst:       TypeSignature{reflect.Int},
				Converter: anyIntConverter,
			},
			{
				Src:       TypeSignature{intKind},
				Dst:       TypeSignature{reflect.Int8},
				Converter: anyIntConverter,
			},
			{
				Src:       TypeSignature{intKind},
				Dst:       TypeSignature{reflect.Int16},
				Converter: anyIntConverter,
			},
			{
				Src:       TypeSignature{intKind},
				Dst:       TypeSignature{reflect.Int32},
				Converter: anyIntConverter,
			},
			{
				Src:       TypeSignature{intKind},
				Dst:       TypeSignature{reflect.Int64},
				Converter: anyIntConverter,
			},

			{
				Src:       TypeSignature{intKind},
				Dst:       TypeSignature{reflect.Uint},
				Converter: anyIntConverter,
			},
			{
				Src:       TypeSignature{intKind},
				Dst:       TypeSignature{reflect.Uint8},
				Converter: anyIntConverter,
			},
			{
				Src:       TypeSignature{intKind},
				Dst:       TypeSignature{reflect.Uint16},
				Converter: anyIntConverter,
			},
			{
				Src:       TypeSignature{intKind},
				Dst:       TypeSignature{reflect.Uint32},
				Converter: anyIntConverter,
			},
			{
				Src:       TypeSignature{intKind},
				Dst:       TypeSignature{reflect.Uint64},
				Converter: anyIntConverter,
			},
			{
				Src:       TypeSignature{intKind},
				Dst:       TypeSignature{reflect.Uintptr},
				Converter: anyIntConverter,
			},
			{
				Src:       TypeSignature{intKind},
				Dst:       TypeSignature{reflect.Float32},
				Converter: anyIntToFloat32Converter,
			},
			{
				Src:       TypeSignature{intKind},
				Dst:       TypeSignature{reflect.Float64},
				Converter: anyIntToFloat32Converter,
			},
			{
				Src:       TypeSignature{intKind},
				Dst:       TypeSignature{reflect.Complex64},
				Converter: anyIntToComplexConverter,
			},
			{
				Src:       TypeSignature{intKind},
				Dst:       TypeSignature{reflect.Complex128},
				Converter: anyIntToComplexConverter,
			},
			{
				Src:       TypeSignature{intKind},
				Dst:       TypeSignature{reflect.String},
				Converter: anyIntToStringConverter,
			},
			{
				Src:       TypeSignature{intKind},
				Dst:       TypeSignature{reflect.Bool},
				Converter: anyIntToBoolConverter,
			},
		}
	}
*/
func anyIntConverter(value *reflect.Value, _, dst reflect.Type) (*reflect.Value, error) {
	return int64ToTypedInt(getInt64Value(value), dst.Kind()), nil
}

func anyIntToFloat32Converter(value *reflect.Value, _, dst reflect.Type) (*reflect.Value, error) {
	int64Value := getInt64Value(value)
	var result reflect.Value
	switch dst.Kind() {
	case reflect.Float32:
		{
			typed := float32(int64Value)
			result = reflect.ValueOf(typed)
		}
	case reflect.Float64:
		{
			typed := float64(int64Value)
			result = reflect.ValueOf(typed)
		}
	default:
		return nil, errors.New("target type is not float32 or float64")
	}
	return &result, nil
}

func anyIntToStringConverter(value *reflect.Value, _, _ reflect.Type) (*reflect.Value, error) {
	typed := strconv.Itoa(int(getInt64Value(value)))
	result := reflect.ValueOf(typed)
	return &result, nil
}

func anyIntToBoolConverter(value *reflect.Value, _, _ reflect.Type) (*reflect.Value, error) {
	typed := getInt64Value(value) != 0
	result := reflect.ValueOf(typed)
	return &result, nil
}

func anyIntToComplexConverter(value *reflect.Value, _, dst reflect.Type) (*reflect.Value, error) {
	int64value := getInt64Value(value)
	var result reflect.Value
	switch dst.Kind() {
	case reflect.Complex64:
		{
			var typed complex64 = complex(float32(int64value), 0.0)
			result = reflect.ValueOf(typed)
		}
	case reflect.Complex128:
		{
			var typed complex128 = complex(float64(int64value), 0.0)
			result = reflect.ValueOf(typed)
		}
	default:
		return nil, errors.New("target type is not complex64 or complex128")
	}
	return &result, nil
}

func int64ToTypedInt(value int64, targetKind reflect.Kind) *reflect.Value {
	var result reflect.Value
	switch targetKind {
	case reflect.Int:
		{
			typedInt := int(value)
			result = reflect.ValueOf(typedInt)
		}
	case reflect.Int8:
		{
			typedInt := int8(value)
			result = reflect.ValueOf(typedInt)
		}
	case reflect.Int16:
		{
			typedInt := int16(value)
			result = reflect.ValueOf(typedInt)
		}
	case reflect.Int32:
		{
			typedInt := int32(value)
			result = reflect.ValueOf(typedInt)
		}
	case reflect.Int64:
		{
			typedInt := value
			result = reflect.ValueOf(typedInt)
		}

	case reflect.Uint:
		{
			typedInt := uint(value)
			result = reflect.ValueOf(typedInt)
		}
	case reflect.Uint8:
		{
			typedInt := uint8(value)
			result = reflect.ValueOf(typedInt)
		}
	case reflect.Uint16:
		{
			typedInt := uint16(value)
			result = reflect.ValueOf(typedInt)
		}
	case reflect.Uint32:
		{
			typedInt := uint32(value)
			result = reflect.ValueOf(typedInt)
		}
	case reflect.Uint64:
		{
			typedInt := uint64(value)
			result = reflect.ValueOf(typedInt)
		}
	case reflect.Uintptr:
		{
			typedInt := uintptr(value)
			result = reflect.ValueOf(typedInt)
		}

	}

	return &result
}

func getInt64Value(value *reflect.Value) int64 {
	var result int64
	switch value.Kind() {
	case reflect.Int:
		result = value.Int()
	case reflect.Int8:
		result = value.Int()
	case reflect.Int16:
		result = value.Int()
	case reflect.Int32:
		result = value.Int()
	case reflect.Int64:
		result = value.Int()
	case reflect.Uint:
		result = int64(value.Uint())
	case reflect.Uint8:
		result = int64(value.Uint())
	case reflect.Uint16:
		result = int64(value.Uint())
	case reflect.Uint32:
		result = int64(value.Uint())
	case reflect.Uint64:
		result = int64(value.Uint())
	case reflect.Uintptr:
		result = int64(value.Uint())
	}
	return result
}
