package sqlite

import "log"

func CheckViewPrivileges(userID int, postID int) (bool, error) {
	log.Println("CheckViewPrivileges called")
	db, err := OpenDb()
	if err != nil {
		return false, err
	}
	defer db.Close()

	// Check if the user has access to the post by querying the private_posts table
	query := "SELECT * FROM private_post WHERE user_id = $1 AND post_id = $2"
	rows, err := db.Query(query, userID, postID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}

	return false, nil
}
