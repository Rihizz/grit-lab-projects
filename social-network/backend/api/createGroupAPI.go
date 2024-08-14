package api

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/database/sqlite"
)

type CreateGroupRequest struct {
	GroupName        string `json:"name"`
	GroupDescription string `json:"description"`
}

// CreateGroupAPI is the API handler for creating a new group
func CreateGroupAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateGroupAPI called")
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

	// parse the JSON request body into a CreateGroupRequest struct
	var request CreateGroupRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the group name and description from the request
	groupName := request.GroupName
	groupDescription := request.GroupDescription
	// check if the group name is at least 1 character long and at most 50 characters long
	//check if the description is at least 10 character long and at most 500 characters long
	if len(groupName) < 1 || len(groupName) > 50 || len(groupDescription) < 10 || len(groupDescription) > 500 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
	log.Println("session found", session.UserID)

	// get the user data from the database
	userId := session.UserID

	// create the group in the database
	err = sqlite.CreateGroup(groupName, groupDescription, userId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get the group id from the database
	groupId, err := getGroupId(groupName)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// add the creator of the group to the group
	err = sqlite.AddGroupMember(userId, groupId, 1)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return a success status code
	w.WriteHeader(http.StatusOK)
}

func getGroupId(groupName string) (int, error) {
	db, err := sqlite.OpenDb()
	if err != nil {
		return 0, err
	}

	defer db.Close()

	var groupId int
	err = db.QueryRow("SELECT id FROM groups WHERE title = ?", groupName).Scan(&groupId)
	if err != nil {
		return 0, err
	}

	return groupId, nil
}
