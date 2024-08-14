package api

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/database/sqlite"
	"strconv"
)

type GroupMessageRequest struct {
	GroupID string `json:"group"`
}

func ServeGroupMessages(w http.ResponseWriter, r *http.Request) {
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
	if userInfo.UserID == 0 {
		log.Println("Session error user is not logged in")
		w.WriteHeader(401)
		return
	}
	var groupMessageRequest GroupMessageRequest
	err = json.NewDecoder(r.Body).Decode(&groupMessageRequest)
	if err != nil {
		log.Println(err)
		return
	}

	// make group id an int
	groupID, err := strconv.Atoi(groupMessageRequest.GroupID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	isMember, err := sqlite.CheckIfGroupMember(userInfo.UserID, groupID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	if !isMember {
		log.Println("User is not a member of this group")
		w.WriteHeader(403)
		return
	}

	log.Println(groupMessageRequest.GroupID)
	messages, err := sqlite.GetGroupMessages(groupID)
	if err != nil {
		log.Println(err)
		return
	}
	// log.Println(messages)
	// Convert messages to a map
	responseMap := make(map[string]interface{})
	responseMap["messages"] = messages
	//convert messages to json
	//send json to frontend as an array of messages
	ResponseJSON, err := json.Marshal(responseMap)
	if err != nil {
		log.Println(err)
		return
	}

	w.Write(ResponseJSON)
}
