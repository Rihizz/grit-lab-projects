package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"social-network/database/sqlite"
	"strconv"
)

// ServeGroups is the handler for the /groups endpoint

func ServeGroupComments(w http.ResponseWriter, r *http.Request) {
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


	// get the post id from the request
	PostID := r.URL.Query().Get("id")
	PostIDInt, err := strconv.Atoi(PostID)
	if err != nil {
		log.Println(err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	groupId := r.URL.Query().Get("groupId")
	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		log.Println(err)
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the posts from the database
	comments, err := sqlite.GetComments(PostIDInt, true)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	isGroupMember, err := sqlite.CheckIfGroupMember(session.UserID, groupIdInt)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isGroupMember {
		log.Println("User is not a member of the group")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// write the posts to the response
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(struct {
		Comments []sqlite.CommentForResponse `json:"comments"`
	}{Comments: comments})
	if err != nil {
		log.Println("Error encoding JSON:", err)
		fmt.Println("Error encoding JSON:", err)
		return
	}

}
