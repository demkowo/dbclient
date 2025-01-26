package dbclient

import (
	"errors"
	"log"
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
}

func (c *clientMock) Query(query string, args ...any) (rows, error) {
	mock, exists := c.mocks[query]
	if !exists {
		return nil, errors.New("mock not found")
	}

	if mock.Error != nil {
		return nil, mock.Error
	}

	rows := rowsMock{
		Columns: mock.Columns,
		Rows:    mock.Rows,
	}

	return &rows, nil
}

func AddMock(mock Mock) {
	if !isMock {
		log.Println("ignoring AddMock, because mock server is off")
		return
	}

	client, ok := dbClient.(*clientMock)
	if !ok {
		log.Println("invalid type of clientMock")
		return
	}

	if client.mocks == nil {
		client.mocks = make(map[string]Mock, 0)
	}

	client.mocks[mock.Query] = mock
}

func StartMock() {
	isMock = true
}

func StopMock() {
	isMock = false
}
