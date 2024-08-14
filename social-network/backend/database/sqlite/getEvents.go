package sqlite

import (
	"log"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite3 package but we won't use it directly
)

type Event struct {
	Id          int    `json:"id"`
	GroupId     int    `json:"groupId"`
	CreatorId   int    `json:"creatorId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DateTime    string `json:"dateTime"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

func GetEvents(groupId int) ([]Event, error) {
	db, err := OpenDb()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query("SELECT * FROM events WHERE group_id = ?", groupId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err = rows.Scan(&event.Id, &event.GroupId, &event.CreatorId, &event.Title, &event.Description, &event.DateTime, &event.CreatedAt, &event.UpdatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return events, nil
}
