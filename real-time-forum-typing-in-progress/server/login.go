package server

import (
	"log"
	"net/http"
	"pillu/sqlite"
	"text/template"
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

func Signin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	// get username and password from form
	username := r.FormValue("username")
	password := r.FormValue("password")
	userID, err := sqlite.VerifyLogin(username, password)
	if err != nil {
		log.Println("Error verifying login:", err)
		tmpl.Execute(w, M{
			"error": "Error verifying login",
		})
		return
	}
	if userID == -1 {
		log.Println("Invalid username or password")
		tmpl.Execute(w, M{
			"error": "Invalid username or password",
		})
		return
	}

	// if login credentials are valid, create a new session
	SessionToken := randomSessionToken()
	expiresAt := time.Now().Add(1 * time.Hour)
	// set the cookie on the response
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   SessionToken,
		Expires: expiresAt,
		Path:    "/",
	})

	for key, value := range Sessions {
		if value.UserID == userID {
			delete(Sessions, key)
			break
		}
	}

	user, err := sqlite.GetUsernameByID(userID)
	if err != nil {
		log.Println("Error getting username:", err)
		tmpl.Execute(w, M{
			"error": "Error getting username",
		})
		return
	}

	Sessions[SessionToken] = Session{userID, user, expiresAt}

	log.Println("sessions", Sessions)
}

// randomise the session token
func randomSessionToken() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}
