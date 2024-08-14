package sqlite

import (
	"fmt"
	"log"
)

// GetEventResponses returns all responses for the given event

type EventResponse struct {
	UserId   int    `json:"user_id"`
	Response string `json:"response"`
}

func GetEventResponses(eventID int) (responses []EventResponse, err error) {
	log.Println("GetEventResponses called")
	// open the database connection
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
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

	fmt.Println("eventID in geteventresponses", eventID)

	// get all responses for the event
	rows, err := tx.Query(`
	SELECT user_id, response FROM event_responses WHERE event_id = ?
	`, eventID)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}

	// iterate over each row
	for rows.Next() {
		var response EventResponse
		err = rows.Scan(&response.UserId, &response.Response)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return
		}
		// append the response to the responses slice
		responses = append(responses, response)
	}
	fmt.Println("responses in geteventresponses", responses)
	// commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}

	return
}
