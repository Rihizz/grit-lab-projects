package api

import (
	"encoding/json"
	"net/http"
	"social-network/database/sqlite"
	"strconv"
)

type Followers struct {
	Followers []sqlite.User `json:"followers"`
}

func ServeFollowerList(w http.ResponseWriter, r *http.Request) {
	// set the response headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	if r.Method != "GET" && r.Method != "OPTIONS" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// get the user id from the url
	userId := r.URL.Query().Get("id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	followersList, err := sqlite.GetFollowers(userIdInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create the response
	response := Followers{
		Followers: followersList,
	}

	// send the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
