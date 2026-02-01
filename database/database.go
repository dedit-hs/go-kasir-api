package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	db, error := sql.Open("postgres", connectionString)
	if error != nil {
		return nil, error
	}

	error = db.Ping()
	if error != nil {
		return nil, error
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	log.Println("Database connection initialized")
	return db, nil
}
