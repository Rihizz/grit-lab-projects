package handlers

import (
	"fmt"
	"forum/server/login"
	"forum/server/sqlite"
	util "forum/server/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

var uniquePost = map[int]post{}

type post struct {
	PostID   int
	Title    string
	Content  string
	Date     string
	Category string
	Author   string
	Likes    int
	Dislikes int
}

type comment struct {
	CommentID int
	PostID    int
	UserID    int
	Content   string
	Date      string
	Author    string
	Likes     int
	Dislikes  int
}

func viewPost(w http.ResponseWriter, r *http.Request) {
	// Get the user from the session
	if login.CheckSession(w, r) {
		c, err := r.Cookie("session_token")
		if err != nil {
			log.Println(err)
			ErrorHandling(http.StatusInternalServerError, w, r)
			return
		}
		userd := login.Sessions[c.Value]

		// Parse the post ID from the URL
		id, err := parsePostID(r.URL.Path)
		if err != nil {
			log.Println(err)
			ErrorHandling(http.StatusBadRequest, w, r)
			return
		}

		// Search for the post based on the ID
		post, err := SearchPost(id, &uniquePost)
		if err != nil {
			log.Println(err)
			ErrorHandling(http.StatusInternalServerError, w, r)
			return
		}

		// Get the comments for the post
		comments, err := getComments(id)
		if err != nil {
			log.Println(err)
			ErrorHandling(http.StatusInternalServerError, w, r)
			return
		}
		// Set the current post ID
		currentPostID, err = parsePostID(r.URL.Path)
		if err != nil {
			log.Println(err)
			ErrorHandling(http.StatusInternalServerError, w, r)
			return
		}
		// Parse the form
		err = r.ParseForm()
		if err != nil {
			log.Println(err)
			ErrorHandling(http.StatusInternalServerError, w, r)
			return
		}

		// Render the post page
		tmpl, err := template.ParseFiles("templates/detail.html")
		if err != nil {
			ErrorHandling(http.StatusInternalServerError, w, r)
			return
		}
		tmpl.Execute(w, map[string]interface{}{
			"posts":    post,
			"comments": comments,
			"user":     userd,
		})
	} else {
		// Parse the post ID from the URL
		id, err := parsePostID(r.URL.Path)
		if err != nil {
			log.Println(err)
			ErrorHandling(http.StatusBadRequest, w, r)
			return
		}

		// Search for the post based on the ID
		post, err := SearchPost(id, &uniquePost)
		if err != nil {
			log.Println(err)
			ErrorHandling(http.StatusInternalServerError, w, r)
			return
		}

		// Get the comments for the post
		comments, err := getComments(id)
		if err != nil {
			log.Println(err)
			ErrorHandling(http.StatusInternalServerError, w, r)
			return
		}
		// Set the current post ID
		currentPostID, err = parsePostID(r.URL.Path)
		if err != nil {
			log.Println(err)
			ErrorHandling(http.StatusInternalServerError, w, r)
			return
		}
		// Parse the form
		err = r.ParseForm()
		if err != nil {
			log.Println(err)
			ErrorHandling(http.StatusInternalServerError, w, r)
			return
		}

		// Render the post page
		tmpl, err := template.ParseFiles("templates/detail.html")
		if err != nil {
			ErrorHandling(http.StatusInternalServerError, w, r)
			return
		}
		tmpl.Execute(w, map[string]interface{}{
			"posts":    post,
			"comments": comments,
		})
	}
}

var currentPostID int

// parsePostID parses the post ID from the URL path
func parsePostID(path string) (int, error) {
	// Split the URL path on '/'
	parts := strings.Split(path, "/")

	// Return
	if len(parts) < 3 {
		//formating error and return 0
		return 0, fmt.Errorf("invalid path")
	}
	return strconv.Atoi(parts[2])
}

// SearchPost searches for a post based on the ID
func SearchPost(id int, uniquePost *map[int]post) (map[int]post, error) {

	sqlite.OpenDb()

	defer sqlite.Db.Close()

	// Check through the database for the post
	rows, err := sqlite.Db.Query("SELECT * FROM posts WHERE post_id = ?", id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var post_id int
	var user_id int
	var title string
	var content string
	var category string
	var date string
	var author string

	// Add the post to the map
	data := make(map[int]post)
	for rows.Next() {
		err = rows.Scan(&post_id, &user_id, &title, &content, &date, &category, &author)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		data[post_id] = post{
			PostID:   post_id,
			Title:    title,
			Content:  content,
			Date:     date,
			Category: category,
			Author:   author,
			Likes:    util.CountReactionsPosts(post_id, 1),
			Dislikes: util.CountReactionsPosts(post_id, 2),
		}
	}

	return data, nil
}

func getComments(postID int) (map[int]comment, error) {
	sqlite.OpenDb()

	defer sqlite.Db.Close()

	// Check through the database for the comments for the post
	rows, err := sqlite.Db.Query("SELECT * FROM comments WHERE post_id = ?", postID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var comment_id int
	var user_id int
	var post_id int
	var content string
	var date string
	var author string

	// Add the comments to the map
	data := make(map[int]comment)
	for rows.Next() {
		err = rows.Scan(&comment_id, &user_id, &post_id, &content, &date, &author)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		data[comment_id] = comment{
			CommentID: comment_id,
			PostID:    post_id,
			UserID:    user_id,
			Content:   content,
			Date:      date,
			Author:    author,
			Likes:     util.CountReactionsComments(postID, comment_id, 1),
			Dislikes:  util.CountReactionsComments(postID, comment_id, 2),
		}
	}
	return data, nil
}
