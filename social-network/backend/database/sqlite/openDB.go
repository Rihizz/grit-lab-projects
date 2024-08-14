package sqlite

import (
	"database/sql"
	"log"
)

// Declare and initialize the Db variable
var Db = new(sql.DB)

// function to open the database pointing to var db
func OpenDb() (*sql.DB, error) {
	// open db
	db, err := sql.Open("sqlite3", "./database/allData.sqlite3")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Update the Db variable with the new database connection
	Db = db
	return db, nil
}
