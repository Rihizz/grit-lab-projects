package api

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/database/sqlite"
	"strconv"
)

func ServeFollowedUsers(w http.ResponseWriter, r *http.Request) {
	// set cors headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	// if the request method is not GET or OPTIONS, return
	if r.Method != http.MethodGet && r.Method != http.MethodOptions {
		log.Println("Method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// if the request method is OPTIONS, return
	if r.Method == http.MethodOptions {
		log.Println("Method options")
		w.WriteHeader(http.StatusOK)
		return
	}

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

	group := r.URL.Query().Get("group")
	groupID, err := strconv.Atoi(group)
	if err != nil {
		log.Println("Error converting group id to int:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("group id:", groupID)

	members, err := sqlite.GetGroupMembers(groupID)
	if err != nil {
		log.Println("Error getting group members:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID := session.UserID

	userIDstr := strconv.Itoa(userID)

	users, err := sqlite.GetAllContacts(userIDstr)
	if err != nil {
		log.Println("Error getting following for profile:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, member := range members {
		for i, user := range users {
			if member.Id == user.UserID {
				users = append(users[:i], users[i+1:]...)
			break
			}
		}
	}


	log.Println("following:", users)

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the users to the HTTP response
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Println("Error encoding users to JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}