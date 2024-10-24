package sqlhelper

import (
	"context"
	"database/sql"
	"errors"
)

type SqlTxAdapter struct {
	db *sql.Tx
	//db  sqlQuerier
}

func NewSqlTxAdapter(
	db *sql.Tx,
	// db sqlQuerier,
) *SqlTxAdapter {
	return &SqlTxAdapter{
		db: db,
	}
}

func (a *SqlTxAdapter) ParamPlaceholder(index int) string {
	return "?"
}

func (a *SqlTxAdapter) Query(ctx context.Context, query string, args ...interface{}) (DbRows, error) {
	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &SqlRows{
		rows: rows,
	}, nil
}

func (a *SqlTxAdapter) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return a.db.ExecContext(ctx, query, args...)
}

func (a *SqlTxAdapter) RunInTransaction(
	ctx context.Context,
	options sql.TxOptions,
	block func(adapter DbAdapter) error,
) error {
	return errors.New("nested transactions not supported")
}
