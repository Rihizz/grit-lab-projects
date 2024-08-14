package sqlite

import (
	"log"
)

func AddNotification(userID int, SenderID int, content string, nType string) error {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return err
	}

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO notifications (user_id, sender_id, content, type) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, SenderID, content, nType)
	if err != nil {
		log.Println(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
