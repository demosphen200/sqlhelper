package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"sqlhelper/pkg/sqlhelper"
)

type PgxPoolAdapter struct {
	pool *pgxpool.Pool
}

func NewPgxPoolAdapter(pool *pgxpool.Pool) *PgxPoolAdapter {
	return &PgxPoolAdapter{
		pool: pool,
	}
}

func (a *PgxPoolAdapter) ParamPlaceholder(index int) string {
	return fmt.Sprintf("$%d", index+1)
}

func (a *PgxPoolAdapter) Query(ctx context.Context, query string, args ...interface{}) (sqlhelper.DbRows, error) {
	rows, err := a.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &PgRows{
		rows: rows,
	}, nil
}

func (a *PgxPoolAdapter) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	commandTag, err := a.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return NewPgxResult(nil, &commandTag), nil
}

func (a *PgxPoolAdapter) RunInTransaction(
	ctx context.Context,
	options sql.TxOptions,
	block func(adapter sqlhelper.DbAdapter) error,
) (err error) {
	pgxOptions := pgx.TxOptions{}

	switch options.Isolation {
	case sql.LevelDefault:
		pgxOptions.IsoLevel = pgx.ReadCommitted
	case sql.LevelReadCommitted:
		pgxOptions.IsoLevel = pgx.ReadCommitted
	case sql.LevelReadUncommitted:
		pgxOptions.IsoLevel = pgx.ReadUncommitted
	case sql.LevelRepeatableRead:
		pgxOptions.IsoLevel = pgx.RepeatableRead
	case sql.LevelSerializable:
		pgxOptions.IsoLevel = pgx.Serializable
	default:
		return sqlhelper.UnsupportedIsolationLevel
	}

	if options.ReadOnly {
		pgxOptions.AccessMode = pgx.ReadOnly
	} else {
		pgxOptions.AccessMode = pgx.ReadWrite
	}

	tx, err := a.pool.BeginTx(ctx, pgxOptions)
	if err != nil {
		return err
	}
	txAdapter := NewPgxTxAdapter(tx)

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
			err = &sqlhelper.PanicInTransactionError{Value: r}
		}
	}()

	err = block(txAdapter)
	if err == nil {
		_ = tx.Commit(ctx)
	} else {
		_ = tx.Rollback(ctx)
	}
	return err
}
