package postgres

import (
	"database/sql"
	"log"
)

func Open(credentials string) *sql.DB {
	db, err := sql.Open("postgres", credentials)
	if err != nil {
		log.Fatalln("Error connecting to database")
	}
	return db
}