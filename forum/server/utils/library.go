package util

import (
	"database/sql"
	"forum/server/sqlite"
	"log"
	"math/rand"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// random string generator
func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

var Fposts = map[int]Fpost{}

type Fpost struct {
	PostID   int
	UserID   int
	Title    string
	Content  string
	Date     string
	Category string
	Author   string
	Count    int
}

var Fcategories = map[int]Categories{}

type Categories struct {
	Category string
	Count    int
}

// function to get all posts
func GetPosts() {
	// open db
	db, err := sqlite.OpenDb()
	if err != nil {
		return
	}

	//get all posts
	rows, err := db.Query("SELECT post_id, user_id, title, content, date, category, author FROM posts")
	if err != nil {
		log.Println(err)
		return
	}
	var post_id int
	var user_id int
	var title string
	var content string
	var date string
	var category string
	var author string

	for rows.Next() {
		rows.Scan(&post_id, &user_id, &title, &content, &date, &category, &author)
		Fposts[post_id] = Fpost{post_id, user_id, title, content, date, category, author, CountCommentsInPost(post_id)}
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// function is post liked by user
func IsLiked(post_id int, user_id int) bool {
	// open db
	db, err := sql.Open("sqlite3", "database/allData.sqlite3")
	if err != nil {
		log.Print(err)
		return false
	}
	//get all posts
	rows, err := db.Query("SELECT post_id, user_id FROM likes")
	if err != nil {
		log.Println(err)
		return false
	}

	var post_id_db int
	var user_id_db int
	for rows.Next() {
		rows.Scan(&post_id_db, &user_id_db)
		if post_id_db == post_id && user_id_db == user_id {
			return true
		}
	}
	return false
}

// func get category
func GetCategories() {
	// open db
	_, err := sqlite.OpenDb()
	if err != nil {
		log.Print(err)
	}

	defer sqlite.Db.Close()

	//get all posts
	rows, err := sqlite.Db.Query("SELECT category_id, category FROM categories")
	if err != nil {
		log.Println(err)
		return
	}
	var category_id_db int
	var category string
	for rows.Next() {
		rows.Scan(&category_id_db, &category)
		Fcategories[category_id_db] = Categories{category, CountPostsInCat(category)}
	}
}

// function to check if email is valid
func ValidMailAddress(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}

// func to count post in category
func CountPostsInCat(cat string) int {
	db, err := sqlite.OpenDb()
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	// get the amount of posts in category
	count := 0
	for _, v1 := range Fposts {
		catSlice := strings.Split(v1.Category, " ")
		for _, v2 := range catSlice {
			if v2 == cat {
				count++
			}
		}
	}
	return count
}

// func to count comment in post
func CountCommentsInPost(postID int) int {
	db, err := sqlite.OpenDb()

	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	// get the amount of comments in post
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", postID).Scan(&count)
	if err != nil {
		log.Println(err)
	}

	return count
}
