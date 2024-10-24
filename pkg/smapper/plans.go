package smapper

import (
	"reflect"
)

type mapFieldMapping struct {
	srcIndex          int
	dstIndex          int
	srcFieldSignature reflect.Type
	dstFieldSignature reflect.Type
	converter         ValueConverter
}

type mapPlan struct {
	srcType  reflect.Type
	dstType  reflect.Type
	mappings []mapFieldMapping
}

type mapFromSliceMapping struct {
	dstIndex     int
	srcSignature reflect.Type
	dstSignature reflect.Type
	converter    ValueConverter
}

type mapFromSlicePlan struct {
	dstType  reflect.Type
	mappings []mapFromSliceMapping
}
