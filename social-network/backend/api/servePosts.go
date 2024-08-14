package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"social-network/database/sqlite"
	"sort"
	"strconv"
)

const openDbErrorMessage = "Error opening the database, ServePosts.go:"

func ServePosts(w http.ResponseWriter, r *http.Request) {
	const internalServerErrorMessage = "{\"status\": 500, \"message\": \"internal server error\"}"

	// set the response headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	if r.Method != "GET" && r.Method != "OPTIONS" {
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
		fmt.Fprintf(w, "{\"status\": 401, \"message\": \"unauthorized\"}")
		return
	}

	postID := r.URL.Query().Get("id")
	if postID != "" {
		log.Println("request is for single post", postID)
		postIDInt, err := strconv.Atoi(postID)
		if err != nil {
			log.Println(err)
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "{\"status\": 400, \"message\": \"bad request\"}")
			return
		}
		posts, err := fetchSinglePost(postIDInt)
		if err != nil {
			log.Println(err)
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if posts.Id == 0 {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "{\"status\": 404, \"message\": \"post not found\"}")
			return
		}

		// create the response
		response := sqlite.Response{
			Posts: []sqlite.PostForResponse{posts},
		}

		// convert the response to json
		responseJSON, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
			// send a response with the error
			fmt.Fprint(w, internalServerErrorMessage)
		}

		// check if the user is following the author of the post
		following := sqlite.CheckIfFollower(session.UserID, posts.UserId)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "{\"status\": 500, \"message\": \"internal server error\"}")
		}

		if !following && posts.UserId != session.UserID && !postIsPublic(postIDInt) {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "{\"status\": 403, \"message\": \"forbidden\"}")
			return
		}

		// write the response
		w.Write(responseJSON)
		return
	}

	// get the posts
	posts, err := GetPosts(session.UserID)
	if err != nil {
		fmt.Println(err)
		// send a response with the error
		fmt.Fprint(w, internalServerErrorMessage)
	}

	// sort the posts by date
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date > posts[j].Date
	})

	// create the response
	response := sqlite.Response{
		Posts: posts,
	}

	// convert the response to json
	responseJSON, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		// send a response with the error
		fmt.Fprint(w, internalServerErrorMessage)
	}

	// write the response
	w.Write(responseJSON)
}

func GetPosts(activeUserId int) ([]sqlite.PostForResponse, error) {
	db, err := sqlite.OpenDb()
	if err != nil {
		log.Println(openDbErrorMessage, err)
		return nil, err
	}

	defer db.Close()

	// get the posts
	posts := []sqlite.PostForResponse{}

	rows, err := db.Query("SELECT id, user_id, content, author, created_at, image FROM posts WHERE privacy = 0 ORDER BY created_at DESC")
	if err != nil {
		log.Println("Error getting the posts, GetPosts(): ", err)
	}

	defer rows.Close()

	for rows.Next() {
		var post sqlite.PostForResponse
		var imageData []byte
		err := rows.Scan(&post.Id, &post.UserId, &post.Content, &post.FullName, &post.Date, &imageData)
		if err != nil {
			log.Println("Error scanning the posts, GetPosts(): ", err)
		}

		// Encode the image data to base64
		if imageData != nil {
			post.Picture = base64.StdEncoding.EncodeToString(imageData)
		}

		// get the like count
		likeCount, err := sqlite.GetLikes(post.Id)
		if err != nil {
			log.Println("Error getting the like count, GetPosts(): ", err)
		}
		post.LikeCount = likeCount

		// get the likers
		likers, err := getLikersList(post.Id)
		if err != nil {
			log.Println("Error getting the likers, GetPosts(): ", err)
		}
		post.Likers = likers

		posts = append(posts, post)
	}

	privatePosts, err := sqlite.GetPrivatePosts(activeUserId)
	log.Println("private posts: ", privatePosts)
	if err != nil {
		log.Println("Error getting the private posts, GetPosts(): ", err)
	}

	posts = append(posts, privatePosts...)

	return posts, nil
}

func fetchSinglePost(PostID int) (sqlite.PostForResponse, error) {
	db, err := sqlite.OpenDb()
	if err != nil {
		log.Println(openDbErrorMessage, err)
		return sqlite.PostForResponse{}, err
	}

	defer db.Close()

	// get the post
	post := sqlite.PostForResponse{}

	rows, err := db.Query("SELECT id, user_id, content, author, created_at, image FROM posts WHERE id = ?", PostID)

	if err != nil {
		log.Println("Error getting the posts, GetPosts(): ", err)
	}

	defer rows.Close()

	for rows.Next() {
		var imageData []byte
		err := rows.Scan(&post.Id, &post.UserId, &post.Content, &post.FullName, &post.Date, &imageData)
		if err != nil {
			log.Println("Error scanning the posts, GetPosts(): ", err)
		}

		// Encode the image data to base64
		if imageData != nil {
			post.Picture = base64.StdEncoding.EncodeToString(imageData)
			log.Println("has image")
		}

		// get the like count
		likeCount, err := sqlite.GetLikes(post.Id)
		if err != nil {
			log.Println("Error getting the like count, GetPosts(): ", err)
		}
		post.LikeCount = likeCount
	}

	return post, nil
}

func getLikersList(postID int) ([]int, error) {
	db, err := sqlite.OpenDb()
	if err != nil {
		log.Println(openDbErrorMessage, err)
		return nil, err
	}

	defer db.Close()

	// get the likers
	likers := []int{}

	rows, err := db.Query("SELECT user_id FROM reactions WHERE post_id = ?", postID)

	if err != nil {
		log.Println("Error getting the likers, GetPosts(): ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var liker int
		err := rows.Scan(&liker)
		if err != nil {
			log.Println("Error scanning the likers, GetPosts(): ", err)
			return nil, err
		}

		likers = append(likers, liker)
	}

	return likers, nil
}

func postIsPublic(postID int) bool {
	db, err := sqlite.OpenDb()
	if err != nil {
		log.Println(openDbErrorMessage, err)
		return false
	}

	defer db.Close()

	// get the post
	var privacy int

	rows, err := db.Query("SELECT privacy FROM posts WHERE id = ?", postID)

	if err != nil {
		log.Println("Error getting the posts, GetPosts(): ", err)
		return false
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&privacy)
		if err != nil {
			log.Println("Error scanning the posts, GetPosts(): ", err)
			return false
		}
	}

	return privacy == 0
}
