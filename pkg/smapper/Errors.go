package smapper

import (
	"fmt"
	"reflect"
)

type ConverterNotFoundError struct {
	srcSignature reflect.Type
	dstSignature reflect.Type
}

func (error *ConverterNotFoundError) Error() string {
	return fmt.Sprintf(
		"converter not found for src type signature %v and dst type signature %v",
		error.srcSignature,
		error.dstSignature,
	)
}

func NewConverterNotFoundError(
	srcSignature reflect.Type,
	dstSignature reflect.Type,
) *ConverterNotFoundError {
	return &ConverterNotFoundError{
		srcSignature: srcSignature,
		dstSignature: dstSignature,
	}
}

type DestinationNotNullableError struct {
	srcIndex     int
	dstSignature reflect.Type
}

func (error *DestinationNotNullableError) Error() string {
	return fmt.Sprintf(
		"src field with index %v is nil, but destination is not nullable (%v)",
		error.srcIndex,
		error.dstSignature,
	)
}

func NewDestinationNotNullableError(
	srcIndex int,
	dstSignature reflect.Type,
) *DestinationNotNullableError {
	return &DestinationNotNullableError{
		srcIndex:     srcIndex,
		dstSignature: dstSignature,
	}
}

type NotAssignableValueError struct {
	srcFieldNameOrIndex string
	dstFieldName        string
	srcValue            interface{}
	srcSignature        reflect.Type
	dstValue            interface{}
	returnedSignature   reflect.Type
	dstSignature        reflect.Type
}

func (error *NotAssignableValueError) Error() string {
	return fmt.Sprintf(
		"converter returned not assignable value srcField(Name/Index)=%v dstField=%v srcValue=%v%v returnedValue=%v%v expectedType=%v",
		error.srcFieldNameOrIndex,
		error.dstFieldName,
		error.srcValue,
		error.srcSignature,
		error.dstValue,
		error.returnedSignature,
		error.dstSignature,
	)
}

func NewNotAssignableValueError(
	srcFieldNameOrIndex string,
	dstFieldName string,
	srcValue interface{},
	srcSignature reflect.Type,
	dstValue interface{},
	returnedSignature reflect.Type,
	dstSignature reflect.Type,
) *NotAssignableValueError {
	return &NotAssignableValueError{
		srcFieldNameOrIndex: srcFieldNameOrIndex,
		dstFieldName:        dstFieldName,
		srcValue:            srcValue,
		srcSignature:        srcSignature,
		dstValue:            dstValue,
		returnedSignature:   returnedSignature,
		dstSignature:        dstSignature,
	}
}

type ConverterError struct {
	srcFieldNameOrIndex string
	dstFieldName        string
	srcValue            interface{}
	srcSignature        reflect.Type
	//dstValue            interface{}
	//dstValueSignature   TypeSignature
	dstFieldSignature reflect.Type
	nestedError       error
}

func (error *ConverterError) Error() string {
	return fmt.Sprintf(
		"converter error on srcField=%v dstField=%v srcValue=%v%v expectedType=%v: %s",
		//"converter error on srcField=%v dstField=%v srcValue=%v%v returnedValue=%v%v expectedType=%v: %w",
		error.srcFieldNameOrIndex,
		error.dstFieldName,
		error.srcValue,
		error.srcSignature,
		//error.dstValue,
		//error.dstValueSignature,
		error.dstFieldSignature,
		error.nestedError.Error(),
	)
}

func NewConverterError(
	srcFieldNameOrIndex string,
	dstFieldName string,
	srcValue interface{},
	srcSignature reflect.Type,
	dstFieldSignature reflect.Type,
	nestedError error,
) *ConverterError {
	return &ConverterError{
		srcFieldNameOrIndex: srcFieldNameOrIndex,
		dstFieldName:        dstFieldName,
		srcValue:            srcValue,
		srcSignature:        srcSignature,
		dstFieldSignature:   dstFieldSignature,
		nestedError:         nestedError,
	}
}
