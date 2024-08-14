package sqlite

import "log"

func AddNotificationxd(userID int, text string) error {
	
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return err
	}

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO notifications (user_id, text, timestamp) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, text)
	if err != nil {
		return err
	}

	return nil
}