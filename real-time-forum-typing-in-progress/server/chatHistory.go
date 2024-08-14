package server

import (
	"encoding/json"
	"net/http"
	"pillu/sqlite"
	"strconv"
)

func GetHistoryAPI(w http.ResponseWriter, r *http.Request) {
	// Get the username from the query string
	receiver := r.URL.Query().Get("receiver")
	sender := r.URL.Query().Get("sender")
	offsetString := r.URL.Query().Get("offset")

	// Check if the sender is the user that is logged in
	c, err := r.Cookie("session_token")
	if err != nil {
		w.Write([]byte("false"))
		return
	}
	userInfo, ok := Sessions[c.Value]
	if !ok {
		w.Write([]byte("false"))
		return
	}
	if sender != userInfo.Username {
		w.Write([]byte("false"))
		return
	}
	offset, err := strconv.Atoi(offsetString)
	if err != nil {
		offset = 0
	}

	// Get the messages from the database
	messages, err := sqlite.GetMessages(receiver, sender, offset)
	if err != nil {
		w.Write([]byte("false"))
		return
	}

	// Marshal the messages slice to JSON
	resp, err := json.Marshal(messages)
	if err != nil {
		w.Write([]byte("false"))
		return
	}

	// Set the response header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)

}
