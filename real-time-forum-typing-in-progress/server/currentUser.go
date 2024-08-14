package server

import (
	"encoding/json"
	"net/http"
)

func CurrentUserAPI(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		w.Write([]byte("false"))
		return
	}
	sessionToken := c.Value
	if session, ok := Sessions[sessionToken]; ok {
		// Marshal the session.Username to JSON
		Username, err := json.Marshal(session.Username)
		if err != nil {
			http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(Username)
		return
	}
	w.Write([]byte("false"))
}
