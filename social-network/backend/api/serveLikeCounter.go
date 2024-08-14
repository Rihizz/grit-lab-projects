package api

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/database/sqlite"
	"strconv"
)

func ServeLikeCounter(w http.ResponseWriter, r *http.Request) {
	log.Println("ServeLikeCounter called")
	// Enable CORS for all the frontend
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	postID := r.URL.Query().Get("id")
	log.Println("postID:", postID)

	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	likes, err := sqlite.GetLikes(postIDInt)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(likes); err != nil {
		log.Println(err)
		http.Error(w, "Error in ServeLikeCounter", http.StatusInternalServerError)
		return
	}

	log.Println("ServeLikeCounter successfully finished")

}
