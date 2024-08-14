package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"social-network/database/sqlite"
	"strconv"
)

func validateInput(title string, description string, dateTime string) error {
	// titleRegex := regexp.MustCompile(`^[a-zA-Z0-9\s.,!?;: ]{1,50}$`)
	// descriptionRegex := regexp.MustCompile(`^[a-zA-Z0-9\s.,!?;: ]{1,256}$`)
	dateTimeRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}$`) // assuming dateTime follows "YYYY-MM-DDThh:mm" format

	// if !titleRegex.MatchString(title) {
	// 	return errors.New("invalid title")
	// }
	// if !descriptionRegex.MatchString(description) {
	// 	return errors.New("invalid description")
	// }
	if !dateTimeRegex.MatchString(dateTime) {
		return errors.New("invalid dateTime")
	}

	return nil
}

type CreateEventRequest struct {
	GroupId          string `json:"group_id"`
	EventTitle       string `json:"title"`
	EventDescription string `json:"description"`
	EventDateTime    string `json:"date_time"`
}

// CreateEventAPI is the API handler for creating a new event
func CreateEventAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateEventAPI called")
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

	// parse the JSON request body into a CreateEventRequest struct
	var request CreateEventRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println("error decoding json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// validate input
	err = validateInput(request.EventTitle, request.EventDescription, request.EventDateTime)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the event details from the request
	groupId := request.GroupId
	eventTitle := request.EventTitle
	eventDescription := request.EventDescription
	eventDateTime := request.EventDateTime

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

	//print request data and user id
	// log.Println("request data:", groupId, eventTitle, eventDescription, eventDateTime)
	// log.Println("user id:", session.UserID)

	//make groupid into int
	groupIdInt, err := strconv.Atoi(groupId)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create the event in the database
	err = sqlite.CreateEvent(groupIdInt, session.UserID, eventTitle, eventDescription, eventDateTime)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return a success status code
	w.WriteHeader(http.StatusOK)
}
