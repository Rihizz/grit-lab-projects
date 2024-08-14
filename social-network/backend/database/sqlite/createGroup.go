package sqlite

import (
	"log"
	"time"
)

func CreateGroup(groupName string, groupDescription string, creatorId int) error {
	time := time.Now().Format("2006-01-02 15:04:05")

	db, err := OpenDb()
	if err != nil {
		log.Print(err)
		return err
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO groups (creator_id, title, description, created_at ) VALUES (?, ?, ?, ?)", creatorId, groupName, groupDescription, time)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
