package api

import (
	"io"
	"log"
	"net/http"
	"social-network/database/sqlite"
	"strconv"
	Time "time"
)

// request body:
// {
// 	"post_id": 1,
// 	"user_id": "1",
// 	"content": "Hello World!",
// 	"image": "base64 encoded image",
// 	"created_at": "2020-01-01 00:00:00"
// }

type Comment struct {
	PostId    int    `json:"post_id"`
	UserId    string `json:"user_id"`
	Content   string `json:"content"`
	Image     string `json:"image"`
	CreatedAt string `json:"created_at"`
}

func CommentingAPI(w http.ResponseWriter, r *http.Request) {
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
	log.Println(r.Form)

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

	userId := session.UserID
	content := r.FormValue("content")
	// Access the file
	file, _, err := r.FormFile("image")
	var fileContent []byte
	var pic bool

	if err != nil {
		log.Println(err)
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
	createdAt := Time.Now().Format("2006-01-02 15:04:05")
	log.Println("post_id:", postId)

	// modify image to []byte
	// img, err := base64.StdEncoding.DecodeString(image)
	// if err != nil {
	// 	log.Println(err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	if pic {
		log.Println("image is not nil")
		err = sqlite.AddComment(postId, userId, content, fileContent, createdAt, false)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		err = sqlite.AddComment(postId, userId, content, nil, createdAt, false)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
