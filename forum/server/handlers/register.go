package handlers

import (
	"forum/server/login"
	"forum/server/sqlite"
	util "forum/server/utils"
	"log"
	"net/http"
	"text/template"
)

// processing registering requests, and rendering register page also
func Register(w http.ResponseWriter, r *http.Request) {

	var (
		username, email, password string
	)

	// Parse form data
	r.ParseForm()
	// Get form values
	username = r.FormValue("username")
	email = r.FormValue("email")
	password = r.FormValue("psw")

	// Handle HTTP method
	switch r.Method {
	case "GET":
		// Should not be GET, but if it is, redirect to index with error message
		if template, err := template.ParseFiles("templates/index.html"); err != nil {
			log.Println("Register handler 405")
			ErrorHandling(405, w, r)
			return
		} else {
			template.Execute(w, M{
				"cats":  util.Fcategories,
				"posts": util.Fposts,
				"msg":   "Something went wrong. Please try again.",
			})
		}

	case "POST":
		// Parse template
		template, err := template.ParseFiles("templates/index.html")
		// Check if user is already logged in
		if err != nil {
			log.Println("Register handler 500")
			ErrorHandling(500, w, r)
			return
		}
		if login.CheckSession(w, r) {
			c, err := r.Cookie("session_token")
			if err != nil {
				log.Println("Register handler 500")
				ErrorHandling(500, w, r)
				return
			}
			userd := login.Sessions[c.Value]
			template.ExecuteTemplate(w, "index.html", M{
				"msg":   "You are already logged in.",
				"cats":  util.Fcategories,
				"posts": util.Fposts,
				"user":  userd,
			})
			return
		}

		// Check if username already exists
		if exists, err := sqlite.UserExists(username); exists || err != nil {
			template.ExecuteTemplate(w, "index.html", M{
				"msg":   "Username already exists.",
				"cats":  util.Fcategories,
				"posts": util.Fposts,
			})
			return
		}

		// Check if email already exists
		if exists, err := sqlite.EmailExists(email); exists || err != nil {
			template.ExecuteTemplate(w, "index.html", M{
				"msg":   "Email already exists.",
				"cats":  util.Fcategories,
				"posts": util.Fposts,
			})
			return
		}

		// Check if password is secure
		if !PasswordSecurity([]byte(password)) {
			template.ExecuteTemplate(w, "index.html", M{
				"msg":   "Password is not secure enough, it needs to between 6 and 20 characters and contain small and big characters and numbers",
				"cats":  util.Fcategories,
				"posts": util.Fposts,
			})
			return
		}

		if !usernameValidity(username) {
			template.ExecuteTemplate(w, "index.html", M{
				"msg":   "Username is not valid, it needs to be between 3 and 10 characters",
				"cats":  util.Fcategories,
				"posts": util.Fposts,
			})
			return
		}

		if !util.ValidMailAddress(email) {
			template.ExecuteTemplate(w, "index.html", M{
				"msg":   "Email is not valid.",
				"cats":  util.Fcategories,
				"posts": util.Fposts,
			})
			return
		}

		// Hash password
		encryptedPassword, err := util.HashPassword(password)
		if err != nil {
			log.Println("Register handler 500")
			ErrorHandling(500, w, r)
			return
		}

		// Register user in database
		sqlite.RegisterUser(email, username, encryptedPassword, "user")

		// Render success message
		template.ExecuteTemplate(w, "index.html", M{
			"msg":   "registeration succesful",
			"cats":  util.Fcategories,
			"posts": util.Fposts,
		})
		return

	default:
		ErrorHandling(405, w, r)
		return
	}
}

// Makes sure the password is long enough, and contains uppercase, lowercase and numerical characters
func PasswordSecurity(password []byte) bool {
	var (
		minimumLength int = 6
		maximumLength int = 20
		cap               = false
		low               = false
		num               = false
	)

	for _, ch := range password {
		if ch >= 97 && ch <= 122 {
			low = true
		}
		if ch >= 48 && ch <= 57 {
			num = true
		}
		if ch >= 65 && ch <= 90 {
			cap = true
		}
	}

	if len(password) < minimumLength || len(password) > maximumLength {
		return false
	}

	if cap && low && num {
		return true
	} else {
		return false
	}
}

func usernameValidity(username string) bool {
	if len(username) < 3 || len(username) > 10 {
		return false
	}
	for _, ch := range username {
		if ch < 33 || ch > 126 {
			return false
		}
		continue
	}
	return true
}
