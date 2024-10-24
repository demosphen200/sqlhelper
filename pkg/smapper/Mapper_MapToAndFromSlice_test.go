package smapper

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMapper_MapToSlice_ShouldMap(t *testing.T) {
	a := struct {
		Field  int    `map:"key"`
		Field2 string `map:"key2"`
	}{
		Field:  11,
		Field2: "str",
	}
	mapper := MakeMapper()
	assert.Equal(t, []any{11, "str"}, mapper.MapToSlice(&a))
}

type mapFromSliceTestStruct struct {
	Field  int    `map:"key"`
	Field2 string `map:"key2"`
}

func mapFromSliceTest(t *testing.T, mapper *Mapper) {
	var a = mapFromSliceTestStruct{}
	if t != nil {
		assert.NoError(t, mapper.MapFromSlice("src00", []any{"11", "str"}, &a))
		assert.Equal(t, 11, a.Field)
		assert.Equal(t, "str", a.Field2)
	} else {
		if mapper.MapFromSlice("src00", []any{"11", "str"}, &a) != nil {
			panic("map error")
		}
		if a.Field != 11 || a.Field2 != "str" {
			panic("map error")
		}
	}
}

func TestMapper_MapFromSlice_ShouldMap(t *testing.T) {
	mapper := MakeMapper()
	mapFromSliceTest(t, &mapper)
}

func TestMapper_MapFromSlice_ShouldAllowLazyConverterInit(t *testing.T) {
	str := "2"
	a := struct {
		Field *string `map:"key"`
	}{
		Field: &str,
	}
	mapper := MakeMapper()
	assert.NoError(t, mapper.MapFromSlice("src00", []any{nil}, &a))
	assert.Nil(t, a.Field)

	assert.NoError(t, mapper.MapFromSlice("src00", []any{nil}, &a))
	assert.Nil(t, a.Field)

	assert.NoError(t, mapper.MapFromSlice("src00", []any{2}, &a))
	assert.Equal(t, &str, a.Field)

	assert.NoError(t, mapper.MapFromSlice("src00", []any{nil}, &a))
	assert.Nil(t, a.Field)
}

func mapFromSliceWorker(done chan bool, mapper *Mapper) {
	var fake = 1
	for n := 0; n < 500; n++ {
		if mapper == nil {
			newMapper := MakeMapper()
			mapFromSliceTest(nil, &newMapper)
		} else {
			newMapper := NewMapper()
			fake++
			if fake < 0 {
				mapper = newMapper
			}
			mapFromSliceTest(nil, mapper)
		}
	}
	done <- true
}

func TestMapper_MapFromSlice_ShouldWorkInMultiThreadsWithPlan(t *testing.T) {
	threadCount := 30

	mapper := MakeMapper()
	var done = make(chan bool)
	for n := 0; n < threadCount; n++ {
		go mapFromSliceWorker(done, &mapper)
	}
	for n := 0; n < threadCount; n++ {
		<-done
	}
}

func TestMapper_MapFromSlice_ShouldWorkInMultiThreadsWithoutPlan(t *testing.T) {
	threadCount := 30

	var done = make(chan bool)
	for n := 0; n < threadCount; n++ {
		go mapFromSliceWorker(done, nil)
	}
	for n := 0; n < threadCount; n++ {
		<-done
	}
}

func TestMapper_MapFromSlice_CanReturnConverterNotFoundError(t *testing.T) {
	a := struct {
		Field int `map:"key"`
	}{}
	mapper := MakeMapper()
	mapper.Converters = make([]ValueConverterDesc, 0)
	err := mapper.MapFromSlice("id", []any{"str"}, &a)
	if assert.Error(t, err) {
		assert.Equal(t, NewConverterNotFoundError(GetReflectType[string](), GetReflectType[int]()), err)
		//assert.Equal(t, NewConverterNotFoundError([]reflect.Kind{reflect.String}, []reflect.Kind{reflect.Int}), err)
	}
}

func TestMapper_MapFromSlice_CanReturnConverterNotFoundErrorOnLazyInit(t *testing.T) {
	a := struct {
		Field *int `map:"key"`
	}{}
	mapper := MakeMapper()
	mapper.Converters = make([]ValueConverterDesc, 0)
	assert.NoError(t, mapper.MapFromSlice("id", []any{nil}, &a))
	err := mapper.MapFromSlice("id", []any{"1"}, &a)
	if assert.Error(t, err) {
		assert.Equal(t, NewConverterNotFoundError(GetReflectType[string](), GetReflectType[*int]()), err)
		//assert.Equal(t, NewConverterNotFoundError([]reflect.Kind{reflect.String}, []reflect.Kind{reflect.Pointer, reflect.Int}), err)
	}
}

func TestMapper_MapFromSlice_CanReturnDestinationNotNullableError(t *testing.T) {
	a := struct {
		Field int `map:"key"`
	}{}
	mapper := MakeMapper()
	mapper.Converters = make([]ValueConverterDesc, 0)
	err := mapper.MapFromSlice("id", []any{nil}, &a)
	if assert.Error(t, err) {
		assert.Equal(t, NewDestinationNotNullableError(0, GetReflectType[int]() /*TypeSignature{reflect.Int}*/), err)
	}
}

func returningStringConverter(_ *reflect.Value, _, _ reflect.Type) (*reflect.Value, error) {
	str := "str"
	result := reflect.ValueOf(str)
	return &result, nil
}

func returningErrorConverter(_ *reflect.Value, _, _ reflect.Type) (*reflect.Value, error) {
	return nil, errors.New("some error")
}

func TestMapper_MapFromSlice_CanReturnNotAssignableValueError(t *testing.T) {
	a := struct {
		Field int `map:"key"`
	}{}
	mapper := MakeMapper()
	mapper.Converters = []ValueConverterDesc{
		NewConverterDesc[string, int](returningStringConverter),
		/*
			{
				Src:       TypeSignature{reflect.String},
				Dst:       TypeSignature{reflect.Int},
				Converter: returningStringConverter,
			},*/
	}
	err := mapper.MapFromSlice("id", []any{"1"}, &a)
	if assert.Error(t, err) {
		assert.Equal(
			t,
			NewNotAssignableValueError(
				"0",
				"Field",
				"1",
				GetReflectType[string](), //TypeSignature{reflect.String},
				"str",
				GetReflectType[string](), //TypeSignature{reflect.String},
				GetReflectType[int](),    //TypeSignature{reflect.Int},
			),
			err,
		)
	}
}

func TestMapper_MapFromSlice_CanReturnConverterError(t *testing.T) {
	a := struct {
		Field int `map:"key"`
	}{}
	mapper := MakeMapper()
	mapper.Converters = []ValueConverterDesc{
		NewConverterDesc[string, int](returningErrorConverter),
		/*
			{
				Src:       TypeSignature{reflect.String},
				Dst:       TypeSignature{reflect.Int},
				Converter: returningErrorConverter,
			},*/
	}
	err := mapper.MapFromSlice("id", []any{"1"}, &a)
	if assert.Error(t, err) {
		assert.Equal(
			t,
			NewConverterError(
				"0",
				"Field",
				"1",
				GetReflectType[string](), //TypeSignature{reflect.String},
				GetReflectType[int](),    //TypeSignature{reflect.Int},
				errors.New("some error"),
			),
			err,
		)
		//fmt.Println(err.Error())
	}
}
