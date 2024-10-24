package smapper

import (
	"errors"
	"reflect"
	"sqlhelper/pkg/KvCache"
	"strconv"
)

type Mapper struct {
	TagKey            string
	Converters        []ValueConverterDesc
	mapPlans          KvCache.Cache2[reflect.Type, reflect.Type, *mapPlan]
	mapFromSlicePlans KvCache.Cache2[string, reflect.Type, *mapFromSlicePlan]
	//initialized       bool
}

func MakeMapper() Mapper {
	return Mapper{
		TagKey:            "map",
		Converters:        DefaultConverters,
		mapPlans:          KvCache.MakeCache2[reflect.Type, reflect.Type, *mapPlan](),
		mapFromSlicePlans: KvCache.MakeCache2[string, reflect.Type, *mapFromSlicePlan](),
	}
}

func NewMapper() *Mapper {
	m := MakeMapper()
	return &m
}

/*
func (mapper *Mapper) Init() {
	if mapper.initialized {
		return
	}
	if mapper.TagKey == "" {
		mapper.TagKey = "map"
	}
	if mapper.Converters == nil {
		mapper.Converters = DefaultConverters
	}
	//mapper.mapPlans = make(mapPlans)
	//mapper.mapFromSlicePlans = make(mapFromSlicePlans)
	mapper.mapPlans = KvCache.MakeCache2[reflect.Type, reflect.Type, *mapPlan]()
	mapper.mapFromSlicePlans = KvCache.MakeCache2[string, reflect.Type, *mapFromSlicePlan]()
	mapper.initialized = true
}

func (mapper *Mapper) ensureInitialized() {
	if !mapper.initialized {
		mapper.Init()
	}
}*/

func (mapper *Mapper) Map(srcPtr, dstPtr interface{}) error {
	//mapper.ensureInitialized()
	if plan, err := mapper.findOrCreateMapPlan(srcPtr, dstPtr); err != nil {
		return err
	} else {
		err = mapper.mapByPlan(plan, srcPtr, dstPtr)
		return err
	}
}

func (mapper *Mapper) findOrCreateMapPlan(srcPtr, dstPtr interface{}) (*mapPlan, error) {
	var err error
	srcType := reflect.TypeOf(srcPtr).Elem()
	src := reflect.ValueOf(srcPtr).Elem()
	dstType := reflect.TypeOf(dstPtr).Elem()
	dst := reflect.ValueOf(dstPtr).Elem()

	plan, planFound := mapper.mapPlans.Get(srcType, dstType)
	if !planFound {
		plan = &mapPlan{
			srcType:  srcType,
			dstType:  dstType,
			mappings: make([]mapFieldMapping, 0),
		}
	} else {
		return plan, nil
	}

	for t := 0; t < src.NumField(); t++ {
		srcStructField := srcType.Field(t)
		fieldId := srcType.Field(t).Tag.Get(mapper.TagKey)
		if fieldId == "" {
			continue
		}
		dstFieldIndex := mapper.findFieldIndex(dstType, fieldId)
		if dstFieldIndex == -1 {
			continue
		}
		dstStructField := dstType.Field(dstFieldIndex)

		//srcFieldSignature := GetTypeSignature(srcStructField.Type)
		//dstFieldSignature := GetTypeSignature(dstStructField.Type)
		srcFieldSignature := srcStructField.Type
		dstFieldSignature := dstStructField.Type

		dstField := dst.Field(dstFieldIndex)

		converter := mapper.findOrCreateConverter(srcFieldSignature, dstFieldSignature, dstField.Type())

		if converter != nil {
			plan.mappings = append(plan.mappings, mapFieldMapping{
				srcIndex:          t,
				dstIndex:          dstFieldIndex,
				srcFieldSignature: srcFieldSignature,
				dstFieldSignature: dstFieldSignature,
				converter:         converter,
			})
			continue
		}
		err = NewConverterNotFoundError(srcFieldSignature, dstFieldSignature)
		break
	}
	if err == nil {
		mapper.mapPlans.Put(srcType, dstType, plan)
		return plan, nil
	} else {
		return nil, err
	}
}

func (mapper *Mapper) mapByPlan(plan *mapPlan, srcPtr, dstPtr interface{}) error {
	src := reflect.ValueOf(srcPtr).Elem()
	dst := reflect.ValueOf(dstPtr).Elem()
	for _, mapping := range plan.mappings {
		srcField := src.Field(mapping.srcIndex)
		dstField := dst.Field(mapping.dstIndex)
		dstValue, err := mapping.converter(&srcField, mapping.srcFieldSignature, mapping.dstFieldSignature)
		if err != nil {
			return NewConverterError(
				reflect.TypeOf(srcPtr).Elem().Field(mapping.srcIndex).Name,
				reflect.TypeOf(dstPtr).Elem().Field(mapping.dstIndex).Name,
				src.Field(mapping.srcIndex).Interface(),
				mapping.srcFieldSignature,
				mapping.dstFieldSignature,
				err,
			)
		}
		if dstValue.Type().AssignableTo(dstField.Type()) {
			dstField.Set(*dstValue)
		} else {
			return NewNotAssignableValueError(
				reflect.TypeOf(srcPtr).Elem().Field(mapping.srcIndex).Name,
				reflect.TypeOf(dstPtr).Elem().Field(mapping.dstIndex).Name,
				src.Field(mapping.srcIndex).Interface(),
				mapping.srcFieldSignature,
				dstValue.Interface(),
				dstValue.Type(), //GetTypeSignature(dstValue.Type()),
				mapping.dstFieldSignature,
			)
		}
	}
	return nil
}

func (mapper *Mapper) findFieldIndex(typ reflect.Type, fieldId string) int {
	for index := 0; index < typ.NumField(); index++ {
		fieldType := typ.Field(index)
		if fieldType.Tag.Get(mapper.TagKey) == fieldId {
			return index
		}
	}
	return -1
}

func (mapper *Mapper) findOrCreateConverter(src, dst reflect.Type, dstType reflect.Type) ValueConverter {
	if src == dst {
		return SameTypeConverter
	}

	converter := mapper.findConverter(src, dst)
	if converter != nil {
		return converter
	}

	if src.Kind() == reflect.Pointer {
		converter = mapper.findConverter(src.Elem(), dst)
		if converter != nil {
			return mapper.createConverterBySrcUnderlyingType(dstType, converter)
		}
	}
	if dst.Kind() == reflect.Pointer {
		converter = mapper.findConverter(src, dst.Elem())
		if converter != nil {
			return mapper.createConverterByDstUnderlyingType(dstType, converter)
		}
	}
	if src.Kind() == reflect.Pointer && dst.Kind() == reflect.Pointer {
		converter = mapper.findConverter(src.Elem(), dst.Elem())
		if converter != nil {
			return mapper.createConverterBySrcAndDstUnderlyingType(dstType, converter)
		}
	}

	if src.Kind() == reflect.Slice && dst.Kind() == reflect.Slice {
		converter = mapper.findConverter(src.Elem(), dst.Elem())
		if converter != nil {
			return mapper.createConverterBySliceUnderlyingType(dstType, converter)
		}
	}

	return nil
}

func (mapper *Mapper) findConverter(src, dst reflect.Type) ValueConverter {
	if src == dst {
		return SameTypeConverter
	}
	for _, converterDesc := range mapper.Converters {
		if converterDesc.Src == src && converterDesc.Dst == dst {
			return converterDesc.Converter
		}
	}
	return nil
}

func (mapper *Mapper) createConverterBySrcUnderlyingType(
	dstType reflect.Type,
	srcUnderlyingConverter ValueConverter,
) ValueConverter {
	return func(value *reflect.Value, src, dst reflect.Type) (*reflect.Value, error) {
		if value.IsNil() {
			if dst.Kind() == reflect.Pointer {
				return CreateTypedNil(dstType), nil
			} else {
				return nil, errors.New("src value is nil, but dst value is not a pointer")
			}
		} else {
			underlyingValue := value.Elem()
			return srcUnderlyingConverter(&underlyingValue, src.Elem(), dst)
		}
	}
}

func (mapper *Mapper) createConverterByDstUnderlyingType(
	dstType reflect.Type,
	dstUnderlyingConverter ValueConverter,
) ValueConverter {
	return func(value *reflect.Value, src, dst reflect.Type) (*reflect.Value, error) {
		if reflectDstUnderlyingValue, err := dstUnderlyingConverter(value, src, dst.Elem()); err != nil {
			return nil, err
		} else {
			ptrToDstUnderlying := reflect.New(dstType.Elem())
			ptrToDstUnderlying.Elem().Set(*reflectDstUnderlyingValue)
			return &ptrToDstUnderlying, nil
		}
	}
}

func (mapper *Mapper) createConverterBySrcAndDstUnderlyingType(
	dstType reflect.Type,
	srcAndDstUnderlyingConverter ValueConverter,
) ValueConverter {
	return func(value *reflect.Value, src, dst reflect.Type) (*reflect.Value, error) {
		if value.IsNil() {
			return CreateTypedNil(dstType), nil
		} else {
			underlyingValue := value.Elem()
			dstUnderlyingValue, err := srcAndDstUnderlyingConverter(&underlyingValue, src.Elem(), dst.Elem())
			if err != nil {
				return nil, err
			}
			ptrToDstUnderlying := reflect.New(dstType.Elem())
			ptrToDstUnderlying.Elem().Set(*dstUnderlyingValue)
			return &ptrToDstUnderlying, nil
		}
	}
}

func (mapper *Mapper) createConverterBySliceUnderlyingType(
	dstSliceType reflect.Type,
	sliceItemConverter ValueConverter,
) ValueConverter {
	return func(value *reflect.Value, src, dst reflect.Type) (*reflect.Value, error) {
		if value.IsNil() {
			ptrToTypedDstNil := reflect.New(dstSliceType)
			typedDstNil := ptrToTypedDstNil.Elem()
			return &typedDstNil, nil
		} else {
			dstSlice := reflect.MakeSlice(reflect.SliceOf(dstSliceType.Elem()), 0, value.Len())
			for index := 0; index < value.Len(); index++ {
				itemValue := value.Index(index)
				if dstItemValue, err := sliceItemConverter(&itemValue, src.Elem(), dst.Elem()); err != nil {
					return nil, err
				} else {
					dstSlice = reflect.Append(dstSlice, *dstItemValue)
				}
			}
			return &dstSlice, nil
		}
	}
}

func (mapper *Mapper) GetTaggedFieldNames(structPtr interface{}) []string {
	/*if !mapper.initialized {
		mapper.Init()
	}*/
	result := make([]string, 0)
	srcType := reflect.TypeOf(structPtr).Elem()
	for n := 0; n < srcType.NumField(); n++ {
		srcStructField := srcType.Field(n)
		if srcStructField.Tag.Get(mapper.TagKey) != "" {
			result = append(result, srcStructField.Name)
		}
	}
	return result
}

func (mapper *Mapper) GetKeys(structPtr interface{}) []string {
	//mapper.ensureInitialized()
	result := make([]string, 0)
	srcType := reflect.TypeOf(structPtr).Elem()
	for n := 0; n < srcType.NumField(); n++ {
		srcStructField := srcType.Field(n)
		key := srcStructField.Tag.Get(mapper.TagKey)
		if key != "" {
			result = append(result, key)
		}
	}
	return result
}

func (mapper *Mapper) MapToSlice(srcPtr interface{}) []any {
	//mapper.ensureInitialized()

	result := make([]any, 0)
	srcType := reflect.TypeOf(srcPtr).Elem()
	src := reflect.ValueOf(srcPtr).Elem()
	for n := 0; n < srcType.NumField(); n++ {
		srcStructField := srcType.Field(n)
		if !srcStructField.IsExported() {
			continue
		}
		key := srcStructField.Tag.Get(mapper.TagKey)
		if key != "" {
			srcField := src.Field(n)
			result = append(result, srcField.Interface())
		}
	}
	return result
}

func (mapper *Mapper) MapFromSlice(srcId string, src []any, dstPtr interface{}) error {
	//mapper.ensureInitialized()
	if plan, err := mapper.findOrCreateFromSlicePlan(srcId, src, dstPtr); err != nil {
		return err
	} else {
		return mapper.mapFromSliceByPlan(src, dstPtr, plan)
	}
}

func (mapper *Mapper) findOrCreateFromSlicePlan(srcId string, src []any, dstPtr interface{}) (*mapFromSlicePlan, error) {
	dstType := reflect.TypeOf(dstPtr).Elem()
	dst := reflect.ValueOf(dstPtr).Elem()

	if plan, found := mapper.mapFromSlicePlans.Get(srcId, dstType); found {
		return plan, nil
	}
	plan := mapFromSlicePlan{
		dstType:  dstType,
		mappings: make([]mapFromSliceMapping, 0),
	}

	srcIndex := 0
	for n := 0; n < dstType.NumField(); n++ {
		dstStructField := dstType.Field(n)
		dstField := dst.Field(n)
		if !dstStructField.IsExported() {
			continue
		}
		key := dstStructField.Tag.Get(mapper.TagKey)
		if key == "" {
			continue
		}

		var converter ValueConverter
		var srcSignature reflect.Type
		dstSignature := dstStructField.Type //GetTypeSignature(dstStructField.Type)
		if src[srcIndex] != nil {
			srcSignature = reflect.TypeOf(src[srcIndex]) //GetTypeSignature(reflect.TypeOf(src[srcIndex]))
			converter := mapper.findOrCreateConverter(srcSignature, dstSignature, dstField.Type())
			if converter == nil {
				return nil, NewConverterNotFoundError(srcSignature, dstSignature)
			}
		} else {
			if !IsNullable(dstSignature.Kind()) {
				return nil, NewDestinationNotNullableError(srcIndex, dstSignature)
			}
		}
		plan.mappings = append(plan.mappings, mapFromSliceMapping{
			dstIndex:     n,
			srcSignature: srcSignature,
			dstSignature: dstSignature,
			converter:    converter,
		})
		srcIndex++
	}
	mapper.mapFromSlicePlans.Put(srcId, dstType, &plan)
	return &plan, nil
}

func (mapper *Mapper) mapFromSliceByPlan(src []interface{}, dstPtr interface{}, plan *mapFromSlicePlan) error {
	dstStruct := reflect.ValueOf(dstPtr).Elem()

	for srcIndex, mapping := range plan.mappings {
		dstField := dstStruct.Field(mapping.dstIndex)

		if mapping.converter == nil && src[srcIndex] != nil {
			srcSignature := reflect.TypeOf(src[srcIndex]) //GetTypeSignature(reflect.TypeOf(src[srcIndex]))
			converter := mapper.findOrCreateConverter(srcSignature, mapping.dstSignature, dstField.Type())
			if converter == nil {
				return NewConverterNotFoundError(srcSignature, mapping.dstSignature)
			}
			mapping.srcSignature = srcSignature
			mapping.converter = converter
		}

		var dstValue *reflect.Value
		if mapping.converter == nil {
			dstValue = CreateTypedNil(dstField.Type())
		} else {
			srcValue := reflect.ValueOf(src[srcIndex])
			var err error
			dstValue, err = mapping.converter(&srcValue, mapping.srcSignature, mapping.dstSignature)
			if err != nil {
				return NewConverterError(
					strconv.Itoa(srcIndex),
					reflect.TypeOf(dstPtr).Elem().Field(mapping.dstIndex).Name,
					src[srcIndex],
					mapping.srcSignature,
					mapping.dstSignature,
					err,
				)
			}
		}
		if dstValue.Type().AssignableTo(dstField.Type()) {
			dstField.Set(*dstValue)
		} else {
			return NewNotAssignableValueError(
				strconv.Itoa(srcIndex),
				reflect.TypeOf(dstPtr).Elem().Field(mapping.dstIndex).Name,
				src[srcIndex],
				mapping.srcSignature,
				dstValue.Interface(),
				dstValue.Type(), //GetTypeSignature(dstValue.Type()),
				mapping.dstSignature,
			)
		}
	}
	return nil
}
