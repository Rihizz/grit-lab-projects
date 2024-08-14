package api

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/database/sqlite"
	"sort"
)

type messageRequest struct {
	ReceiverID int `json:"receiver_id"`
}

func ServePrivateMessages(w http.ResponseWriter, r *http.Request) {
	// set cors headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method != "POST" && r.Method != "OPTIONS" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	c, err := r.Cookie("session_token")
	if err != nil {
		log.Println("User is not logged in")
		return
	}

	userInfo, ok := Sessions[c.Value]
	if !ok {
		log.Println("Session error user is not logged in")
		return
	}

	var messageRequest messageRequest
	err = json.NewDecoder(r.Body).Decode(&messageRequest)
	if err != nil {
		log.Println(err)
		return
	}

	messages, err := sqlite.GetPrivateMessages(userInfo.UserID, messageRequest.ReceiverID)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(messages)

	// sort messages by date newest first using the sort package
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt < messages[j].CreatedAt
	})

	//convert messages to json
	//send json to frontend
	ResponseJSON, err := json.Marshal(messages)
	if err != nil {
		log.Println(err)
		return
	}

	w.Write(ResponseJSON)
}
