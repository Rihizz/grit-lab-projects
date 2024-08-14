package server

import (
	"encoding/json"
	"net/http"
	"pillu/sqlite"
)

func UsersAPI(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
	if !CheckSession(w, r) {
		w.Write([]byte("false"))
		return
	}
	users, err := sqlite.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Marshal the online users slice to JSON
	resp, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}

	// Set the response header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
