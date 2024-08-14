package api

import (
	"encoding/json"
	"net/http"
	"social-network/database/sqlite"
	"time"

	"github.com/gofrs/uuid"
)

// map to store sessions
var Sessions = map[string]Session{}

// struct to store session data
type Session struct {
	UserID   int
	Username string
	expiry   time.Time
}

var SessionToken string

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
}

func LoginAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Parse the JSON request body into a RegisterRequest struct
	var request RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	email := request.Email
	password := request.Password
	userID, nickname, err := sqlite.VerifyUser(email, password)
	// if there is an make the response message "Invalid credentials" and send it to the client
	if err != nil {
		response := LoginResponse{
			Message: "Invalid credentials",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	if userID == 0 {
		response := LoginResponse{
			Message: "Invalid credentials",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// create a new session token if credentials are valid
	sessionToken := randomSessionToken()
	expiresAt := time.Now().Add(24 * time.Hour)

	// set the cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
		Path:    "/",
	})

	Sessions[sessionToken] = Session{userID, nickname, expiresAt}

	// Send a success response to the client
	response := LoginResponse{
		Message: "User logged in successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// randomise the session token
func randomSessionToken() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}
