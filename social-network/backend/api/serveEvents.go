package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"social-network/database/sqlite"
	"strconv"
)

type ServeEventsRequest struct {
	GroupId string `json:"groupId"`
}

func ServeEvents(w http.ResponseWriter, r *http.Request) {
	log.Println("ServeEvents called")

	// set cors headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	// if the request method is not POST or OPTIONS, return
	if r.Method != http.MethodPost && r.Method != http.MethodOptions {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// if the request method is OPTIONS, return
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// parse the JSON request body into a ServeEventsRequest struct
	var request ServeEventsRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println(err, "error parsing the request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the group id from the request
	groupId := request.GroupId
	//make id into int
	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		log.Println(err, "error converting group id to int")
		w.WriteHeader(http.StatusBadRequest)
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

	// check if the user is a member of the group
	isMember, err := sqlite.CheckIfGroupMember(userID, groupIdInt)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isMember {
		log.Println("user is not a member of the group")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "{\"status\": 403, \"message\": \"forbidden\"}")
		return
	}

	// get the events from the database
	events, err := sqlite.GetEvents(groupIdInt)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create the JSON response body
	response := struct {
		Events []sqlite.Event `json:"events"`
	}{
		Events: events,
	}

	// write the JSON response body
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return a success status code
	//w.WriteHeader(http.StatusOK)
}
