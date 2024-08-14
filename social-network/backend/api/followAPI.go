package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"social-network/database/sqlite"
	"strconv"
)

type FollowRequest struct {
	Followee int    `json:"followee"`
	Follower int    `json:"follower"`
	Privacy  string `json:"privacy"`
}

func FollowAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != "POST" && r.Method != "OPTIONS" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "{\"status\": 405, \"message\": \"method not allowed\"}")
		return
	}

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"status\": 200, \"message\": \"success\"}")
		return
	}

	// get the user id from the session
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

	// Read the request body
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse the JSON data
	var followReq FollowRequest
	err = json.Unmarshal(requestBody, &followReq)
	if err != nil {
		log.Println("Error unmarshalling request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("followReq:", followReq)

	privacy, err := strconv.Atoi(followReq.Privacy)
	if err != nil {
		log.Println("Error converting privacy to int:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Insert the data into your database
	err = sqlite.AddFollower(followReq.Follower, followReq.Followee, privacy)
	if err != nil {
		log.Println("Error adding follower:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{\"status\": 200, \"message\": \"success\"}")

}
