package sqlite

import (
	"fmt"
	"log"
)

// AddPosts adds a post to the database
func AddPosts2(userid int, content string, created string, author string, privacyStatus string, filename string, picture []byte) error {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(filename)
	defer db.Close()
	// create a prepared SQL statement
	stmt, err := db.Prepare("INSERT INTO posts(user_id, content, author, privacy, created_at, image) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		return err
	}
	// execute the prepared SQL statement
	_, err = stmt.Exec(userid, content, author, privacyStatus, created, picture)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("New post was created")
	return nil
}
