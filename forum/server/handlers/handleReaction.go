package handlers

import (
	"fmt"
	"forum/server/login"
	"forum/server/sqlite"
	"net/http"
	"strconv"
)

func AddReaction(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
	if !login.CheckSession(w, r) {
		ErrorHandling(401, w, r)
		return
	}
	// Get the user from the session
	c, err := r.Cookie("session_token")
	if err != nil {
		ErrorHandling(http.StatusInternalServerError, w, r)
		return
	}
	// parse the form data
	r.ParseForm()
	// Get the current user data
	userd := login.Sessions[c.Value]
	// gettin an id for what comment or post the reaction is for. if 0, it's for a post else it's for a comment
	likeid := r.FormValue("reaction")
	for i, v := range likeid {
		if v == 'e' {
			likeid = likeid[i+1:]
			break
		}
	}
	actualLikeId, err := strconv.Atoi(likeid)
	if err != nil {
		fmt.Println("Error: ", err)
		ErrorHandling(http.StatusInternalServerError, w, r)
		return
	}
	// seeing the reaction type
	reaction := r.FormValue("reaction")
	for i, v := range reaction {
		if v == 'e' {
			reaction = reaction[:i+1]
			break
		}
	}
	if reaction == "like" {
		// checks if same reaction already exists
		if sqlite.CheckReaction(currentPostID, actualLikeId, userd.UserID, 1) {
			sqlite.DeleteReaction(currentPostID, actualLikeId, userd.UserID)
		} else {
			// function automatically updates the reaction if it already exists
			sqlite.AddReaction(currentPostID, actualLikeId, userd.UserID, 1)
		}
	} else if reaction == "dislike" {
		// checks if same reaction already exists
		if sqlite.CheckReaction(currentPostID, actualLikeId, userd.UserID, 2) {
			sqlite.DeleteReaction(currentPostID, actualLikeId, userd.UserID)
		} else {
			// function automatically updates the reaction if it already exists
			sqlite.AddReaction(currentPostID, actualLikeId, userd.UserID, 2)
		}
	}

	http.Redirect(w, r, "/post-id/"+fmt.Sprint(currentPostID), http.StatusSeeOther)
}
