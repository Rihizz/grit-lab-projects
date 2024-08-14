package sqlite

import "log"

func AddFollower(follower, followee, privacy int) error {
	// open the database connection
	db, err := OpenDb()
	if err != nil {
		return err
	}

	// defer the closing of the database connection
	defer db.Close()

	// start a new transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	// check if the user is already following the user
	exists, err := checkFollow(follower, followee)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	// if the user is already following the user, return an error
	if exists {
		log.Println("User is already following the user")
		tx.Rollback()
		return err
	}

	status := 2
	if privacy == 1 {
		status = 1
	}
	log.Println("follower: ", follower)
	log.Println("followee: ", followee)

	// insert user into table
	statement, err := tx.Prepare("INSERT INTO followers (user_id, follower_id, status) VALUES (?, ?, ?)")
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	// execute the prepared statement and insert the new user
	result, err := statement.Exec(followee, follower, status)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	// check if the insert operation succeeded
	lastInsertID, err := result.LastInsertId()
	if err != nil || lastInsertID == 0 {
		log.Println(err)
		tx.Rollback()
		return err
	}

	// commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	return nil

}

// Function to check if a user is already following the user
func checkFollow(follower, followee int) (bool, error) {
	// open the database connection
	db, err := OpenDb()
	if err != nil {
		return false, err
	}

	// defer the closing of the database connection
	defer db.Close()

	// query the database
	rows, err := db.Query("SELECT * FROM followers WHERE user_id = ? AND follower_id = ?", followee, follower)
	if err != nil {
		return false, err
	}

	// close the statement
	defer rows.Close()

	// check if the user exists
	if rows.Next() {
		return true, nil
	}

	return false, nil
}