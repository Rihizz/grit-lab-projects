package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pillu/sqlite"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type M map[string]interface{}

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	users, err := sqlite.GetUsers()
	if err != nil {
		log.Println("Home - 500. Can't get users.")
		return
	}

	if r.URL.Path == "/" && r.Method == "GET" {
		if CheckSession(w, r) {
			c, err := r.Cookie("session_token")
			if err != nil {
				log.Println("Home - 500. Can't get cookie.")
				return
			}
			sessionToken := c.Value
			session := Sessions[sessionToken]
			tmpl.Execute(w, M{
				"session":  session,
				"sessions": Sessions,
				"users":    users,
			})
			return
		} else {
			tmpl.Execute(w, M{
				"sessions": Sessions,
				"users":    users,
			})
			return
		}
	}
	if r.URL.Path == "/login" && r.Method == "POST" {
		login(w, r)
	}
	if r.URL.Path == "/register" && r.Method == "POST" {
		register(w, r)
	}
	if r.URL.Path == "/logout" && r.Method == "GET" {
		logout(w, r)
	}
	if r.URL.Path == "/createPost" && r.Method == "POST" {
		createPost(w, r)
	}
	if r.URL.Path == "/cum" && r.Method == "POST" {
		cum(w, r)
	}
}
func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userData := map[string]string{
		"username": r.FormValue("username"),
		"password": r.FormValue("password"),
	}

	log.Println(userData["username"], " is trying to login")
	Signin(w, r)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func register(w http.ResponseWriter, r *http.Request) {
	var userData map[string]string
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Sanitize input
	userData["username"] = sanitizeInput(userData["username"])
	userData["firstName"] = sanitizeInput(userData["firstName"])
	userData["lastName"] = sanitizeInput(userData["lastName"])
	userData["email"] = sanitizeInput(userData["email"])
	userData["password"] = sanitizeInput(userData["password"])

	if !isValidInput(userData) {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if username is already taken
	_, err = sqlite.GetUserIdByUsername(userData["username"])
	if err == nil {
		http.Error(w, "Username is already taken", http.StatusBadRequest)
		return
	}

	// Check if email is already taken
	_, err = sqlite.GetUserIdByEmail(userData["email"])
	if err != nil && err != sqlite.ErrEmailNotFound {
		http.Error(w, "Email is already taken", http.StatusBadRequest)
		return
	}

	err = sqlite.RegisterUser(userData["email"], userData["username"], userData["password"], userData["age"], userData["gender"], userData["firstName"], userData["lastName"])
	if err != nil {
		// Check if the error is related to the unique constraint on email
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			http.Error(w, "Email is already taken", http.StatusBadRequest)
			return
		} else {
			log.Println("Register - 500. Can't register user.", err)
			http.Error(w, "Can't register user", http.StatusInternalServerError)
			return
		}
	}
	BroadcastNewUser(userData["username"])
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		log.Println("Logout - 500. Can't get cookie.")
		return
	}
	sessionToken := c.Value
	delete(Sessions, sessionToken)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("CreatePost - 500. Can't parse form.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Check for "<" and ">" characters in title, content, and category fields
	for _, field := range []string{"title", "content", "category"} {
		fieldValue := strings.TrimSpace(r.FormValue(field))
		if fieldValue == "" || strings.HasPrefix(fieldValue, " ") || strings.HasSuffix(fieldValue, " ") || strings.Contains(fieldValue, "\n") || strings.Contains(fieldValue, "<") || strings.Contains(fieldValue, ">") {
			errMsg := "Don't use newlines or spaces or > < in the " + field + " 8=====D"
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}
	}
	postData := map[string]string{
		"content":  r.FormValue("content"),
		"title":    r.FormValue("title"),
		"category": r.FormValue("category"),
		"date":     string(time.Now().Format("2006-01-02 15:04:05")),
	}
	c, err := r.Cookie("session_token")
	if err != nil {
		log.Println("Not logged in.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	author := Sessions[c.Value].Username
	user_id := Sessions[c.Value].UserID
	sqlite.AddPost(postData["title"], postData["content"], postData["date"], user_id, author, postData["category"])
}

func cum(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cumData := map[string]string{
		"post_id": r.FormValue("post_id"),
		"content": r.FormValue("content"),
	}
	// Check for disallowed characters in the input values
	disallowed := []string{"<", ">", "\n"}
	for field, value := range cumData {
		for _, d := range disallowed {
			if strings.Contains(value, d) {
				http.Error(w, fmt.Sprintf("Invalid %s", field), http.StatusBadRequest)
				return
			}
		}
		// Check if the field starts or ends with a space or newline
		if strings.HasPrefix(value, " ") || strings.HasPrefix(value, "\n") ||
			strings.HasSuffix(value, " ") || strings.HasSuffix(value, "\n") {
			http.Error(w, fmt.Sprintf("Invalid %s", field), http.StatusBadRequest)
			return
		}
	}
	// check that the content is not empty or only whitespace
	if strings.TrimSpace(cumData["content"]) == "" {
		http.Error(w, "Empty content", http.StatusBadRequest)
		return
	}

	cumtent := cumData["content"]

	if len(cumtent) > 140 {
		http.Error(w, "Cum too long", http.StatusBadRequest)
		return
	}

	c, err := r.Cookie("session_token")
	if err != nil {
		log.Println("Not logged in.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	log.Println(c.Value)
	log.Println(Sessions[c.Value].Username)
	author := Sessions[c.Value].Username
	log.Println("form inc")
	log.Println(r.Form)
	userID, err := sqlite.GetUserIdByUsername(author)
	if err != nil {
		log.Println("Cum - 500. Can't get user_id.")
		return
	}
	postID, err := strconv.Atoi(cumData["post_id"])
	if err != nil {
		log.Println("Cum - 500. Can't convert post_id to int.")
		return
	}
	sqlite.AddComment(postID, userID, cumData["content"], string(time.Now().Format("2006-01-02 15:04:05")), author)
}

func sanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\n", "")
	return input
}

func isValidInput(data map[string]string) bool {
	// Define a slice of fields to validate
	fields := []string{"username", "firstName", "lastName", "email", "password"}
	for _, field := range fields {
		input := data[field]
		// Check that input is not empty
		if input == "" {
			return false
		}
		// Check that input does not contain "<" or ">"
		if strings.Contains(input, "<") || strings.Contains(input, ">") {
			return false
		}
		// Check that input does not start or end with a space or newline
		if strings.HasPrefix(input, " ") || strings.HasPrefix(input, "\n") || strings.HasSuffix(input, " ") || strings.HasSuffix(input, "\n") {
			return false
		}
	}
	return true
}
