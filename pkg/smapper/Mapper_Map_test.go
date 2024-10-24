package smapper

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strconv"
	"testing"
)

func TestMapper_Map_ShouldMapBaseTypes(t *testing.T) {
	mapBaseTypesTest(t, nil)
}

func mapBaseTypesTest(t *testing.T, mapper *Mapper) {
	mapIntTest[int](t, 123, 1, 123, mapper)
	mapIntTest[int8](t, 123, 1, 123, mapper)
	mapIntTest[int16](t, 123, 1, 123, mapper)
	mapIntTest[int32](t, 123, 1, 123, mapper)
	mapIntTest[int64](t, 123, 1, 123, mapper)
	mapIntTest[uint](t, 123, 1, 123, mapper)
	mapIntTest[uint8](t, 123, 1, 123, mapper)
	mapIntTest[uint16](t, 123, 1, 123, mapper)
	mapIntTest[uint32](t, 123, 1, 123, mapper)
	mapIntTest[uint64](t, 123, 1, 123, mapper)
	mapIntTest[uintptr](t, 123, 1, 123, mapper)

	mapSingleTypeTwoWayTest[bool, string](t, true, "true", mapper)
	mapSingleTypeTwoWayTest[bool, string](t, false, "false", mapper)

	mapSingleTypeTwoWayTest[float32, float64](t, 22, 22, mapper)
}

func mapIntTest[Src any](t *testing.T, srcValue Src, boolSrcValue Src, expected int, mapper *Mapper) {
	mapSingleTypeTwoWayTest[Src, int](t, srcValue, expected, mapper)
	mapSingleTypeTwoWayTest[Src, int8](t, srcValue, int8(expected), mapper)
	mapSingleTypeTwoWayTest[Src, int16](t, srcValue, int16(expected), mapper)
	mapSingleTypeTwoWayTest[Src, int32](t, srcValue, int32(expected), mapper)
	mapSingleTypeTwoWayTest[Src, int64](t, srcValue, int64(expected), mapper)
	mapSingleTypeTwoWayTest[Src, uint](t, srcValue, uint(expected), mapper)
	mapSingleTypeTwoWayTest[Src, uint8](t, srcValue, uint8(expected), mapper)
	mapSingleTypeTwoWayTest[Src, uint16](t, srcValue, uint16(expected), mapper)
	mapSingleTypeTwoWayTest[Src, uint32](t, srcValue, uint32(expected), mapper)
	mapSingleTypeTwoWayTest[Src, uint64](t, srcValue, uint64(expected), mapper)
	mapSingleTypeTwoWayTest[Src, uintptr](t, srcValue, uintptr(expected), mapper)
	mapSingleTypeTwoWayTest[Src, string](t, srcValue, strconv.Itoa(expected), mapper)
	mapSingleTypeTwoWayTest[Src, bool](t, boolSrcValue, true, mapper)

	mapSingleTypeTest[Src, complex64](t, srcValue, complex64(complex(float32(expected), 0)), mapper)
	mapSingleTypeTest[Src, complex128](t, srcValue, complex128(complex(float64(expected), 0)), mapper)
}

func mapSingleTypeTwoWayTest[Src any, Dst any](t *testing.T, srcValue Src, dstValue Dst, mapper *Mapper) {
	mapSingleTypeTest[Src, Dst](t, srcValue, dstValue, mapper)
	mapSingleTypeTest[Dst, Src](t, dstValue, srcValue, mapper)
}

func mapSingleTypeTest[SrcType any, DstType any](t *testing.T, value SrcType, expected DstType, mapper *Mapper) {
	a := struct {
		Field SrcType `map:"key"`
	}{
		Field: value,
	}
	b := struct {
		Field DstType `map:"key"`
	}{}
	if mapper == nil {
		m := MakeMapper()
		mapper = &m
	}
	if t != nil {
		assert.NoError(t, mapper.Map(&a, &b))
		assert.Equal(t, expected, b.Field)
	} else {
		if mapper.Map(&a, &b) != nil {
			panic("not equal 1")
		}
		if !reflect.DeepEqual(expected, b.Field) {
			panic("not equal 2")
		}
	}
}

func TestMapper_Map_ShouldMapBySrcUnderlyingType(t *testing.T) {
	var i = 3
	var str = "str"
	a := struct {
		Field  *int `map:"key"`
		Field2 *int `map:"key2"`
	}{
		Field: &i,
	}
	b := struct {
		Field  string  `map:"key"`
		Field2 *string `map:"key2"`
	}{
		Field2: &str,
	}
	mapper := MakeMapper()
	mapper.Converters = UnionSlices([]ValueConverterDesc{
		NewConverterDesc[int, *string](
			func(_ *reflect.Value, _, _ reflect.Type) (*reflect.Value, error) {
				str := "str"
				ptr := &str
				result := reflect.ValueOf(ptr)
				return &result, nil
			},
		),
		/*
			{
				Src: TypeSignature{reflect.Int},
				Dst: TypeSignature{reflect.Pointer, reflect.String},
				Converter: func(_ *reflect.Value, _, _ TypeSignature) (*reflect.Value, error) {
					str := "str"
					ptr := &str
					result := reflect.ValueOf(ptr)
					return &result, nil
				},
			},*/
	}, mapper.Converters)

	assert.NoError(t, mapper.Map(&a, &b))
	assert.Equal(t, b.Field, strconv.Itoa(i))
	assert.Nil(t, b.Field2)

}

func TestMapper_Map_ShouldMapByDstUnderlyingType(t *testing.T) {
	var i = 3
	a := struct {
		Field int `map:"key"`
	}{
		Field: i,
	}
	b := struct {
		Field *string `map:"key"`
	}{}
	mapper := MakeMapper()
	assert.NoError(t, mapper.Map(&a, &b))
	assert.NotNil(t, b.Field)
	if b.Field != nil {
		assert.Equal(t, *b.Field, strconv.Itoa(i))
	}
}

func TestMapper_Map_ShouldMapBySrcAndDstUnderlyingTypes(t *testing.T) {
	var i = 3
	var str = "str"
	a := struct {
		Field  *int `map:"key"`
		Field2 *int `map:"key2"`
	}{
		Field: &i,
	}
	b := struct {
		Field  *string `map:"key"`
		Field2 *string `map:"key2"`
	}{
		Field2: &str,
	}
	mapper := MakeMapper()
	assert.NoError(t, mapper.Map(&a, &b))
	assert.NotNil(t, b.Field)
	if b.Field != nil {
		assert.Equal(t, *b.Field, strconv.Itoa(i))
	}
	assert.Nil(t, b.Field2)
}

func TestMapper_Map_ShouldSkipNotTaggedFields(t *testing.T) {
	a := struct {
		Field int `map:"key"`
	}{
		Field: 1,
	}
	b := struct {
		Ignored string
		Field   string `map:"key"`
	}{
		Ignored: "existing value",
	}
	mapper := MakeMapper()
	assert.NoError(t, mapper.Map(&a, &b))
	assert.Equal(t, "existing value", b.Ignored)
}

func TestMapper_Map_FieldOrderShouldNotAffect(t *testing.T) {
	a := struct {
		Field1 int `map:"key1"`
		Field2 int `map:"key2"`
	}{
		Field1: 1,
		Field2: 2,
	}
	b := struct {
		Ignored string
		Field2  string `map:"key2"`
		Field1  string `map:"key1"`
	}{
		Ignored: "existing value",
	}
	mapper := MakeMapper()
	assert.NoError(t, mapper.Map(&a, &b))
	assert.Equal(t, "2", b.Field2)
	assert.Equal(t, "1", b.Field1)
}

func TestMapper_Map_ShouldMapNotInitializedPointers(t *testing.T) {
	a := struct {
		Field *int `map:"key"`
	}{}
	str := "str"
	b := struct {
		Field *string `map:"key"`
	}{
		Field: &str,
	}
	mapper := MakeMapper()
	assert.NoError(t, mapper.Map(&a, &b))
	assert.Nil(t, b.Field)
}

func TestMapper_Map_ShouldMapSlice(t *testing.T) {
	a := struct {
		Field []int `map:"key"`
	}{
		Field: []int{1, 2, 3},
	}
	b := struct {
		Field []string `map:"key"`
	}{}
	mapper := MakeMapper()
	assert.NoError(t, mapper.Map(&a, &b))
	assert.NotNil(t, b.Field)
	if b.Field != nil {
		assert.Equal(t, true, reflect.DeepEqual(b.Field, []string{"1", "2", "3"}))
	}
}

func TestMapper_Map_ShouldMapNilSlice(t *testing.T) {
	a := struct {
		Field []int `map:"key"`
	}{
		Field: nil,
	}
	b := struct {
		Field []string `map:"key"`
	}{
		Field: []string{"111"},
	}
	mapper := MakeMapper()
	assert.NoError(t, mapper.Map(&a, &b))
	assert.Nil(t, b.Field)
}

func TestMapper_Map_CanReturnConverterNotFoundError(t *testing.T) {
	a := struct {
		Field int `map:"key"`
	}{
		Field: 1,
	}
	b := struct {
		Field string `map:"key"`
	}{}
	mapper := MakeMapper()
	mapper.Converters = make([]ValueConverterDesc, 0)
	err := mapper.Map(&a, &b)
	if assert.Error(t, err) {
		assert.Equal(t, NewConverterNotFoundError(GetReflectType[int](), GetReflectType[string]()), err)
		//assert.Equal(t, NewConverterNotFoundError([]reflect.Kind{reflect.Int}, []reflect.Kind{reflect.String}), err)
	}
}

func TestMapper_Map_CanReturnNotAssignableValueError(t *testing.T) {
	a := struct {
		Field string `map:"key"`
	}{
		Field: "1",
	}
	b := struct {
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
	err := mapper.Map(&a, &b)
	if assert.Error(t, err) {
		assert.Equal(
			t,
			NewNotAssignableValueError(
				"Field",
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

func TestMapper_Map_CanReturnConverterError(t *testing.T) {
	a := struct {
		Field string `map:"key"`
	}{
		Field: "1",
	}
	b := struct {
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
	err := mapper.Map(&a, &b)
	if assert.Error(t, err) {
		assert.Equal(
			t,
			NewConverterError(
				"Field",
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

func TestMapper_Map_ShouldWorkInMultiThreadsWithPlan(t *testing.T) {
	threadCount := 50

	var mapper = NewMapper()
	var done = make(chan bool, 1)
	for n := 0; n < threadCount; n++ {
		go mapWorker(done, mapper)
	}
	for n := 0; n < threadCount; n++ {
		<-done
	}
}

func TestMapper_Map_ShouldWorkInMultiThreadingWithoutPlan(t *testing.T) {
	threadCount := 50

	var done = make(chan bool, 1)
	for n := 0; n < threadCount; n++ {
		go mapWorker(done, nil)
	}
	for n := 0; n < threadCount; n++ {
		<-done
	}
}

func mapWorker(done chan bool, mapper *Mapper) {
	for n := 0; n < 50; n++ {
		mapBaseTypesTest(nil, mapper)
	}
	done <- true
}

func BenchmarkMapper_Map_WithPlan(b *testing.B) {
	aa := struct {
		Field int `map:"key"`
	}{
		Field: 3,
	}
	bb := struct {
		Field string `map:"key"`
	}{}

	mapper := MakeMapper()
	for n := 0; n < 10000; n++ {
		err := mapper.Map(&aa, &bb)
		if err != nil {
			panic(" ")
		}
	}
}

func BenchmarkMapper_Map_WithoutPlan(b *testing.B) {
	aa := struct {
		Field int `map:"key"`
	}{
		Field: 3,
	}
	bb := struct {
		Field string `map:"key"`
	}{}

	for n := 0; n < 10000; n++ {
		mapper := MakeMapper()
		err := mapper.Map(&aa, &bb)
		if err != nil {
			panic(" ")
		}
	}
}
