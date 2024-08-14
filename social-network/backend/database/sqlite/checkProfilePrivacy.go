package sqlite

import "log"

func CheckProfilePrivacy(userID int) (int, error) {
	db, err := OpenDb()
	if err != nil {
		log.Println("Error opening the database, CheckProfilePrivacy(): ", err)
		return 0, err
	}

	defer db.Close()

	// Check if the user has access to the post by querying the private_posts table
	query := "SELECT privacy FROM users WHERE user_id = $1"
	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println("Error checking the profile privacy, CheckProfilePrivacy(): ", err)
		return 0, err
	}

	defer rows.Close()

	var privacy int

	if rows.Next() {
		err := rows.Scan(&privacy)
		if err != nil {
			log.Println("Error scanning the profile privacy, CheckProfilePrivacy(): ", err)
			return 0, err
		}
	}

	return privacy, nil
}
