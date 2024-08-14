package sqlite

func GetLikes(PostID int) (int, error) {
	db, err := OpenDb()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	// check if user has liked the post
	var likes int
	_ = db.QueryRow("SELECT COUNT(*) FROM reactions WHERE post_id = ? AND reaction_type = 1", PostID).Scan(&likes)
	return likes, nil
}
