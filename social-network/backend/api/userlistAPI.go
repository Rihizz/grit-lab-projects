package api

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/database/sqlite"
)

func UserListAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("UserListAPI called")
	// Enable CORS for all the frontend
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	//print userID from the form data
	userID := r.FormValue("userID")
	log.Println(userID)
	log.Println("UserListAPI userID HERE YAY")
	userlist, err := sqlite.GetAllContacts(userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error in UserListAPI", http.StatusInternalServerError)
		return
	}

	log.Println("Userlist: ", userlist)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userlist); err != nil {
		log.Println(err)
		http.Error(w, "Error in UserListAPI", http.StatusInternalServerError)
		return
	}

	log.Println("UserListAPI successfully finished")

}
