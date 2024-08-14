package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/mail"
	"os"
	"regexp"
	"social-network/database/sqlite"
	"strconv"
	"time"
)

const defaultProfilePicturePath = "./database/pictures/pepega.png"

type RegisterRequest struct {
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
	Birthdate string `json:"birthday"`
	AboutMe   string `json:"aboutme"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Privacy   string `json:"privacy"`
}
type RegisterResponse struct {
	Message string `json:"message"`
}

func RegisterAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("RegisterAPI called")
	// Enable CORS for all the frontend
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Parse the multipart form data
	err := r.ParseMultipartForm(32 << 20) // Set the maximum memory to 32MB
	if err != nil {
		log.Println(err)
		http.Error(w, "Error parsing multipart form data", http.StatusBadRequest)
		return
	}

	// Access form fields
	email := r.FormValue("email")
	password := r.FormValue("password")
	nickname := r.FormValue("nickname")
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	birthday := r.FormValue("birthday")
	aboutMe := r.FormValue("aboutMe")
	fullname := firstName + " " + lastName
	privacy := r.FormValue("privacy")
	privacyInt, err := strconv.Atoi(privacy)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error converting privacy to int", http.StatusInternalServerError)
		return
	}

	// Input validation
	// if !IsValidEmail(email) {
	// 	http.Error(w, "Invalid email format", http.StatusBadRequest)
	// 	return
	// }
	// if !IsAlphaNumericOnly(nickname) || !IsAlphaNumericOnly(firstName) || !IsAlphaNumericOnly(lastName) {
	// 	http.Error(w, "Nickname, FirstName, and LastName should only consist of alphanumeric characters", http.StatusBadRequest)
	// 	return
	// }
	// if !IsValidPassword(password) {
	// 	http.Error(w, "Invalid password. It must be at least 8 characters long, contain at least one uppercase character and one special character", http.StatusBadRequest)
	// 	return
	// }
	if !IsValidDate(birthday) {
		http.Error(w, "Invalid Date, please write in DD/MM/YYYY format", http.StatusBadRequest)
		return
	}
	if !IsValidAboutMe(aboutMe) {
		http.Error(w, "AboutMe is too long. It should have a maximum length of 500 characters", http.StatusBadRequest)
		return
	}

	// Access uploaded file
	file, fileHeader, err := r.FormFile("profilePicture")
	var fileContent []byte
	var fileName string

	// Use default profile picture if no file was uploaded
	if err != nil {
		log.Println("No profile picture uploaded, using default")

		defaultFile, err := os.Open(defaultProfilePicturePath)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error opening default profile picture", http.StatusInternalServerError)
			return
		}
		defer defaultFile.Close()

		fileContent, err = io.ReadAll(defaultFile)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error reading default profile picture", http.StatusInternalServerError)
			return
		}

		fileName = "default_profile_picture.jpg"
	} else {
		defer file.Close()

		// Read the file's content
		fileContent, err = io.ReadAll(file)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error reading profile picture", http.StatusInternalServerError)
			return
		}

		fileName = fileHeader.Filename
	}

	//change birthday format to YYYY-MM-DD to match expected format in database
	// birthday = birthday[6:10] + "-" + birthday[3:5] + "-" + birthday[0:2]

	// Perform registration logic (e.g. insert user into database)
	err = sqlite.RegisterUser(email, nickname, password, birthday, aboutMe, firstName, lastName, fullname, fileName, fileContent, privacyInt)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	// Send a success response to the client
	response := RegisterResponse{
		Message: "User registered successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsAlphaNumericOnly(s string) bool {
	// Matches only alphanumeric characters
	alphaNumRegexp := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return alphaNumRegexp.MatchString(s)
}

func IsValidPassword(password string) bool {
	// At least 8 characters long, with at least one uppercase character and one special character
	//passRegexp := regexp.MustCompile(`^(?=.*[A-Z])(?=.*[^a-zA-Z\d]).{8,}$`)
	//range over string to check for uppercase and special characters
	var hasUpper, hasSpecial bool
	for _, c := range password {
		if c >= 'A' && c <= 'Z' {
			hasUpper = true
		} else if c >= '!' && c <= '/' || c >= ':' && c <= '@' || c >= '[' && c <= '`' || c >= '{' && c <= '~' {
			hasSpecial = true
		}
	}
	return len(password) >= 8 && hasUpper && hasSpecial
	//

	//return passRegexp.MatchString(password)
}

func IsValidDate(birthdate string) bool {
	// Check if birthdate is in format DD/MM/YYYY
	_, err := time.Parse("02/01/2006", birthdate)
	return err == nil
}

func IsValidAboutMe(aboutMe string) bool {
	// AboutMe should have a maximum length of 500 characters
	return len(aboutMe) <= 500
}
