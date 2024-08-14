package api

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/database/sqlite"
	"strconv"
)

func JoinGroupAPI(w http.ResponseWriter, r *http.Request) {
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

	// parse the JSON request body
	var request struct {
		GroupId string `json:"group_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	groupId, err := strconv.Atoi(request.GroupId)
	if err != nil {
		log.Println(err)
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

	// get the user id from the session
	userId := session.UserID

	// check if the user is already in the group
	// if the user is already in the group, return
	inGroup, err := userInGroup(userId, groupId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if inGroup {
		w.WriteHeader(http.StatusAlreadyReported)
		return
	}

	// insert the user into the group
	sqlite.AddGroupMember(userId, groupId, 0)

	// add notification to the group owner here

	w.WriteHeader(http.StatusOK)
}

func userInGroup(userID, groupID int) (bool, error) {
	db, err := sqlite.OpenDb()
	if err != nil {
		log.Println(err)
		return false, err
	}

	defer db.Close()

	// check if the user is already in the group
	rows, err := db.Query("SELECT * FROM group_members WHERE user_id = ? AND group_id = ?", userID, groupID)
	if err != nil {
		log.Println(err)
		return false, err
	}

	defer rows.Close()

	// if the user is already in the group, return
	if rows.Next() {
		return true, nil
	}

	return false, nil
}
