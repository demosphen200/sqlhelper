package pgsql

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"sqlhelper/pkg/sqlhelper"
)

func NewPgxPoolSqlHelper(
	pool *pgxpool.Pool,
	tableName string,
) *sqlhelper.SqlHelper {
	adapter := NewPgxPoolAdapter(pool)
	return sqlhelper.NewSqlHelper(adapter, tableName)
}
