package core

import (
	"database/sql"
	"fmt"

	// in order to be able to connect to a mysql database.
	_ "github.com/go-sql-driver/mysql"
)

// CreateDatabaseConnection returns a new connection with the provided credentials
func CreateDatabaseConnection(user string, password string, databaseName string) (*sql.DB, error) {
	dbinfo := fmt.Sprintf("%s:%s@/%s", user, password, databaseName)
	db, err := sql.Open("mysql", dbinfo)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}
