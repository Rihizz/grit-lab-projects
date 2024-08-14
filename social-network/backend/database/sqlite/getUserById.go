package sqlite

import (
	"database/sql"
	"log"
)

// User struct
type User struct {
	Id       int    `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

// GetUserById gets a users basic info by id
func GetUserById(id int) (User, error) {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return User{}, err
	}

	defer db.Close()

	// get the user
	var user User

	row := db.QueryRow("SELECT user_id, fullname, email FROM users WHERE user_id = ?", id)

	switch err := row.Scan(&user.Id, &user.FullName, &user.Email); err {
	case sql.ErrNoRows:
		log.Println("No rows were returned!")
		return user, err
	case nil:
		return user, nil
	default:
		log.Println("Error getting the user, GetUserById(): ", err)
		return user, err
	}
}
