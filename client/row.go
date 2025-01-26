package dbclient

import "database/sql"

type dbRows struct {
	rows *sql.Rows
}

type rows interface {
	Next() bool
	Close() error
	Scan(descriptions ...interface{}) error
}

func (r *dbRows) Next() bool {
	return r.rows.Next()
}

func (r *dbRows) Close() error {
	return r.rows.Close()
}

func (r *dbRows) Scan(descriptions ...interface{}) error {
	return r.rows.Scan(descriptions...)
}
