package sqlite

import (
	"fmt"
	"log"
)

type Message struct {
	Sender    string
	Receiver  string
	Text      string
	Timestamp string
	ID        int
}

// AddComment adds a new comment to the database.
func AddMessage(sender, receiver, text, timestamp string) error {
	// Open the connection to the database.
	log.Println("inside addMessage.go")
	log.Println(sender, receiver, text, timestamp)
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Use a prepared statement to insert the comment into the database.
	stmt, err := db.Prepare("INSERT INTO messages (sender, receiver, text, timestamp) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the prepared statement.
	_, err = stmt.Exec(sender, receiver, text, timestamp)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}
