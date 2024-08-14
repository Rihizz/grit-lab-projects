package sqlite

import (
	"log"
)

// function to register a user
func RegisterUser(email, username, password, role string) {
	// open the database connection
	db, err := OpenDb()
	if err != nil {
		return
	}

	// defer the closing of the database connection
	defer db.Close()

	// start a new transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}

	// insert user into table
	statement, err := tx.Prepare("INSERT INTO users (email, username, password, role) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}

	// execute the prepared statement and insert the new user
	result, err := statement.Exec(email, username, password, role)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}

	// check if the insert operation succeeded
	lastInsertID, err := result.LastInsertId()
	if err != nil || lastInsertID == 0 {
		log.Println(err)
		tx.Rollback()
		return
	}

	// commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}
}
