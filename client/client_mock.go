package dbclient

import (
	"fmt"
	"log"
	"reflect"
	"time"
)

var (
	isMock bool
)

type clientMock struct {
	mocks map[string]Mock
}

type Mock struct {
	Query         string
	Args          []interface{}
	Columns       []string
	Rows          [][]interface{}
	ExpectedValue map[string]interface{}
	Error         error
	RowsErr       error
	ScanErr       error
}

func AddMock(mock Mock) {
	if !isMock {
		log.Println("mock server is off, therefore ignoring AddMock")
		return
	}

	client, ok := dbClient.(*clientMock)
	if !ok {
		log.Fatal("invalid type of clientMock")
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
		log.Fatal("mock not found for query: " + query)
	}

	compareArgs(mock.Args, args)

	if mock.Error != nil {
		return nil, mock.Error
	}

	res := resultMock{Args: mock.Args}

	checkExpectedValue(mock.ExpectedValue, args)

	return &res, nil
}

func (c *clientMock) Query(query string, args ...any) (rows, error) {
	mock, exists := c.mocks[query]
	if !exists {
		log.Fatal("mock not found for query: " + query)
	}

	compareArgs(mock.Args, args)

	if mock.Error != nil {
		return nil, mock.Error
	}

	rows := rowsMock{
		Columns: mock.Columns,
		Rows:    mock.Rows,
		RowsErr: mock.RowsErr,
		ScanErr: mock.ScanErr,
	}

	checkExpectedValue(mock.ExpectedValue, args)

	return &rows, nil
}

func (c *clientMock) QueryRow(query string, args ...any) row {
	mock, exists := c.mocks[query]
	if !exists {
		log.Fatal("mock not found for query: " + query)
	}

	compareArgs(mock.Args, args)

	if mock.Error != nil {
		return &rowMock{Error: mock.Error}
	}

	checkExpectedValue(mock.ExpectedValue, args)

	return &rowMock{
		Args:    mock.Args,
		Columns: mock.Columns,
		Rows:    mock.Rows,
		Error:   mock.Error,
	}
}

func compareArgs(expected, actual []interface{}) {
	// TODO: solve issue with time.Now() executed in both test case and tested method
	if len(expected) > len(actual) {
		log.Fatal("number of arguments in mock can't be higher than number of arguments in method:\n", expected, "\n", actual)
	} else if len(expected) < len(actual) {
		log.Fatal("number of arguments in method can't be higher than number of arguments in mock:\n", expected, "\n", actual)
	}

	ok := true
	for i, v := range expected {
		if reflect.TypeOf(v) == reflect.TypeOf(time.Now()) {
			continue
		}
		if !reflect.DeepEqual(v, actual[i]) {
			ok = false
			fmt.Println("\t", i+1, v, " != ", actual[i])
		}
	}

	if !ok {
		log.Fatal("invalid mock, args in mock and method are not equal: ")
	}
}

func checkExpectedValue(expectedValue map[string]interface{}, args []interface{}) {
	fmt.Println("=== mock client checkExpectedValue ===")

	if expectedValue["index"] != nil && len(args) != 0 {
		index := expectedValue["index"]
		i, ok := index.(int)
		if !ok {
			log.Fatal("ExpectedValue[\"index\"] should be of type int")
		}

		fmt.Println("===", args[i], expectedValue["value"])
		if !reflect.DeepEqual(args[i], expectedValue["value"]) {
			log.Fatalf("expected expectedVaule type should be equal %v while is %v", expectedValue["value"], args[i])
		}
	}
}
