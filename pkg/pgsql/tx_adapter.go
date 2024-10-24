package pgsql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"sqlhelper/pkg/sqlhelper"
)

type PgxTxAdapter struct {
	tx pgx.Tx
}

func NewPgxTxAdapter(tx pgx.Tx) *PgxTxAdapter {
	return &PgxTxAdapter{
		tx: tx,
	}
}

func (a *PgxTxAdapter) ParamPlaceholder(index int) string {
	return fmt.Sprintf("$%d", index+1)
}

func (a *PgxTxAdapter) Query(ctx context.Context, query string, args ...interface{}) (sqlhelper.DbRows, error) {
	rows, err := a.tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &PgRows{
		rows: rows,
	}, nil
}

func (a *PgxTxAdapter) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	commandTag, err := a.tx.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return NewPgxResult(a.tx.Conn(), &commandTag), nil
}

func (a *PgxTxAdapter) RunInTransaction(
	ctx context.Context,
	options sql.TxOptions,
	block func(adapter sqlhelper.DbAdapter) error,
) error {
	return errors.New("nested transactions not supported")
}
