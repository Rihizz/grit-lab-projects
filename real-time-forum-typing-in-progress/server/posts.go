package server

import (
	"encoding/json"
	"net/http"
	"pillu/sqlite"
)

func PostsAPI(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
	if !CheckSession(w, r) {
		w.Write([]byte("pls login retard"))
		return
	}
	// Get posts
	posts, err := sqlite.GetPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Marshal the posts slice to JSON
	postResp, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}

	// Set the response header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(postResp)

}
