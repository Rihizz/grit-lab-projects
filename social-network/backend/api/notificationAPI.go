package api

import (
	"encoding/json"
	"net/http"
	"social-network/database/sqlite"
	"strconv"
)

func NotificationAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	userIDstr := r.URL.Query().Get("id")

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	notifications, err := sqlite.GetNotifications(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(notifications)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
