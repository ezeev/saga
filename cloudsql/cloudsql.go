// Package cloudsql provides connections to Google Cloud SQL
package cloudsql

import (
	"database/sql"
	"os"
	"log"
)


// ConnStr constructs a Google Cloud SQL Connection String.
func ConnStr() string {

	//connectionName := mustGetenv("CLOUDSQL_CONNECTION_NAME")
	//user := mustGetenv("CLOUDSQL_USER")
	//password := os.Getenv("CLOUDSQL_PASSWORD") // NOTE: password may be empty
	//dbName := mustGetenv("CLOUDSQL_DB_NAME")

	conn_str := mustGetenv("CLOUDSQL_CONN_STR")
	return conn_str
}
// CloudSQLConnection Opens a Cloud SQL Connection and returns a pointer to it.
func CloudSQLConnection() (*sql.DB, error) {

	connStr := ConnStr()
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// mustGetevn is a helper function that returns an environment variable for the provided key, or an error if the variable does not exist.
func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("%s environment variable not set.", k)
	}
	return v
}
