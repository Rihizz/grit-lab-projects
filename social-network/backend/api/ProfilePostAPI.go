package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"social-network/database/sqlite"
	"strconv"
)

func ProfilePostAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	if r.Method != "GET" && r.Method != "OPTIONS" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "{\"status\": 405, \"message\": \"method not allowed\"}")
		return
	}

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"status\": 200, \"message\": \"success\"}")
		return
	}

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"status\": 400, \"message\": \"id is required\"}")
		return
	}

	if id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"status\": 400, \"message\": \"id is required\"}")
		return
	}

	// get the user id from the session
	// check if the request cookie is in the sessions map
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println("Error getting cookie:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	session, ok := Sessions[cookie.Value]
	if !ok {
		log.Println("session not found")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	posts, err := sqlite.GetPostsByUser(id, session.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"status\": 500, \"message\": \"internal server error\"}")
		return
	}

	// reverse the posts to show the newest posts first
	for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
		posts[i], posts[j] = posts[j], posts[i]
	}

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		log.Printf("Could not encode posts to JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"status\": 500, \"message\": \"could not encode posts to JSON\"}")
	}
}
