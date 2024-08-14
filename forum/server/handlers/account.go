package handlers

import (
	"fmt"
	"forum/server/login"
	"forum/server/sqlite"
	"net/http"
	"text/template"
)

var PostByUser = map[string]postByUser{}

type postByUser struct {
	PostID   int
	UserID   int
	Title    string
	Content  string
	Date     string
	Author   string
	Category string
}

type UserLikedPosts struct {
	PostID   int
	UserID   int
	Title    string
	Content  string
	Date     string
	Author   string
	Category string
}

type UserLikedComment struct {
	CommentID int
	PostID    int
	UserID    int
	Content   string
	Date      string
	Author    string
}

func viewAccount(w http.ResponseWriter, r *http.Request) {
	// Get the user from the session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("Error: ", err)
		ErrorHandling(http.StatusInternalServerError, w, r)
		return
	}
	userd := login.Sessions[cookie.Value]

	// Get the user's posts
	posts, err := getPostsByUser(userd.Username)
	if err != nil {
		fmt.Println("Error: ", err)
		ErrorHandling(http.StatusInternalServerError, w, r)
		return
	}

	// Get the user's liked posts
	likedPosts, err := getLikedPostsByUser(userd.UserID)
	if err != nil {
		fmt.Println("Error: ", err)
		ErrorHandling(http.StatusInternalServerError, w, r)
		return
	}
	// Get the user's liked comments
	likedComments, err := getLikedCommentsByUser(userd.UserID)
	if err != nil {
		fmt.Println("Error: ", err)
		ErrorHandling(http.StatusInternalServerError, w, r)
		return
	}

	// Execute the template
	t, err := template.ParseFiles("templates/account.html")
	if err != nil {
		fmt.Println("Error: ", err)
		ErrorHandling(http.StatusInternalServerError, w, r)
		return
	}

	if !login.CheckSession(w, r) {
		ErrorHandling(401, w, r)
		return
	}

	t.Execute(w, map[string]interface{}{
		"Posts":         posts,
		"LikedPosts":    likedPosts,
		"LikedComments": likedComments,
		"user":          userd,
	})
}

func getPostsByUser(username string) ([]postByUser, error) {
	db, err := sqlite.OpenDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	posts := []postByUser{}
	rows, err := db.Query("SELECT * FROM posts WHERE author = ?", username)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		p := postByUser{}
		err := rows.Scan(&p.PostID, &p.UserID, &p.Title, &p.Content, &p.Date, &p.Category, &p.Author)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

// get user liked posts
func getLikedPostsByUser(userid int) ([]UserLikedPosts, error) {
	db, err := sqlite.OpenDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	posts := []UserLikedPosts{}
	rows, err := db.Query(`
		SELECT posts.* 
		FROM posts 
		JOIN reactions ON reactions.post_id = posts.post_id 
		WHERE reactions.user_id = ? AND reactions.value = ? AND reactions.comment_id = ?`, userid, 1, 0)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := UserLikedPosts{}
		err := rows.Scan(&p.PostID, &p.UserID, &p.Title, &p.Content, &p.Date, &p.Author, &p.Category)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func getLikedCommentsByUser(userid int) ([]UserLikedComment, error) {
	db, err := sqlite.OpenDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	posts := []UserLikedComment{}
	rows, err := db.Query(`
		SELECT comments.* 
		FROM comments 
		JOIN reactions ON reactions.comment_id = comments.comment_id 
		WHERE reactions.user_id = ? AND reactions.value = ?`, userid, 1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := UserLikedComment{}
		err := rows.Scan(&p.CommentID, &p.PostID, &p.UserID, &p.Content, &p.Date, &p.Author)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}
