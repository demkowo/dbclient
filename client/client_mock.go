package dbclient

import (
	"errors"
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
	Error         map[string]error
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

	if err := compareArgs(mock.Args, args); err != nil {
		return nil, err
	}

	if mock.Error["Exec"] != nil {
		return nil, mock.Error["Exec"]
	}

	res := resultMock{Args: mock.Args}

	if err := checkExpectedValue(mock.ExpectedValue, args); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *clientMock) Query(query string, args ...any) (rows, error) {
	mock, exists := c.mocks[query]
	if !exists {
		return nil, errors.New("mock not found for query: " + query)
	}

	if err := compareArgs(mock.Args, args); err != nil {
		return nil, err
	}

	if mock.Error["Query"] != nil {
		return nil, mock.Error["Query"]
	}

	rows := rowsMock{
		Columns: mock.Columns,
		Rows:    mock.Rows,
		Error:   mock.Error,
	}

	if err := checkExpectedValue(mock.ExpectedValue, args); err != nil {
		return nil, err
	}

	return &rows, nil
}

func (c *clientMock) QueryRow(query string, args ...any) row {
	mock, exists := c.mocks[query]
	if !exists {
		log.Fatal("mock not found for query: " + query)
	}

	if err := compareArgs(mock.Args, args); err != nil {
		mock.Error["Scan"] = err
		return nil
	}

	if mock.Error["QueryRow"] != nil {
		return &rowMock{
			Error: mock.Error,
		}
	}

	if err := checkExpectedValue(mock.ExpectedValue, args); err != nil {
		mock.Error["Scan"] = err
		return nil
	}

	return &rowMock{
		Args:    mock.Args,
		Columns: mock.Columns,
		Rows:    mock.Rows,
		Error:   mock.Error,
	}
}

func compareArgs(expected, actual []interface{}) error {
	// TODO: solve issue with time.Now() executed in both test case and tested method
	if len(expected) > len(actual) {
		return errors.New(fmt.Sprint("number of arguments in mock can't be higher than number of arguments in method:\n", expected, "\n", actual))
	} else if len(expected) < len(actual) {
		return errors.New(fmt.Sprint("number of arguments in method can't be higher than number of arguments in mock:\n", expected, "\n", actual))
	}

	ok := true
	for i, v := range expected {
		if reflect.TypeOf(v) == reflect.TypeOf(time.Now()) {
			continue
		}
		if !reflect.DeepEqual(v, actual[i]) {
			ok = false
		}
	}

	if !ok {
		return errors.New("invalid mock, args in mock and method are not equal: ")
	}
	return nil
}

func checkExpectedValue(expectedValue map[string]interface{}, args []interface{}) error {
	if expectedValue["index"] != nil && len(args) != 0 {
		index := expectedValue["index"]
		i, ok := index.(int)
		if !ok {
			return errors.New("ExpectedValue[\"index\"] should be of type int")
		}

		if !reflect.DeepEqual(args[i], expectedValue["value"]) {
			return fmt.Errorf("\n\texpected vaule: %v\n\treceived value: %v", expectedValue["value"], args[i])
		}
	}
	return nil
}
