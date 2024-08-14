package server

import (
	"encoding/json"
	"net/http"
)

func OnlineUsersAPI(w http.ResponseWriter, r *http.Request) {
	// Create a slice of usernames for online users
	var onlineUsers []string
	for _, conn := range Connections {
		onlineUsers = append(onlineUsers, conn.Username)
	}

	// Marshal the online users slice to JSON
	resp, err := json.Marshal(onlineUsers)
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}

	// Set the response header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
