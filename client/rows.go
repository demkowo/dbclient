package dbclient

import (
	"database/sql"
)

type rows interface {
	Next() bool
	Close() error
	Scan(destination ...interface{}) error
	Err() error
}

type dbRows struct {
	rows *sql.Rows
}

func (r *dbRows) Next() bool {
	return r.rows.Next()
}

func (r *dbRows) Close() error {
	return r.rows.Close()
}

func (r *dbRows) Scan(destination ...interface{}) error {
	return r.rows.Scan(destination...)
}

func (r *dbRows) Err() error {
	return r.rows.Err()
}
