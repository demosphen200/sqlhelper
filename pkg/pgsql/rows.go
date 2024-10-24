package pgsql

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"sqlhelper/internal/utils"
)

type PgRows struct {
	rows pgx.Rows
}

func (rs *PgRows) RawRows() any {
	return rs.rows
}

func (rs *PgRows) Close() error {
	rs.rows.Close()
	return nil
}

func (rs *PgRows) Next() bool {
	return rs.rows.Next()
}

func (rs *PgRows) Scan(dest ...any) error {
	return rs.rows.Scan(dest...)
}

func (rs *PgRows) Columns() ([]string, error) {
	fields := rs.rows.FieldDescriptions()
	return utils.MapSlice(fields, func(f pgconn.FieldDescription) string {
		return f.Name
	}), nil
}
