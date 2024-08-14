package sqlite

import (
	"log"
)

func CheckIfFollower(activeUser, profileId int) bool {
	if activeUser == profileId {
		return true
	}

	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return false
	}

	defer db.Close()

	var followerId int
	err = db.QueryRow("SELECT follower_id FROM followers WHERE user_id = ? AND follower_id = ? AND status = 2", profileId, activeUser).Scan(&followerId)
	if err != nil {
		log.Println(err)
		return false
	}

	if followerId == 0 {
		return false
	}

	log.Println("all good")

	return true
}
