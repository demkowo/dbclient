package dbclient

import (
	"database/sql"
)

type row interface {
	Scan(destinations ...interface{}) error
	Err() error
}

type dbRow struct {
	row *sql.Row
}

func (r *dbRow) Scan(destinations ...interface{}) error {
	return r.row.Scan(destinations...)
}

func (r *dbRow) Err() error {
	return r.row.Err()
}
