package sqlite

import (
	"log"
)

func GetUsersInGroup(groupID int) ([]User, error) {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer db.Close()

	// get the users
	var users []User

	rows, err := db.Query("SELECT user_id, fullname, email FROM users WHERE user_id IN (SELECT user_id FROM group_users WHERE group_id = ?)", groupID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user User

		err = rows.Scan(&user.Id, &user.FullName, &user.Email)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
