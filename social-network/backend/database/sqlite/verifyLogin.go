package sqlite

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func VerifyUser(email, passwordInput string) (int, string, error) {
	// open db
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return 0, "", err
	}

	// defer the closing of the database connection
	defer db.Close()

	//check if username exists
	if exists, err := emailExists(email); !exists || err != nil {
		log.Println(err)
		return 0, "", err
	}

	//check if password matches
	rows, err := db.Query("SELECT user_id, email, fullname, password FROM users WHERE email = ?", email)
	if err != nil {
		return 0, "", err
	}
	defer rows.Close()

	var userId int
	var fullname, password string
	if rows.Next() {
		rows.Scan(&userId, &email, &fullname, &password)
		if CheckPasswordHash(passwordInput, password) {
			return userId, fullname, nil
		}
	}

	return 0, "", fmt.Errorf("username or password does not match")
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// function to check if a email exists
func emailExists(email string) (bool, error) {
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
