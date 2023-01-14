package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type DB struct {
	Db *sql.DB
}

// OpenDB Take the database name and return a DB reference to data source
func OpenDB(host string, port int, user string, password string, dbname string) (*sql.DB, error) {
	// DB Connection config
	cfg := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", cfg)
	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	return db, nil
}
