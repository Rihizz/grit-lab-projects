package sqlite

import (
	"log"
)

func AddGroupMember(userID, GroupID, status int) error {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return err
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO group_members (group_id, user_id, status) VALUES (?, ?, ?)", GroupID, userID, status)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
