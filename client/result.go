package dbclient

type result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}
