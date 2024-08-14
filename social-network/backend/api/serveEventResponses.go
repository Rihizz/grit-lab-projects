package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"social-network/database/sqlite"
)

type ServeEventResponsesRequest struct {
	EventId int `json:"EventId"`
}

func ServeEventResponses(w http.ResponseWriter, r *http.Request) {
	log.Println("ServeEventResponses called")

	// set cors headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	// if the request method is OPTIONS, return
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		fmt.Println("Options")
		return
	}
	// if the request method is not POST or OPTIONS, return
	if r.Method != http.MethodPost && r.Method != http.MethodOptions {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println("Not post / options")
		return
	}

	//get the event id from the request body
	var request ServeEventResponsesRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println("error decoding json")
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("request", request)
	fmt.Println("request", request.EventId)

	//get the responses from the database
	responses, err := sqlite.GetEventResponses(request.EventId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//fmt.Println("responses", responses)
	//call GetUserById for each response and print the full names with the responses
	type EventResponseWithFullName struct {
		UserId   int    `json:"user_id"`
		Response string `json:"response"`
		FullName string `json:"full_name"`
	}
	var responsesWithFullName []EventResponseWithFullName
	for _, response := range responses {
		user, err := sqlite.GetUserById(response.UserId)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		responseWithFullName := EventResponseWithFullName{
			UserId:   response.UserId,
			Response: response.Response,
			FullName: user.FullName,
		}
		responsesWithFullName = append(responsesWithFullName, responseWithFullName)
	}
	fmt.Println("responsesWithFullName", responsesWithFullName)

	//write the responses to the response
	err = json.NewEncoder(w).Encode(responsesWithFullName)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
