package smapper

import (
	"reflect"
)

func SameTypeConverter(value *reflect.Value, _, _ reflect.Type) (*reflect.Value, error) {
	return value, nil
}

func UnionSlices[T any](slices ...[]T) []T {
	var union = make([]T, 0)
	for _, slice := range slices {
		union = append(union, slice...)
	}
	return union
}

func GetReflectType[T any]() reflect.Type {
	var item *T
	return reflect.TypeOf(item).Elem()
}

func NewConverterDesc[SRC comparable, DST comparable](
	converter ValueConverter,
) ValueConverterDesc {
	return ValueConverterDesc{
		Src:       GetReflectType[SRC](),
		Dst:       GetReflectType[DST](),
		Converter: converter,
	}
}

var DefaultConverters = UnionSlices(
	buildIntConverters[int](),
	buildIntConverters[int8](),
	buildIntConverters[int16](),
	buildIntConverters[int32](),
	buildIntConverters[int64](),
	buildIntConverters[uint](),
	buildIntConverters[uint8](),
	buildIntConverters[uint16](),
	buildIntConverters[uint32](),
	buildIntConverters[uint64](),
	buildIntConverters[uintptr](),
	StringConverters,
	ComplexConverters,
	FloatConverters,
	BoolConverters,
)
