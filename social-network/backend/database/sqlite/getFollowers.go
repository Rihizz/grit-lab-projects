// Gets the followers of a user

package sqlite

func GetFollowers(userId int) ([]User, error) {
	db, err := OpenDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query("SELECT follower_id FROM followers WHERE user_id = ?", userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var followers []User

	for rows.Next() {
		var followerId int
		err := rows.Scan(&followerId)
		if err != nil {
			return nil, err
		}

		follower, err := GetUserById(followerId)
		if err != nil {
			return nil, err
		}

		followers = append(followers, follower)
	}

	return followers, nil
}
