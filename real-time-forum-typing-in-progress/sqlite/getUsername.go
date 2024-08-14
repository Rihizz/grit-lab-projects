package sqlite

import "log"

func GetUsernameByID(id int) (string, error) {
	// open db
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return "", err
	}

	// defer the closing of the database connection
	defer db.Close()

	//check if username or email exists
	rows, err := db.Query("SELECT username FROM users WHERE user_id = ?", id)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer rows.Close()

	var username string
	if rows.Next() {
		rows.Scan(&username)
		return username, nil
	}

	return "", err
}
