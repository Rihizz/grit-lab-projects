package sqlite

import "log"

type UserList struct {
	UserID   int    `json:"id"`
	Fullname string `json:"fullname"`
}

// Based on received user_id, returns a list of all users that the person is following or followed by
func GetAllContacts(userID string) ([]UserList, error) {
	log.Println("GetAllContacts called")
	db, err := OpenDb()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var users []UserList

	//query the table followers and get all the users that the user is following, or is followed by
	//the colums are user_id and follower_id
	//then populate the userlist that is to be returned with all the matches except the user_id itself

	// Query to get all the user_ids that the user is following
	followingQuery := `SELECT follower_id FROM followers WHERE user_id = ? AND status = 2`
	rows, err := db.Query(followingQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var followerID int
		err = rows.Scan(&followerID)
		if err != nil {
			return nil, err
		}

		user, err := GetUserById(followerID)
		if err != nil {
			log.Println("Error getting user by ID: ", err)
			continue
		}

		users = append(users, UserList{UserID: user.Id, Fullname: user.FullName})
	}

	// Query to get all the user_ids that are followed by the user
	followedByQuery := `SELECT user_id FROM followers WHERE follower_id = ? AND status = 2`
	rows, err = db.Query(followedByQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var followedByID int
		err = rows.Scan(&followedByID)
		if err != nil {
			return nil, err
		}

		user, err := GetUserById(followedByID)
		if err != nil {
			log.Println("Error getting user by ID: ", err)
			continue
		}

		users = append(users, UserList{UserID: user.Id, Fullname: user.FullName})
	}

	//remove duplicates from the userlist
	users = removeDuplicates(users)

	return users, nil
}

func removeDuplicates(users []UserList) []UserList {
	// Use map to record duplicates as we find them.
	encountered := map[UserList]bool{}
	result := []UserList{}

	for v := range users {
		if encountered[users[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[users[v]] = true
			// Append to result slice.
			result = append(result, users[v])
		}
	}
	// Return the new slice.
	return result
}
