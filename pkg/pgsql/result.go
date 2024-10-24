package pgsql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PgResultStub struct {
	conn       *pgx.Conn
	commandTag *pgconn.CommandTag
}

func NewPgxResult(
	conn *pgx.Conn,
	commandTag *pgconn.CommandTag,
) *PgResultStub {
	return &PgResultStub{
		conn:       conn,
		commandTag: commandTag,
	}
}

func (r *PgResultStub) LastInsertId() (int64, error) {
	if r.conn != nil {
		var id int64
		err := r.conn.QueryRow(context.Background(), "select LASTVAL()").Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	}
	return 0, errors.New("connection is nil (not in transaction?)")
}

func (r *PgResultStub) RowsAffected() (int64, error) {
	if r.commandTag != nil {
		return r.commandTag.RowsAffected(), nil
	}
	return 0, errors.New("cannot get RowsAffected - no command tag")
}
