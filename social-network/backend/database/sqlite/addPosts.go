package sqlite

import (
	"fmt"
	"log"
)

// AddPosts adds a post to the database
func AddPosts(userid int, content string, created string, author string, privacyStatus string) error {
	db, err := OpenDb()
	if err != nil {
		log.Println("Error in opening AddPosts line 10: ", err)
		return err
	}
	defer db.Close()
	// create a prepared SQL statement
	stmt, err := db.Prepare("INSERT INTO posts(user_id, content, author, privacy, created_at) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Error in preparing AddPosts line 12: ", err)
		return err
	}
	// execute the prepared SQL statement
	_, err = stmt.Exec(userid, content, author, privacyStatus, created)
	if err != nil {
		log.Println("Error in executing AddPosts line 18: ", err)
		return err
	}
	fmt.Println("New post was created")
	return nil
}
