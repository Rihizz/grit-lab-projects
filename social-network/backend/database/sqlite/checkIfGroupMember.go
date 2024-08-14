package sqlite

import (
	"log"
)

func CheckIfGroupMember(userID, groupID int) (bool, error) {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return false, err
	}

	defer db.Close()

	rows, err := db.Query("SELECT user_id FROM group_members WHERE user_id = ? AND group_id = ? AND status = 1", userID, groupID)
	if err != nil {
		log.Println(err)
		return false, err
	}

	defer rows.Close()

	if !rows.Next() {
		return false, nil
	}

	return true, nil
}
