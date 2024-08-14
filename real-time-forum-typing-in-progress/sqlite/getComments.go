package sqlite

type Comment struct {
	ID      int
	PostID  int
	UserID  int
	Content string
	Date    string
	Author  string
}

func GetComments(postID int) ([]Comment, error) {
	db, err := OpenDb()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT comment_id, post_id, user_id, content, date, author FROM comments WHERE post_id = ? ORDER BY date DESC", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.Date, &comment.Author); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
