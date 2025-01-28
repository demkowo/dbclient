package dbclient

type resultMock struct {
	Args []interface{}
}

// LastInsertId() is not implemented and always returns 0 and nil.
// It's created to ensure that resultMock satisfies the sql.Result interface.
func (r *resultMock) LastInsertId() (int64, error) {
	return 0, nil
}

// RowsAffected() is not implemented and always returns 0 and nil.
// It's created to ensure that resultMock satisfies the sql.Result interface.
func (r *resultMock) RowsAffected() (int64, error) {
	return 0, nil
}
