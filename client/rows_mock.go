package dbclient

import (
	"errors"
	"fmt"
	"reflect"
)

type rowsMock struct {
	Columns []string
	Rows    [][]interface{}
	RowsErr error

	rowIndex int
}

func (r *rowsMock) Next() bool {
	return r.rowIndex < len(r.Rows)
}

func (r *rowsMock) Close() error {
	return nil
}

func (r *rowsMock) Scan(destinations ...interface{}) error {
	if r.rowIndex >= len(r.Rows) {
		return errors.New("no more rows to scan")
	}

	row := r.Rows[r.rowIndex]
	if len(row) != len(destinations) {
		return errors.New("invalid destinations length")
	}

	for i, v := range row {
		destVal := reflect.ValueOf(destinations[i])
		if destVal.Kind() != reflect.Ptr || destVal.IsNil() {
			return fmt.Errorf("destination at index %d is not a valid pointer", i)
		}

		expectedType := destVal.Elem().Type()
		actualValue := r.Rows[r.rowIndex][i]
		actualType := reflect.TypeOf(actualValue)
		if expectedType != actualType {
			return fmt.Errorf("type mismatch for column '%s': expected %v, got %v",
				r.Columns[i], expectedType, actualType)
		}

		destVal.Elem().Set(reflect.ValueOf(v))
	}

	r.rowIndex++
	return nil
}

func (r *rowsMock) Err() error {
	return r.RowsErr
}
