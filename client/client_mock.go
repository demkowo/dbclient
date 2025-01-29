package dbclient

import (
	"errors"
	"log"
	"reflect"
)

var (
	isMock bool
)

type clientMock struct {
	mocks map[string]Mock
}

type Mock struct {
	Query   string
	Args    []interface{}
	Error   error
	Columns []string
	Rows    [][]interface{}
	RowsErr error
}

func AddMock(mock Mock) {
	if !isMock {
		log.Println("mock server is off, therefore ignoring AddMock")
		return
	}

	client, ok := dbClient.(*clientMock)
	if !ok {
		log.Println("invalid type of clientMock")
		return
	}

	if client.mocks == nil {
		client.mocks = make(map[string]Mock)
	}

	client.mocks[mock.Query] = mock

}

func StartMock() {
	isMock = true
}

func StopMock() {
	isMock = false
}

func (c *clientMock) Close() error {
	return nil
}

func (c *clientMock) Exec(query string, args ...any) (result, error) {
	mock, exists := c.mocks[query]
	if !exists {
		return nil, errors.New("mock not found for query: " + query)
	}

	if !compareArgs(mock.Args, args) {
		return nil, errors.New("mock not found for query and args: " + query)
	}

	if mock.Error != nil {
		return nil, mock.Error
	}

	res := resultMock{Args: mock.Args}

	return &res, nil
}

func (c *clientMock) Query(query string, args ...any) (rows, error) {
	mock, exists := c.mocks[query]
	if !exists {
		return nil, errors.New("mock not found")
	}

	if !compareArgs(mock.Args, args) {
		return nil, errors.New("mock not found for query and args: " + query)
	}

	if mock.Error != nil {
		return nil, mock.Error
	}

	rows := rowsMock{
		Columns: mock.Columns,
		Rows:    mock.Rows,
		RowsErr: mock.RowsErr,
	}

	return &rows, nil
}

func (c *clientMock) QueryRow(query string, args ...any) row {
	mock, exists := c.mocks[query]
	if !exists {
		return &rowMock{RowsErr: errors.New("mock not found for query: " + query)}
	}

	if !compareArgs(mock.Args, args) {
		return &rowMock{
			RowsErr: errors.New("mock not found for query and args: " + query),
		}
	}

	if mock.Error != nil {
		return &rowMock{Error: mock.Error}
	}

	return &rowMock{
		Args:    mock.Args,
		Columns: mock.Columns,
		Rows:    mock.Rows,
		Error:   mock.Error,
	}
}

func compareArgs(expected, actual []interface{}) bool {
	if len(expected) != len(actual) {
		return false
	}
	for i, v := range expected {
		if !reflect.DeepEqual(v, actual[i]) {
			return false
		}
	}
	return true
}
