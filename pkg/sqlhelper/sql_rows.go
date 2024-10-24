package sqlhelper

import "database/sql"

type SqlRows struct {
	rows *sql.Rows
}

func (rs *SqlRows) RawRows() any {
	return rs.rows
}

func (rs *SqlRows) Close() error {
	return rs.rows.Close()
}

func (rs *SqlRows) Next() bool {
	return rs.rows.Next()
}

func (rs *SqlRows) Scan(dest ...any) error {
	return rs.rows.Scan(dest...)
}

func (rs *SqlRows) Columns() ([]string, error) {
	return rs.rows.Columns()
}
