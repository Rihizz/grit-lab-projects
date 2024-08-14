package sqlite

import "log"

func AcceptFollow(userId int, followerId int) error {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return err
	}

	defer db.Close()

	_, err = db.Exec("UPDATE followers SET status = 2 WHERE user_id = ? AND follower_id = ?", userId, followerId)
	if err != nil {
    log.Fatal(err)
	}

	return nil
}