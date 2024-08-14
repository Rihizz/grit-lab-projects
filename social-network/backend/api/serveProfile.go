package api

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"social-network/database/sqlite"
	"strconv"
	"strings"
)

// struct for the response
type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	AboutMe   string `json:"aboutme"`
	Birthday  string `json:"birthday"`
	Avatar    string `json:"avatar"`
	Privacy   string `json:"privacy"`
}

func ServeUser(w http.ResponseWriter, r *http.Request) {
	log.Println("ServeUser")

	// set the response headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	if r.Method == "OPTIONS" {
		log.Println("Method options")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		log.Println("Method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Extract userId from the request URL
	urlPath := r.URL.Path
	pathParts := strings.Split(urlPath, "/")
	userIdStr := pathParts[len(pathParts)-1]
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		log.Println("Error converting userId to int:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("userId:", userId)

	// get the user data from the database
	user, err := getUser(userId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("user:", user.AboutMe, user.Birthday, user.Email, user.FirstName, user.LastName, user.Nickname)

	// check if the active user is following the user
	// get the user id from the session
	// check if the request cookie is in the sessions map
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

	// get the user id from the session
	activeUserId := session.UserID

	// check if the active user is following the user
	isFollowing := sqlite.CheckIfFollower(activeUserId, userId)

	followers, err := sqlite.GetFollowersForProfile(userId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	following, err := sqlite.GetFollowingForProfile(userId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert followers and following to JSON strings
	followersJson, err := json.Marshal(followers)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	followingJson, err := json.Marshal(following)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isFollowing && activeUserId != userId && user.Privacy == "1" {
		log.Println("active user is not following the user")
		w.WriteHeader(http.StatusOK)
		response := struct {
			Status  int    `json:"status"`
			Message string `json:"message"`
			User    User   `json:"user"`
		}{Status: 200, Message: "success", User: User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Avatar:    user.Avatar,
			Privacy:   user.Privacy,
		}}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshaling response:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(responseJSON)
		return
	}

	// Send the response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": 200, "message": "success", "user": {"firstName": "` + user.FirstName + `", "lastName": "` + user.LastName + `", "birthday": "` + user.Birthday + `", "email": "` + user.Email + `", "nickname": "` + user.Nickname + `", "aboutme": "` + user.AboutMe + `", "avatar": "` + user.Avatar + `", "privacy": "` + user.Privacy + `", "isFollowing": ` + strconv.FormatBool(isFollowing) + `, "followers": ` + string(followersJson) + `, "following": ` + string(followingJson) + `}}`))
}

func getUser(userId int) (User, error) {
	db, err := sqlite.OpenDb()
	if err != nil {
		log.Println(err)
		return User{}, err
	}
	defer db.Close()

	var user User
	var avatar []byte

	err = db.QueryRow("SELECT firstname, lastname, birthdate, email, nickname, aboutme, avatar, privacy FROM users WHERE user_id = ?", userId).Scan(&user.FirstName, &user.LastName, &user.Birthday, &user.Email, &user.Nickname, &user.AboutMe, &avatar, &user.Privacy)
	if err != nil {
		log.Println(err)
		return user, err
	}

	user.Avatar = base64.StdEncoding.EncodeToString(avatar)

	return user, nil
}
