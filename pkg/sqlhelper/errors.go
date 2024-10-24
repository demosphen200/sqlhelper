package sqlhelper

import (
	"errors"
	"fmt"
)

type MustBePtrToStructError struct{}

func (err *MustBePtrToStructError) Error() string {
	return "ptrToModel must be pointer to struct"
}

type MustBePtrToSliceError struct{}

func (err *MustBePtrToSliceError) Error() string {
	return "ptrToSlice must be pointer to slice"
}

type MustBePtrToStructOrSliceError struct{}

func (err *MustBePtrToStructOrSliceError) Error() string {
	return "ptrToResult must be pointer to struct or pointer to slice"
}

type TableNameNotSetError struct{}

func (err *TableNameNotSetError) Error() string {
	return "table name not set"
}

type DbNotSetError struct{}

func (err *DbNotSetError) Error() string {
	return "db not set"
}

type IdFieldAndArgCountNotMatchError struct{}

func (err *IdFieldAndArgCountNotMatchError) Error() string {
	return "id field count and count of id args count does not match"
}

type NoRowsReturnedError struct{}

func (err *NoRowsReturnedError) Error() string {
	return "no rows returned"
}

type MoreThanOneRowReturnedError struct{}

func (err *MoreThanOneRowReturnedError) Error() string {
	return "query returned more than 1 record"
}

type TypeConverterNotFoundError struct {
	Name string
}

func (err *TypeConverterNotFoundError) Error() string {
	return fmt.Sprintf("converter \"%s\" not found", err.Name)
}

type PanicInTransactionError struct {
	Value any
}

func (err *PanicInTransactionError) Error() string {
	return fmt.Sprintf("panic in transaction: %v", err.Value)
}

var UnsupportedIsolationLevel = errors.New("unsupported isolation level")
