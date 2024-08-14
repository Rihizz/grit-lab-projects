package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"social-network/database/sqlite"
	"strconv"
	"strings"
	"time"
)

// struct for the post
type Post struct {
	Content     string `json:"content"`
	Privacy     string `json:"privacy"`
	Picture     string `json:"picture"`
	UsersWhoSee string `json:"users_who_see"`
}

// AddPosts adds a post to the database
func ServePosting(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != "POST" && r.Method != "OPTIONS" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "{\"status\": 405, \"message\": \"method not allowed\"}")
		return
	}

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"status\": 200, \"message\": \"success\"}")
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
	log.Println("session found", session.UserID)

	// Parse the multipart form data
	err = r.ParseMultipartForm(10 << 20) // 10 MB max size
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"status\": 400, \"message\": \"bad request\"}")
		log.Println(err)
		return
	}

	// Get the form values
	content := r.FormValue("content")
	privacy := r.FormValue("privacy")
	//validate that content is 10-500 characters long
	if len(content) < 10 || len(content) > 500 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"status\": 400, \"message\": \"bad request\"}")
		return
	}

	privacyInt, err := strconv.Atoi(privacy)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"status\": 500, \"message\": \"internal server error\"}")
		return
	}

	// Access the file
	file, fileHeader, err := r.FormFile("picture")
	var fileContent []byte
	var fileName string
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
			fileName = fileHeader.Filename
			pic = true
		}
	}

	// get the user data from the database
	userId := session.UserID
	poster, err := sqlite.GetUserById(userId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !pic {
		profilePrivacy, err := sqlite.CheckProfilePrivacy(userId)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if profilePrivacy == 1 && privacy == "0" {
			// the user has a private profile and is trying to post a public post
			// post will be stored as private post even though the user tried to post it as public
			privacy = "1"
		}

		err = sqlite.AddPosts(userId, content, time.Now().Format("2006-01-02 15:04:05"), poster.FullName, privacy)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "{\"status\": 500, \"message\": \"internal server error\"}")
			return
		}
	} else {
		profilePrivacy, err := sqlite.CheckProfilePrivacy(userId)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if profilePrivacy == 1 && privacy == "0" {
			// the user has a private profile and is trying to post a public post
			// post will be stored as private post even though the user tried to post it as public
			privacy = "1"
		}
		err = sqlite.AddPosts2(userId, content, time.Now().Format("2006-01-02 15:04:05"), poster.FullName, privacy, fileName, fileContent)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "{\"status\": 500, \"message\": \"internal server error\"}")
			return
		}
	}

	if privacyInt == 2 {
		usersWhoSee := r.FormValue("users_who_see")
		log.Println("users who see:", usersWhoSee)

		// Decode the JSON array
		var usersWhoSeeArr []string
		err := json.Unmarshal([]byte(usersWhoSee), &usersWhoSeeArr)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "{\"status\": 500, \"message\": \"internal server error\"}")
			return
		}

		// Convert the array of strings to an array of integers
		var usersWhoSeeInt []int
		for _, val := range usersWhoSeeArr {
			// Trim leading and trailing quotes from the string
			trimmedVal := strings.Trim(val, "\"")

			intVal, err := strconv.Atoi(trimmedVal)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "{\"status\": 500, \"message\": \"internal server error\"}")
				return
			}
			usersWhoSeeInt = append(usersWhoSeeInt, intVal)
		}

		// get the post id of the post that was just added
		postId, err := getPostIdByUserIdAndDate(userId, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "{\"status\": 500, \"message\": \"internal server error\"}")
			return
		}

		// add the users who can see the post and the post id to the private_posts table
		err = addIdToPrivatePosts(usersWhoSeeInt, session.UserID, postId)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "{\"status\": 500, \"message\": \"internal server error\"}")
			return
		}
	}

	// Return the response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{\"status\": 200, \"message\": \"success\"}")
}

func getPostIdByUserIdAndDate(userId int, date string) (int, error) {
	db, err := sqlite.OpenDb()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	defer db.Close()

	var postId int

	err = db.QueryRow("SELECT id FROM posts WHERE user_id = ? AND created_at = ?", userId, date).Scan(&postId)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return postId, nil
}

func addIdToPrivatePosts(usersWhoSee []int, userId, postId int) error {
	db, err := sqlite.OpenDb()
	if err != nil {
		log.Println(err)
		return err
	}

	defer db.Close()

	// add active user to the list
	usersWhoSee = append(usersWhoSee, userId)

	// add the users who can see the post and the post id to the private_posts table
	for _, user := range usersWhoSee {
		_, err := db.Exec("INSERT INTO private_post (user_id, post_id) VALUES (?, ?)", user, postId)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
