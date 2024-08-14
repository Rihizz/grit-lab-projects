package sqlite

import (
	"log"
	"time"
)

func CreateEvent(groupId int, creatorId int, eventTitle string, eventDescription string, eventDateTime string) error {
	time := time.Now().Format("2006-01-02 15:04:05")

	db, err := OpenDb()
	if err != nil {
		log.Print(err)
		return err
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO events (group_id, creator_id, title, description, date_time, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)", groupId, creatorId, eventTitle, eventDescription, eventDateTime, time, time)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
