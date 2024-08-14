package sqlite

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// function to open db with all tables
func CreateDB() error {
	// open db
	db, err := OpenDb()
	if err != nil {
		return err
	}

	// defer the closing of the database connection
	defer db.Close()

	// create tables
	if err := createUserTbl(); err != nil {
		log.Println(err)
		return err
	}
	if err := createPostsTbl(); err != nil {
		log.Println(err)
		return err
	}
	if err := createCommentsTbl(); err != nil {
		log.Println(err)
		return err
	}
	if err := createMessagesTbl(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// fuction to open database users
func createUserTbl() error {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()
	// create table
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (user_id INTEGER PRIMARY KEY, email TEXT UNIQUE, username TEXT UNIQUE, password TEXT, age TEXT, gender TEXT, firstName TEXT, lastName TEXT)")
	if err != nil {
		log.Println(err)
		return err
	}
	if _, err := statement.Exec(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// create a new database for posts
func createPostsTbl() error {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()
	// create table
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS posts (post_id INTEGER PRIMARY KEY, user_id INTEGER, title TEXT, content TEXT, date TEXT, category TEXT, author TEXT, FOREIGN KEY(user_id) REFERENCES user(user_id))")
	if err != nil {
		log.Println(err)
		return err
	}
	if _, err := statement.Exec(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func createCommentsTbl() error {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()
	// create table
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS comments (comment_id INTEGER PRIMARY KEY, post_id INTEGER, user_id INTEGER , content TEXT, date TEXT, author TEXT, FOREIGN KEY(user_id) REFERENCES user(user_id), FOREIGN KEY(post_id) REFERENCES posts(post_id))")
	if err != nil {
		log.Println(err)
		return err
	}
	if _, err := statement.Exec(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func createMessagesTbl() error {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()
	// create table
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS messages (message_id INTEGER PRIMARY KEY, sender TEXT, receiver TEXT, text TEXT, timestamp TEXT)")
	if err != nil {
		log.Println(err)
		return err
	}
	if _, err := statement.Exec(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
