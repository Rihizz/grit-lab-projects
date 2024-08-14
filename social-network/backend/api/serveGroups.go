package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social-network/database/sqlite"
)

func ServeGroups(w http.ResponseWriter, r *http.Request) {
	//set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	// if the request method is not GET or OPTIONS, return
	if r.Method != http.MethodGet && r.Method != http.MethodOptions {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// if the request method is OPTIONS, return
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// get the groups from the database
	groups, err := sqlite.GetGroups()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// write te groups to the response with group:groups
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(struct {
		Groups []sqlite.Group `json:"groups"`
	}{Groups: groups})
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

}
