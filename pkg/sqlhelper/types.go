package sqlhelper

import (
	"context"
	"database/sql"
)

type DbRow struct {
	RawRow any
}

type DbRows interface {
	RawRows() any
	Close() error
	Next() bool
	Columns() ([]string, error)
	Scan(dest ...any) error
}

type DbAdapter interface {
	ParamPlaceholder(index int) string
	Query(ctx context.Context, query string, args ...interface{}) (DbRows, error)
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	RunInTransaction(
		ctx context.Context,
		options sql.TxOptions,
		block func(adapter DbAdapter) error,
	) error
}
