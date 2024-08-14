package sqlite

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDB() error {
	// open db
	db, err := OpenDb()
	if err != nil {
		return err
	}

	// defer the closing of the database connection
	defer db.Close()

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Println(err)
		return err
	}

	// Create a migrate instance using the file-based migration source
	// Use the source URI directly in migrate.NewWithDatabaseInstance()
	m, err := migrate.NewWithDatabaseInstance("file://./database/migrations", "sqlite3", driver)
	if err != nil {
		log.Println(err)
		return err
	}

	// apply all migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Println(err)
		return err
	}

	return nil
}

// 	// create tables
// 	if err := createUserTbl(); err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	if err := createPostsTbl(); err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	if err := createCommentsTbl(); err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	if err := createMessagesTbl(); err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	return nil
// }

// // fuction to open database users
// func createUserTbl() error {
// 	db, err := OpenDb()
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	defer db.Close()
// 	// create table
// 	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (user_id INTEGER PRIMARY KEY, email TEXT UNIQUE, nickname TEXT UNIQUE, password TEXT, aboutme TEXT, birthdate TEXT, firstName TEXT, lastName TEXT)")
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	if _, err := statement.Exec(); err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	return nil
// }

// // create a new database for posts
// func createPostsTbl() error {
// 	db, err := OpenDb()
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	defer db.Close()
// 	// create table
// 	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS posts (post_id INTEGER PRIMARY KEY, user_id INTEGER, title TEXT, content TEXT, date TEXT, category TEXT, author TEXT, FOREIGN KEY(user_id) REFERENCES user(user_id))")
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	if _, err := statement.Exec(); err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	return nil
// }

// func createCommentsTbl() error {
// 	db, err := OpenDb()
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	defer db.Close()
// 	// create table
// 	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS comments (comment_id INTEGER PRIMARY KEY, post_id INTEGER, user_id INTEGER , content TEXT, date TEXT, author TEXT, FOREIGN KEY(user_id) REFERENCES user(user_id), FOREIGN KEY(post_id) REFERENCES posts(post_id))")
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	if _, err := statement.Exec(); err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	return nil
// }

// func createMessagesTbl() error {
// 	db, err := OpenDb()
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	defer db.Close()
// 	// create table
// 	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS messages (message_id INTEGER PRIMARY KEY, sender_id INTEGER, receiver_id INTEGER, text TEXT, timestamp TEXT, FOREIGN KEY(sender_id) REFERENCES user(user_id), FOREIGN KEY(receiver_id) REFERENCES user(user_id))")
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	if _, err := statement.Exec(); err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	return nil
// }
