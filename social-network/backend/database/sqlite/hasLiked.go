package sqlite

import "log"

func HasLiked(userID, postID int) (int, error) {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	defer db.Close()

	// check if user has liked the post
	var status int
	_ = db.QueryRow("SELECT reaction_type FROM reactions WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&status)
	if status == 0 {
		return 0, nil
	}

	if status == 1 {
		return 1, nil
	} else {
		return 2, nil
	}
}
