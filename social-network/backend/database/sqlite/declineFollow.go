package sqlite

import "log"

func DeclineFollow(userId int, followerId int) error {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return err
	}

	defer db.Close()

	_, err = db.Exec("UPDATE followers SET status = 0 WHERE user_id = ? AND follower_id = ?", userId, followerId)
	if err != nil {
    log.Fatal(err)
	}

	return nil
}