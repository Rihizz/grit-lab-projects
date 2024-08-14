package sqlite

import (
	"log"
)

// function to add a post
func AddPost(title, content, date string, user_id int, author string, category string) error {
	// open db
	db, err := OpenDb()
	if err != nil {
		return err
	}

	// defer the closing of the database connection
	defer db.Close()

	// start a new transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	// insert post into table
	statement, err := tx.Prepare("INSERT INTO posts (title, content, user_id, category, date, author) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	// execute the prepared statement and insert the new post
	_, err = statement.Exec(title, content, user_id, category, date, author)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	// commit the transaction
	return tx.Commit()
}
