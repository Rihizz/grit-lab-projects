package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"social-network/database/sqlite"
)

func UnfollowAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("Unfollow API called")
	// set the headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	// check if the method is allowed
	if r.Method != "POST" && r.Method != "OPTIONS" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "{\"status\": 405, \"message\": \"method not allowed\"}")
		return
	}

	// check if the method is OPTIONS
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"status\": 200, \"message\": \"success\"}")
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

	// get the user id from the session
	userID := session.UserID

	// get the id of the user to unfollow
	var data map[string]int
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("Error decoding the request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := data["id"]
	log.Println("User", userID, "is trying to unfollow user", id)
	// check if the user is trying to unfollow himself
	if id == userID {
		log.Println("User is trying to unfollow himself")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if the user is already following the user
	isFollowing := sqlite.CheckIfFollower(userID, id)

	if !isFollowing {
		log.Println("User is not following the user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// unfollow the user
	err = sqlite.Unfollow(id, userID)
	if err != nil {
		log.Println("Error unfollowing the user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{\"status\": 200, \"message\": \"success\"}")
}
