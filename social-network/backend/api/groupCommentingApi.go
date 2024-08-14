package api

import (
	"io"
	"log"
	"net/http"
	"social-network/database/sqlite"
	"strconv"
	Time "time"
)

func GroupCommentingAPI(w http.ResponseWriter, r *http.Request) {
	// set cors headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	// if the request method is not POST or OPTIONS, return
	if r.Method != http.MethodPost && r.Method != http.MethodOptions {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// if the request method is OPTIONS, return
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// get the user id from the session
	//check if the request cookie is in the sessions map
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println("Error getting cookie:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	session, ok := Sessions[cookie.Value]
	if !ok {
		log.Println("session not found")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Parse the multipart form data
	err = r.ParseMultipartForm(10 << 20) // 10 MB max size
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}


	postIdStr := r.FormValue("post_id")
	if postIdStr == "" {
		log.Println("post_id is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	content := r.FormValue("content")

	// validate that content is 1-500 characters long
	if len(content) < 1 || len(content) > 500 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdAt := Time.Now().Format("2006-01-02 15:04:05")
	log.Println("post_id:", postId)

	// Access the file
	file, _, err := r.FormFile("picture")
	var fileContent []byte
	var pic bool

	if err != nil {
		pic = false
	} else {
		defer file.Close()
		fileContent, err = io.ReadAll(file)
		if err != nil {
			log.Println(err)
			pic = false
		} else {
			pic = true
		}
	}

	if pic {
		err = sqlite.AddComment(postId, session.UserID, content, fileContent, createdAt, true)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		err = sqlite.AddComment(postId, session.UserID, content, nil, createdAt, true)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
