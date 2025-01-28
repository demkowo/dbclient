package dbclient

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/lib/pq"
)

const (
	goEnvironment = "GO_ENVIRONMENT"
	production    = "production"
)

var (
	dbClient DbClient
)

type DbClient interface {
	Exec(query string, args ...any) (result, error)
	Query(query string, args ...any) (rows, error)
	QueryRow(query string, args ...any) row
}

type client struct {
	db *sql.DB
}

func Open(driverName string, dataSourceName string) (DbClient, error) {
	if !isProduction() && isMock {
		dbClient = &clientMock{}
		return dbClient, nil
	}

	if driverName == "" {
		return nil, errors.New("invalid driver name")
	}

	database, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	dbClient := &client{db: database}

	return dbClient, nil
}

func (c *client) Exec(query string, args ...any) (result, error) {
	return c.db.Exec(query, args...)
}

func (c *client) Query(query string, args ...any) (rows, error) {
	rows, err := c.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	res := dbRows{rows: rows}

	return &res, nil
}

func (c *client) QueryRow(query string, args ...any) row {
	row := c.db.QueryRow(query, args...)
	return &dbRow{row: row}
}

func isProduction() bool {
	return os.Getenv(goEnvironment) == production
}
