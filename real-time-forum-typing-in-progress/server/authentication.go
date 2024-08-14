package server

import "net/http"

func isAuthenticated(r *http.Request) bool {
	c, err := r.Cookie("session_token")
	if err != nil {
		return false
	}
	for m := range Sessions {
		if m == c.Value {
			return true
		}
	}
	return false
}

func LoggedInHandler(w http.ResponseWriter, r *http.Request) {
	if isAuthenticated(r) {
		w.Write([]byte("true"))
	} else {
		w.Write([]byte("false"))
	}
}
