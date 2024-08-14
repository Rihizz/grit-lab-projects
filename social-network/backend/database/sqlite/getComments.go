package sqlite

import (
	"encoding/base64"
	"log"
)

type CommentForResponse struct {
	CommentID int    `json:"comment_id"`
	PostID    int    `json:"post_id"`
	UserID    int    `json:"user_id"`
	FullName  string `json:"full_name"`
	Content   string `json:"content"`
	Image     string `json:"image"`
	CreatedAt string `json:"created_at"`
}

// GetComments gets all the comments from the database
func GetComments(PostID int, group bool) ([]CommentForResponse, error) {
	log.Println("GetComments")
	log.Println("PostID:", PostID)
	db, err := OpenDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	// create a slice of comments
	comments := []CommentForResponse{}

	if !group {

		// get all the comments from the database
		rows, err := db.Query("SELECT id, post_id, user_id, content, image, created_at FROM comments WHERE post_id = ?", PostID)
		if err != nil {
			return nil, err
		}

		// iterate over the rows
		for rows.Next() {
			// create a comment
			comment := CommentForResponse{}
			imageData := []byte{}

			// scan the rows into the comment
			err := rows.Scan(&comment.CommentID, &comment.PostID, &comment.UserID, &comment.Content, &imageData, &comment.CreatedAt)
			if err != nil {
				return nil, err
			}

			if imageData != nil {
				// encode the image data to base64
				comment.Image = base64.StdEncoding.EncodeToString(imageData)
			}
			log.Println("comment.UserID:", comment.UserID)
			db.QueryRow("SELECT fullname FROM users WHERE user_id = ?", comment.UserID).Scan(&comment.FullName)

			// append the comment to the slice
			comments = append(comments, comment)
		}
	} else {
		rows, err := db.Query("SELECT id, group_post_id, user_id, content, image, created_at FROM group_comments WHERE group_post_id = ?", PostID)
		if err != nil {
			return nil, err
		}

		// iterate over the rows
		for rows.Next() {
			// create a comment
			comment := CommentForResponse{}
			imageData := []byte{}

			// scan the rows into the comment
			err := rows.Scan(&comment.CommentID, &comment.PostID, &comment.UserID, &comment.Content, &imageData, &comment.CreatedAt)
			if err != nil {
				return nil, err
			}

			if imageData != nil {
				// encode the image data to base64
				comment.Image = base64.StdEncoding.EncodeToString(imageData)
			}

			log.Println("comment.UserID:", comment.UserID)
			db.QueryRow("SELECT fullname FROM users WHERE user_id = ?", comment.UserID).Scan(&comment.FullName)

			// append the comment to the slice
			comments = append(comments, comment)
		}
	}

	return comments, nil
}
