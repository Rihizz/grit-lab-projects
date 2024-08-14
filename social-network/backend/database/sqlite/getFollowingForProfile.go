package sqlite

func GetFollowingForProfile(UserID int) ([]User, error) {
	db, err := OpenDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query("SELECT user_id, fullname FROM users WHERE user_id IN (SELECT user_id FROM followers WHERE follower_id = ? AND status = 2)", UserID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var following []User

	for rows.Next() {
		var follow User
		err := rows.Scan(&follow.Id, &follow.FullName)
		if err != nil {
			return nil, err
		}

		following = append(following, follow)
	}

	return following, nil
}