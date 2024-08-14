package util

import (
	"forum/server/sqlite"
	"log"
)

// counts all reactions for a post
func CountReactionsPosts(postID, reactionValue int) int {
	db, err := sqlite.OpenDb()

	if err != nil {
		log.Println(err)
		return 0
	}

	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM reactions WHERE post_id = ? AND value = ? AND comment_id = ?", postID, reactionValue, 0).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

// counts all reactions for a comment
func CountReactionsComments(postID, commentID, reactionValue int) int {
	db, err := sqlite.OpenDb()

	if err != nil {
		log.Println(err)
		return 0
	}

	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM reactions WHERE post_id = ? AND value = ? AND comment_id = ?", postID, reactionValue, commentID).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}
