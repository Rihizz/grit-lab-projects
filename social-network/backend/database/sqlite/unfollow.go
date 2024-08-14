package sqlite

import "log"

func Unfollow(userID, activeUser int) error {
	log.Println("Unfollowing user", userID, "from user", activeUser)
	db, err := OpenDb()
	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec("DELETE FROM followers WHERE user_id = ? AND follower_id = ?", userID, activeUser)
	if err != nil {
		return err
	}

	return nil
}
