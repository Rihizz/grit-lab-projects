package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"social-network/database/sqlite"
)

type EventResponseRequest struct {
	EventId  int    `json:"eventId"`
	UserId   int    `json:"userID"`
	Response string `json:"response"`
}

func EventResponse(w http.ResponseWriter, r *http.Request) {
	log.Println("EventResponseAPI called")
	// set cors headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// if the request method is not POST or OPTIONS, return
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request EventResponseRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println("error decoding json")
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// print request data: event id, user id, response
	log.Println("request data:", request.EventId, request.UserId, request.Response)

	//add response to database
	sqlite.AddEventResponse(request.EventId, request.UserId, request.Response)

	// return a success status code
	w.WriteHeader(http.StatusOK)

}
