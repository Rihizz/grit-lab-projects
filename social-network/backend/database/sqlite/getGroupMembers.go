package sqlite

import (
	"log"
)

type Member struct {
	Id       int    `json:"id"`
	FullName string `json:"full_name"`
}

func GetGroupMembers(groupID int) ([]Member, error) {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query("SELECT user_id FROM group_members WHERE group_id = ? AND status = 1", groupID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	members := make([]Member, 0)

	for rows.Next() {
		var member Member
		err := rows.Scan(&member.Id)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		// get the user's full name
		row := db.QueryRow("SELECT fullname FROM users WHERE user_id = ?", member.Id)
		err = row.Scan(&member.FullName)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil
}
