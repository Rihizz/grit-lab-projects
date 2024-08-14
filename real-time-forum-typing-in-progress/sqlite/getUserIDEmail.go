package sqlite

import (
	"errors"
	"log"
)

var ErrEmailNotFound = errors.New("Email not found")

func GetUserIdByEmail(email string) (int, error) {
	// open db
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return -1, err
	}

	// defer the closing of the database connection
	defer db.Close()

	//check if email exists
	rows, err := db.Query("SELECT user_id FROM users WHERE email = ?", email)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	defer rows.Close()

	var id int
	if rows.Next() {
		rows.Scan(&id)
		return id, nil
	}

	return -1, ErrEmailNotFound
}
