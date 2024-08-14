package api

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/database/sqlite"
)

type SearchUserRequest struct {
	SearchQuery string `json:"searchQuery"`
}

func SearchUserAPI(w http.ResponseWriter, r *http.Request) {
	// Enable CORS for all the frontend
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	// if the request method is not GET or OPTIONS, return
	if r.Method != http.MethodPost && r.Method != http.MethodOptions {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// if the request method is OPTIONS, return
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// unmarshal the request body
	var searchUserRequest SearchUserRequest
	err := json.NewDecoder(r.Body).Decode(&searchUserRequest)
	if err != nil {
		log.Println(err, "error parsing the request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the search query from the request
	searchQuery := searchUserRequest.SearchQuery

	// get the users from the database
	users, err := sqlite.SearchUser(searchQuery)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("users:", users)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Println(err)
		http.Error(w, "Error in SearchUserAPI", http.StatusInternalServerError)
		return
	}
}
