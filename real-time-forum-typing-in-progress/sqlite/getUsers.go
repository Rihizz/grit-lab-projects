package sqlite

type User struct {
	ID       int
	Username string
}

func GetUsers() ([]User, error) {
	db, err := OpenDb()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT user_id, username FROM users ORDER BY username COLLATE NOCASE ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
