package sqlite

import (
	"log"
)

type GroupPost struct {
	Id        int
	GroupId   int
	UserId    int
	FullName  string
	Post      string
	CreatedAt string
	Image    []byte
}

func GetGroupPosts(groupId int) ([]GroupPost, error) {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query("SELECT id, group_id, user_id, content, image, created_at FROM group_posts WHERE group_id = ?", groupId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	var posts []GroupPost

	for rows.Next() {
		var post GroupPost
		err := rows.Scan(&post.Id, &post.GroupId, &post.UserId, &post.Post, &post.Image, &post.CreatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		// get the user's full name
		row := db.QueryRow("SELECT fullname FROM users WHERE user_id = ?", post.UserId)
		err = row.Scan(&post.FullName)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}
