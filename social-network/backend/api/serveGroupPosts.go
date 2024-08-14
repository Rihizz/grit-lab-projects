package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"social-network/database/sqlite"
	"strconv"
)

// struct for the posts
type GroupPostForResponse struct {
	PostId    int    `json:"id"`
	GroupId   int    `json:"group_id"`
	UserId    int    `json:"user_id"`
	FullName  string `json:"full_name"`
	Content   string `json:"content"`
	Picture   []byte `json:"picture"`
	Date      string `json:"date"`
	LikeCount int    `json:"like_count"`
	Likers    []int  `json:"likers"`
}

func ServeGroupPosts(w http.ResponseWriter, r *http.Request) {
	log.Println("ServeGroupPosts called")
	// Enable CORS for all the frontend
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	groupID := r.URL.Query().Get("id")
	log.Println("groupID:", groupID)

	groupPostID := r.URL.Query().Get("group-postID")
	log.Println("groupPostID:", groupPostID)

	groupIDInt, err := strconv.Atoi(groupID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"status\": 400, \"message\": \"bad request\"}")
		return
	}

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

	// check if the active user is a member of the group
	isMember, err := sqlite.CheckIfGroupMember(session.UserID, groupIDInt)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("isMember:", isMember)

	if groupPostID != "" {
		groupPostIDInt, err := strconv.Atoi(groupPostID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		xd, err := sqlite.GetUserById(session.UserID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		post, err := getGroupPost(groupPostIDInt, groupIDInt, xd.FullName)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if post.PostId == 0 {
			log.Println("Post not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if !isMember {
			log.Println("User is not a member of the group")
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "{\"status\": 403, \"message\": \"forbidden\"}")
			return
		}

		actualPoster, err := sqlite.GetUserById(post.UserId)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		post.FullName = actualPoster.FullName

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(post); err != nil {
			log.Println(err)
			http.Error(w, "Error in ServeGroupPosts", http.StatusInternalServerError)
			return
		}
		return
	}

	posts, err := sqlite.GetGroupPosts(groupIDInt)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isMember {
		log.Println("User is not a member of the group")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		log.Println(err)
		http.Error(w, "Error in ServeGroupPosts", http.StatusInternalServerError)
		return
	}

	log.Println("ServeGroupPosts successfully finished")

}

func getGroupPost(groupPostID, groupId int, fullName string) (GroupPostForResponse, error) {
	db, err := sqlite.OpenDb()
	if err != nil {
		return GroupPostForResponse{}, err
	}

	defer db.Close()

	// Get the post
	post := GroupPostForResponse{}
	err = db.QueryRow("SELECT id, user_id, group_id, content, image, created_at FROM group_posts WHERE id = ?", groupPostID).Scan(&post.PostId, &post.UserId, &post.GroupId, &post.Content, &post.Picture, &post.Date)
	if err != nil {
		return GroupPostForResponse{}, err
	}

	post.FullName = fullName

	return post, nil
}
