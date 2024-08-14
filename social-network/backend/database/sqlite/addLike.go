package sqlite

import "log"

func AddLike(userID, postID int) error {
	log.Println("add like")
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return err
	}

	defer db.Close()

	likeStatus, err := HasLiked(userID, postID)
	if err != nil {
		log.Println(err)
		return err
	}

	// check if user has liked the post
	if likeStatus == 0 {
		log.Println("has not liked")
		_, err = db.Exec("INSERT INTO reactions (user_id, post_id, reaction_type) VALUES (?, ?, ?)", userID, postID, 1)
		if err != nil {
			log.Println(err)
			return err
		}
	} else if likeStatus == 1 {
		_, err = db.Exec("DELETE FROM reactions WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
