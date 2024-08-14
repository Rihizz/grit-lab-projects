package api

import (
	"log"
	"net/http"
)

func LogoutAPI( w http.ResponseWriter, r *http.Request) {
	log.Println("logout API")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Read the session_token cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Remove the session token from the Sessions map
	delete(Sessions, cookie.Value)

	// Remove the cookie
	cookie = &http.Cookie{
		Name:   "session_token",
		Value:  "",
		MaxAge: -1,
		Path:   "/", // the same as the cookie to be deleted
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}