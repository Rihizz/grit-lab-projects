package handlers

import (
	"forum/server/sqlite"
	util "forum/server/utils"
	"net/http"
	"strconv"
)

// Handles the requests.
func Handler(w http.ResponseWriter, r *http.Request) {
	// Get the posts from the database
	util.GetPosts()
	util.GetCategories()
	// Parse the post ID from the URL
	id, _ := parsePostID(r.URL.Path)

	// Search for the post with the given id
	uniquePost, err := SearchPost(id, &uniquePost)
	if err != nil {
		ErrorHandling(404, w, r)
		return
	}

	// Check if the URL path is allowed
	if allowedURLs(r.URL.Path, uniquePost) {
		viewPost(w, r)
		return
	}

	// get category from the database
	categorySlice := sqlite.GetCategories()

	// Check if the URL path with catName endpoint is allowed
	if allowedBoards(r.URL.Path, categorySlice) {
		viewCertainCat(w, r)
		return
	}

	switch r.URL.Path {
	case "/":
		Home(w, r)
		return
	case "/login/":
		Login(w, r)
		return
	case "/loginhandle/":
		Login(w, r)
		return
	case "/register/":
		Register(w, r)
		return
	case "/favicon.ico":
		return
	case "/logout/":
		Logout(w, r)
		return
	case "/createpost/":
		CreatePost(w, r)
		return
	case "/postpost/":
		CreatePost(w, r)
		return
	case "/postcomment/":
		postComment(w, r)
		return
	case "/account/":
		viewAccount(w, r)
		return
	case "/reactionhandler/":
		AddReaction(w, r)
		return
	default:
		ErrorHandling(404, w, r)
		return
	}
}

func allowedURLs(s string, data map[int]post) bool {
	// loop through the map and check if the postID in url is in database
	for _, v := range data {
		value := strconv.Itoa(v.PostID)
		if s == "/post-id/"+value {
			return true
		}
	}
	return false
}

func allowedBoards(s string, categories []string) bool {
	// loop through the map and check if the category/board in url is in database
	for _, v := range categories {
		if s == "/"+v {
			return true
		}
	}
	return false
}
