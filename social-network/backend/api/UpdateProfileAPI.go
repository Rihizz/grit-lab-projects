// api.go
package api

import (
	// ... other imports ...
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"social-network/database/sqlite"
)

// ... other code ...

func UpdateProfileAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("UpdateProfileAPI called")
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
	userID := r.FormValue("userId")
	email := r.FormValue("email")
	nickname := r.FormValue("nickname")
	aboutMe := r.FormValue("aboutMe")
	newPassword := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")
	privacy := r.FormValue("privacy")
	//print received data
	log.Println(userID, email, nickname, aboutMe, privacy)

	// Access uploaded file
	file, fileHeader, err := r.FormFile("avatar")
	var fileContent []byte
	var fileName string

	// Check if an avatar was uploaded
	if err == nil {
		defer file.Close()

		// Read the file's content
		fileContent, err = io.ReadAll(file)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error reading avatar", http.StatusInternalServerError)
			return
		}

		fileName = fileHeader.Filename

	}
	fmt.Println("Filename:", fileName)                    // Print the filename
	fmt.Println("File content length:", len(fileContent)) // Print the content length
	if newPassword != "" && newPassword != confirmPassword {
		http.Error(w, "Error updating user profile: Passwords do not match", http.StatusBadRequest)
		return
	}

	// Perform update logic (e.g. update user in the database)
	err = sqlite.UpdateUserProfile(userID, email, nickname, aboutMe, fileName, fileContent, newPassword, privacy)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error updating user profile", http.StatusInternalServerError)
		return
	}

	// Send a success response to the client
	response := RegisterResponse{
		Message: "User profile updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}
