package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func CheckLoginStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("login status check")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	// This is done to allow the frontend to send an OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Read the session_token cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println("login status check failed")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Println("cookie value", cookie.Value)
	log.Println("Sessions", Sessions)

	// Check if the session token exists in the Sessions map
	if session, ok := Sessions[cookie.Value]; ok {
		log.Println("login status check success")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"userID":   session.UserID,
			"nickname": session.Username,
		})
	} else {
		log.Println("login status check failed")
		w.WriteHeader(http.StatusUnauthorized)
	}
}
