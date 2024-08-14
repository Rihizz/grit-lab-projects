package server

import "net/http"

func CheckSession(w http.ResponseWriter, r *http.Request) bool {
	// Check if session exists
	if len(Sessions) > 0 {
		cookies := r.Cookies()
		for _, c := range cookies {
			if c.Name == "session_token" {
				for m := range Sessions {
					if m == c.Value {
						return true
					}
				}
			}
		}
	}
	return false
}
