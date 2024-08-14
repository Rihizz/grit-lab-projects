package sqlite

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// function to compare login credentials
func VerifyLogin(loginIdentifier, passwordInput string) (int, error) {
	// open db
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return -1, err
	}

	// defer the closing of the database connection
	defer db.Close()

	//check if username or email exists
	if exists, err := UserExists(loginIdentifier); !exists || err != nil {
		log.Println(err)
		return -1, err
	}

	//check if password matches
	rows, err := db.Query("SELECT user_id, email, username, password FROM users WHERE username = ? OR email = ?", loginIdentifier, loginIdentifier)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	var user_id int
	var email string
	var username string
	var password string
	if rows.Next() {
		rows.Scan(&user_id, &email, &username, &password)
		if CheckPasswordHash(passwordInput, password) {
			return user_id, nil
		}
	}

	return -1, fmt.Errorf("username or password does not match")
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// function to check if a user exists
func UserExists(loginIdentifier string) (bool, error) {
	// open db
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return false, err
	}

	// defer the closing of the database connection
	defer db.Close()

	//check if user exists by username or email
	rows, err := db.Query("SELECT user_id FROM users WHERE username = ? OR email = ?", loginIdentifier, loginIdentifier)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

// function to check if a email exists
func EmailExists(email string) (bool, error) {
	// open db
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return false, err
	}

	// defer the closing of the database connection
	defer db.Close()

	//check if email exists
	rows, err := db.Query("SELECT user_id FROM users WHERE email = ?", email)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}
