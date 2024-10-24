package smapper

import "reflect"

func IsNullable(kind reflect.Kind) bool {
	return kind == reflect.Pointer || kind == reflect.Struct || kind == reflect.Slice
}

func CreateTypedNil(ptrType reflect.Type) *reflect.Value {
	ptrToTypedNil := reflect.New(ptrType)
	typedNil := ptrToTypedNil.Elem()
	return &typedNil
}

func IsPtrToStruct(value interface{}) bool {
	if value == nil {
		return false
	}
	if reflect.TypeOf(value).Kind() != reflect.Pointer {
		return false
	}
	return reflect.TypeOf(value).Elem().Kind() == reflect.Struct
}

func IsPtrToSliceOfStruct(value interface{}) bool {
	if value == nil {
		return false
	}
	if reflect.TypeOf(value).Kind() != reflect.Pointer {
		return false
	}
	if reflect.TypeOf(value).Elem().Kind() != reflect.Slice {
		return false
	}
	return reflect.TypeOf(value).Elem().Elem().Kind() == reflect.Struct
}

func EnumTaggedStructFields(
	structType reflect.Type,
	tag string,
	fn func(index int, structField *reflect.StructField, tagValue string) error,
) error {
	if structType == nil || structType.Kind() != reflect.Struct {
		panic("structPtr must be pointer to struct")
	}
	for index := 0; index < structType.NumField(); index++ {
		structField := structType.Field(index)
		tagValue := structField.Tag.Get(tag)
		if tagValue == "" {
			continue
		}
		if err := fn(index, &structField, tagValue); err != nil {
			return err
		}
	}
	return nil
}
