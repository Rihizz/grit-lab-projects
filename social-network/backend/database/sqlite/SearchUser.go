package sqlite

import (
	"log"
)

func SearchUser(searchQuery string) ([]User, error) {
	db, err := OpenDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	var users []User

	rows, err := db.Query("SELECT user_id, fullname FROM users WHERE (firstname LIKE ? OR lastname LIKE ? OR fullname LIKE ?) OR email LIKE ?", searchQuery+"%", searchQuery+"%", searchQuery+"%", "%"+searchQuery+"%")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.FullName)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
