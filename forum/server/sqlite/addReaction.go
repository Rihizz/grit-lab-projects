package sqlite

import (
	"log"
)

// function to add a reaction
func AddReaction(postid int, commentid int, user_id int, value int) {
	//open db
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return
	}

	//defer close until end of function
	defer db.Close()

	//if reaction already exists, update it
	check, err := reactionExists(postid, commentid, user_id)
	if err != nil {
		log.Println(err)
		return
	}

	//if reaction exists, update it
	if check {
		statement, err := db.Prepare("UPDATE reactions SET value = ? WHERE post_id = ? AND comment_id = ? AND user_id = ?")
		if err != nil {
			log.Println(err)
		}
		statement.Exec(value, postid, commentid, user_id)
		return
	}

	//insert reaction into table
	statement, err := db.Prepare("INSERT INTO reactions (post_id, comment_id, user_id, value) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
	}
	statement.Exec(postid, commentid, user_id, value)
}

func reactionExists(postid int, commentid int, userid int) (bool, error) {
	// open db
	db, err := OpenDb()
	if err != nil {
		return false, err
	}
	defer db.Close()

	// check if reaction exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM reactions WHERE post_id = ? AND comment_id = ? AND user_id = ?", postid, commentid, userid).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// function to delete a reaction
func DeleteReaction(postid int, commentid int, user_id int) {
	//open db
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return
	}

	//defer close until end of function
	defer db.Close()

	//delete reaction from table
	statement, err := db.Prepare("DELETE FROM reactions WHERE post_id = ? AND comment_id = ? AND user_id = ?")
	if err != nil {
		log.Println(err)
	}
	statement.Exec(postid, commentid, user_id)
}

// function to check if reaction is 1
func CheckReaction(postid int, commentid int, user_id int, value int) bool {
	//open db
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return false
	}

	//defer close until end of function
	defer db.Close()

	//check if reaction is 1 or 2
	err = db.QueryRow("SELECT value FROM reactions WHERE post_id = ? AND comment_id = ? AND user_id = ? AND value = ?", postid, commentid, user_id, value).Scan(&value)

	return err == nil
}
