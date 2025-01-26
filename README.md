# dbclient Library

A simple, extensible Go library that abstracts away direct database interactions with PostgreSQL. 
It offers both a production client (using the standard `database/sql` package) and a mock client for testing or development environments.

## Features

- **Mock support**: Toggle mocking with `StartMock()` and `StopMock()` to return predefined queries and results.
- **Production-ready**: Leverages `database/sql` to handle real database connections.
- **Lightweight abstraction**: Provides `DBClient` interface for easy substitution in your project.
- **Type safety checks**: Ensures scanned column values match the expected types during mocking.

## Usage Example

> **Note**: The following code snippet (similar to `main.go`) is **not** part of the library. It’s only to illustrate how you might use the `dbclient` library in your own application.

```go
package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	dbclient "github.com/demkowo/dbclient/client"
)

var (
	db dbclient.DbClient
)

type User struct {
	Id    int
	Email string
}

func init() {
	// Turn on mock mode for testing or development
	dbclient.StartMock()

	var err error
	db, err = dbclient.Open("postgres", os.Getenv("DB_CLIENT"))
	if err != nil {
		log.Panicln("can't open database connection:", err)
	}
}

func main() {
	user, err := GetUser(1)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("ID:", user.Id)
	fmt.Println("Email:", user.Email)
}

func GetUser(id int) (*User, error) {
	// Add a mock result for demonstration
	dbclient.AddMock(dbclient.Mock{
		// Error: errors.New("error creating query"),
        Query:   "SELECT id, email FROM users WHERE id=$1",
		Args:    []interface{}{id},
		Columns: []string{"id", "email"},
		Rows: [][]interface{}{
			{1, "email-1@test.com"},
			{2, "email-2@test.com"},
		},
	})

	rows, err := db.Query("SELECT id, email FROM users WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user User
	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Email); err != nil {
			return nil, err
		}
		if user.Id == id {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}
```

## How It Works

### Package Initialization
- By default, the library checks the environment to decide if it should run in production mode or mock mode (based on `GO_ENVIRONMENT` variable).
- If you explicitly call `StartMock()`, the library routes all `Open` calls to a `clientMock`.

### Opening a Database Connection
- `Open(driverName, dataSourceName)`:
    - If running in mock mode, returns a `clientMock`.
    - Otherwise, creates an actual database client (`client`) using `database/sql`.

### Query Execution
- `SqlClient.Query(query string, args ...any)` is the main entry point:
    - **Production** (`client`): Executes a real `db.Query` and returns a `dbRows` wrapper around sql.Rows.
    - **Mock** (`clientMock`): Checks if a matching mock query is registered. If found, returns `rowsMock` with the predefined data.

### Row Handling
- Both `dbRows` (production) and `rowsMock` (mock) implement the same `rows` interface:
    - `Next()`: Retrieves the next row.
    - `Scan(...)`: Copies column data into destination variables.
    - `Close()`: Closes the underlying rows resource.

## Workflow Overview
1. **Initialize:** Decide on production or mock mode.
2. **Open Connection:** `dbclient.Open(...)` to get a `DbClient`.
3. **Optionally Add Mocks:** If in mock mode, use `dbclient.AddMock(...)` to register - predefined responses.
4. **Execute Queries:** Call `Query(...)` on the `DbClient`.
5. **Process Results:** Use the returned `rows` interface to iterate and scan.

## Potential Improvements
- **Connection Pool Configuration:** Extend `Open` to allow customizing max open/idle - connections, timeouts, etc.
- **Logging:** Implement structured logging with log levels for debugging or - production monitoring.
- **Support for Other Drivers:** Currently tailored for PostgreSQL. You could - generalize for MySQL, SQLite, or others.
- **Error Handling Enhancements:** Provide more contextual information (e.g., query - string, arguments) on errors.
- **Transaction Support:** Wrap queries in transactions for atomic operations (e.g., - `BEGIN`, `COMMIT`, `ROLLBACK`).