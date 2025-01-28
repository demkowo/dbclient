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

type client struct {
	db *sql.DB
}

type DbClient interface {
	Query(query string, args ...any) (rows, error)
	Exec(query string, args ...any) (result, error)
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

	dbClient := &client{
		db: database,
	}

	return dbClient, nil
}

func (c *client) Query(query string, args ...any) (rows, error) {
	rows, err := c.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	res := dbRows{
		rows: rows,
	}

	return &res, nil
}

func (c *client) Exec(query string, args ...any) (result, error) {
	return c.db.Exec(query, args...)
}

func isProduction() bool {
	return os.Getenv(goEnvironment) == production
}
