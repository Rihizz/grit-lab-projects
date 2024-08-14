package sqlite

import (
	"encoding/base64"
	"log"
)

// struct for the response
type Response struct {
	Posts []PostForResponse `json:"posts"`
}

// struct for the posts
type PostForResponse struct {
	Id        int    `json:"id"`
	UserId    int    `json:"user_id"`
	FullName  string `json:"full_name"`
	Content   string `json:"content"`
	Picture   string `json:"picture"`
	Date      string `json:"date"`
	LikeCount int    `json:"like_count"`
	Likers    []int  `json:"likers"`
}

type FollowingList struct {
	Folloiwing []int `json:"following"`
}

func GetPrivatePosts(activeUserId int) ([]PostForResponse, error) {
	db, err := OpenDb()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	followingList, err := GetFollowingList(activeUserId)
	if err != nil {
		return nil, err
	}

	// add the active user to the list
	followingList = append(followingList, activeUserId)

	var posts []PostForResponse

	for _, userId := range followingList {
		log.Println("checking user: ", userId)
		rows, err := db.Query("SELECT id, user_id, content, author, created_at, image FROM posts WHERE user_id = ? AND privacy = 1", userId)
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			var post PostForResponse
			var imageData []byte
			err := rows.Scan(&post.Id, &post.UserId, &post.Content, &post.FullName, &post.Date, &imageData)
			if err != nil {
				return nil, err
			}

			// Encode the image data to base64
			if imageData != nil {
				post.Picture = base64.StdEncoding.EncodeToString(imageData)
			}

			// get the like count
			likeCount, err := GetLikes(post.Id)
			if err != nil {
				return nil, err
			}
			post.LikeCount = likeCount

			// // get the likers
			// likers, err := getLikersList(post.Id)
			// if err != nil {
			//     return nil, err
			// }
			// post.Likers = likers
			posts = append(posts, post)
		}
	}

	// get the special posts
	specialPosts, err := getSpecialPosts(activeUserId)
	if err != nil {
		return nil, err
	}

	posts = append(posts, specialPosts...)

	return posts, nil
}

func GetFollowingList(activeUserId int) ([]int, error) {
	db, err := OpenDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query("SELECT user_id FROM followers WHERE follower_id = ?", activeUserId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var followingList []int

	for rows.Next() {
		var userId int
		err := rows.Scan(&userId)
		if err != nil {
			return nil, err
		}

		followingList = append(followingList, userId)
	}
	log.Println("followingList: ", followingList)
	return followingList, nil
}

func getSpecialPosts(activeUser int) ([]PostForResponse, error) {
	db, err := OpenDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	// get the posts
	rows, err := db.Query("SELECT id, user_id, content, author, created_at, image FROM posts WHERE privacy = 2")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []PostForResponse

	for rows.Next() {
		// check if the user has access to the post
		if err != nil {
			return nil, err
		}

		var post PostForResponse
		var imageData []byte
		err := rows.Scan(&post.Id, &post.UserId, &post.Content, &post.FullName, &post.Date, &imageData)
		if err != nil {
			return nil, err
		}

		// Encode the image data to base64
		if imageData != nil {
			post.Picture = base64.StdEncoding.EncodeToString(imageData)
		}

		canView, err := CheckViewPrivileges(activeUser, post.Id)
		if err != nil {
			return nil, err
		}
		if canView {
			posts = append(posts, post)
		}
	}

	return posts, nil
}
