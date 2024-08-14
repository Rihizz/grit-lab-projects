package login

import (
	"fmt"
	"forum/server/sqlite"
	"net/http"
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
	debug := false
	// get username and password from form
	username := r.FormValue("username")
	password := r.FormValue("psw")
	userID, err := sqlite.VerifyLogin(username, password)

	if debug {
		fmt.Println("username:", username)
		fmt.Println("password:", password)
		fmt.Println("userID:", userID)
		fmt.Println("error:", err)
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

	Sessions[SessionToken] = Session{userID, username, expiresAt}

	if debug {
		//print the session map
		fmt.Println("Sessions:", Sessions)
		fmt.Println("==============================")
	}

	if debug {
		fmt.Println("sessionToken:", SessionToken)
		fmt.Println("expiresAt:", expiresAt)
		fmt.Println("==============================")
	}
}

func CheckSession(w http.ResponseWriter, r *http.Request) bool {
	// Check if session exists
	if len(Sessions) > 0 {
		cookies := r.Cookies()
		for _, c := range cookies {
			if c.Name == "session_token" {
				for m := range Sessions {
					if m == c.Value {
						return true
					}
				}
			}
		}
	}
	return false
}

// randomise the session token
func randomSessionToken() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}
