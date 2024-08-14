package sqlite

import (
	"log"
)

func UpdateFollowerStatus(userID, activeUser, acceptStatus int) (err error) {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return err
	}

	defer db.Close()

	_, err = db.Exec("UPDATE followers SET status = ? WHERE user_id = ? AND follower_id = ?", acceptStatus, userID, activeUser)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
