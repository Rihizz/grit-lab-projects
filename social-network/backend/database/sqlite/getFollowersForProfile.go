package sqlite

func GetFollowersForProfile(userId int) ([]User, error) {
	db, err := OpenDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query("SELECT user_id, fullname FROM users WHERE user_id IN (SELECT follower_id FROM followers WHERE user_id = ? AND status = 2)", userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var followers []User

	for rows.Next() {
		var follower User
		err := rows.Scan(&follower.Id, &follower.FullName)
		if err != nil {
			return nil, err
		}

		followers = append(followers, follower)
	}

	return followers, nil
}
