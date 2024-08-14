package sqlite

import "log"

func GetUserIdByUsername(username string) (int, error) {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer db.Close()

	var userId int
	err = db.QueryRow("SELECT user_id FROM users WHERE username = ?", username).Scan(&userId)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return userId, nil
}
