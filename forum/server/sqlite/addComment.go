package sqlite

import (
	"fmt"
)

// AddComment adds a new comment to the database.
func AddComment(postid, userid int, content, date, author string) error {
	// Open the connection to the database.
	db, err := OpenDb()
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Use a prepared statement to insert the comment into the database.
	stmt, err := db.Prepare("INSERT INTO comments (post_id, user_id, content, date, author) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the prepared statement.
	_, err = stmt.Exec(postid, userid, content, date, author)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}
