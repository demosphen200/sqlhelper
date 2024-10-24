package sqlhelper

import (
	"context"
	"database/sql"
)

type SqlAdapter struct {
	db *sql.DB
}

func NewSqlAdapter(
	db *sql.DB,
) *SqlAdapter {
	return &SqlAdapter{
		db: db,
	}
}

func (a *SqlAdapter) ParamPlaceholder(index int) string {
	return "?"
}

func (a *SqlAdapter) Query(ctx context.Context, query string, args ...interface{}) (DbRows, error) {
	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &SqlRows{
		rows: rows,
	}, nil
}

func (a *SqlAdapter) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return a.db.ExecContext(ctx, query, args...)
}

func (a *SqlAdapter) RunInTransaction(
	ctx context.Context,
	options sql.TxOptions,
	block func(adapter DbAdapter) error,
) (err error) {
	tx, err := a.db.BeginTx(ctx, &options)
	if err != nil {
		return err
	}
	txAdapter := NewSqlTxAdapter(tx)

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			err = &PanicInTransactionError{Value: r}
		}
	}()

	err = block(txAdapter)
	if err == nil {
		_ = tx.Commit()
	} else {
		_ = tx.Rollback()
	}
	return err
}
