package sqlite

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(email, nickname, password, age, aboutme, firstName, lastName, fullname, fileName string, avatar []byte, privacy int) error {
	// open the database connection
	db, err := OpenDb()
	if err != nil {
		return err
	}

	// defer the closing of the database connection
	defer db.Close()

	// start a new transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	password, err = hashPassword(password)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	// insert user into table
	statement, err := tx.Prepare("INSERT INTO users (email, nickname, password, birthdate, aboutme, firstName, lastName, fullname, avatarname, avatar, privacy) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	// execute the prepared statement and insert the new user
	result, err := statement.Exec(email, nickname, password, age, aboutme, firstName, lastName, fullname, fileName, avatar, privacy)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	// check if the insert operation succeeded
	lastInsertID, err := result.LastInsertId()
	if err != nil || lastInsertID == 0 {
		log.Println(err)
		tx.Rollback()
		return err
	}

	// commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
