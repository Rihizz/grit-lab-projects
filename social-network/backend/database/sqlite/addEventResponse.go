package sqlite

import "log"

func AddEventResponse(eventID int, userID int, response string) {

	// open the database connection
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return
	}

	// defer the closing of the database connection
	defer db.Close()
	//if response is "going" set the response to 1, if response is "not going" set the response to 0
	if response == "going" {
		response = "1"
	} else if response == "not going" {
		response = "0"
	}

	// start a new transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}
	//adds the response to the database, if the response already exists, it will be replaced
	statement, err := tx.Prepare(`
	INSERT OR REPLACE INTO event_responses (event_id, user_id, response) 
	VALUES (?, ?, ?)
`)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return
	}

	// execute the prepared statement and insert the new user
	result, err := statement.Exec(eventID, userID, response)
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
	log.Println("Event response added successfully")

}
