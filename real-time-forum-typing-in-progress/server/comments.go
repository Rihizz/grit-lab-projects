package server

import (
	"encoding/json"
	"net/http"
	"pillu/sqlite"
	"strconv"
)

func CommentsAPI(w http.ResponseWriter, r *http.Request) {
	// Get the post ID from the query string
	postIDStr := r.URL.Query().Get("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Get the comments for the post
	comments, err := sqlite.GetComments(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal the comments slice to JSON
	resp, err := json.Marshal(comments)
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}

	// Set the response header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
