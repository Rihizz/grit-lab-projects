// Handles the likes and unlikes of posts
// 0 = no reaction
// 1 = has liked and unliked
// 2 = active like

package api

import (
	"fmt"
	"log"
	"net/http"
	"social-network/database/sqlite"
	"strconv"
)

func HandlingLikes(w http.ResponseWriter, r *http.Request) {
	// set cors headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != "POST" && r.Method != "OPTIONS" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "{\"status\": 405, \"message\": \"method not allowed\"}")
		return
	}

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	//check if the request cookie is in the sessions map
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
	log.Println("session found", session.UserID)

	// get post id from url
	postId := r.URL.Query().Get("id")
	if postId == "" {
		fmt.Println("post_id is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	postIDInt, err := strconv.Atoi(postId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if the user has already liked the post
	liked, err := sqlite.HasLiked(session.UserID, postIDInt)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("liked:", liked)

	// like or dislike the post
	err = sqlite.AddLike(session.UserID, postIDInt)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return success
	w.WriteHeader(http.StatusOK)

	log.Println("handling likes success")
}
