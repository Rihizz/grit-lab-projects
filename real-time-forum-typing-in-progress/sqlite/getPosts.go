package sqlite

type Post struct {
	ID       int
	UserID   int
	Title    string
	Content  string
	Date     string
	Category string
	Author   string
}

func GetPosts() ([]Post, error) {
	db, err := OpenDb()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT post_id, user_id, title, content, date, category, author FROM posts ORDER BY date DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Date, &post.Category, &post.Author); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
