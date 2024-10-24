package smapper

import (
	"reflect"
)

type ValueConverter = func(value *reflect.Value, src, dst reflect.Type) (*reflect.Value, error)

type ValueConverterDesc struct {
	Src       reflect.Type
	Dst       reflect.Type
	Converter ValueConverter
}
