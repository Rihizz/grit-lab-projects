package sqlite

import (
	"encoding/base64"
	"log"
)

type Post struct {
	Id      int
	UserId  int
	Content string
	Author  string
	Date    string
	Image   []byte
}

func GetPostsByUser(userID, activeUser int) ([]PostForResponse, error) {
	db, err := OpenDb()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query("SELECT id, user_id, content, author, created_at, image FROM posts WHERE user_id = ? AND privacy = 0", userID)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer rows.Close()

	var posts []PostForResponse

	for rows.Next() {
		var post PostForResponse
		var imageData []byte
		err := rows.Scan(&post.Id, &post.UserId, &post.Content, &post.FullName, &post.Date, &imageData)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		// Encode the image data to base64
		if imageData != nil {
			post.Picture = base64.StdEncoding.EncodeToString(imageData)
		}

		posts = append(posts, post)
	}

	privatePosts, err := getPrivateProfilePosts(userID, activeUser)
	if err != nil {
		log.Print(err)
		return posts, err
	}

	posts = append(posts, privatePosts...)

	specialPosts, err := getSpecialProfilePosts(activeUser, userID)
	if err != nil {
		log.Print(err)
		return posts, err
	}

	posts = append(posts, specialPosts...)

	return posts, nil
}

func getPrivateProfilePosts(userID, activeUser int) ([]PostForResponse, error) {
	// check if the user has access to the post
	canView := CheckIfFollower(activeUser, userID)

	if !canView && activeUser != userID {
		log.Println("User does not have access to the post, getPrivateProfilePosts(): ")
		return nil, nil
	}

	db, err := OpenDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	// get the posts that have privacy status 1
	rows, err := db.Query("SELECT id, user_id, content, author, created_at, image FROM posts WHERE privacy = 1 AND user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []PostForResponse

	for rows.Next() {
		post := PostForResponse{}
		var imageData []byte
		err := rows.Scan(&post.Id, &post.UserId, &post.Content, &post.FullName, &post.Date, &imageData)
		if err != nil {
			return nil, err
		}

		// Encode the image data to base64
		if imageData != nil {
			post.Picture = base64.StdEncoding.EncodeToString(imageData)
		}

		log.Println("private post: ", post)

		posts = append(posts, post)
	}

	log.Println("private posts authorized, getPrivateProfilePosts()")
	log.Println("private posts: ", posts)

	return posts, nil
}

func getSpecialProfilePosts(activeUser, userId int) ([]PostForResponse, error) {
	db, err := OpenDb()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer db.Close()

	// get the posts that have privacy status 2 and are made by user with id userId
	rows, err := db.Query("SELECT id, user_id, content, author, created_at, image FROM posts WHERE privacy = 2 AND user_id = ?", userId)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer rows.Close()

	var posts []PostForResponse

	for rows.Next() {
		post := PostForResponse{}
		var imageData []byte
		err := rows.Scan(&post.Id, &post.UserId, &post.Content, &post.FullName, &post.Date, &imageData)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		// Encode the image data to base64
		if imageData != nil {
			post.Picture = base64.StdEncoding.EncodeToString(imageData)
		}

		canView, err := CheckViewPrivileges(activeUser, post.Id)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		if !canView && activeUser != userId {
			log.Println("User does not have access to the post, getSpecialProfilePosts(): ")
			continue
		}

		posts = append(posts, post)
	}

	return posts, nil
}
