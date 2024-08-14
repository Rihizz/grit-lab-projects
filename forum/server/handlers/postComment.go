package handlers

import (
	"fmt"
	"forum/server/login"
	"forum/server/sqlite"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

func postComment(w http.ResponseWriter, r *http.Request) {

	var (
		comment string
		userID  int
		author  string
		timeNow string
	)
	// Check the session
	if !login.CheckSession(w, r) {
		ErrorHandling(http.StatusUnauthorized, w, r)
		return
	}

	// Get the user from the session
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		ErrorHandling(http.StatusInternalServerError, w, r)
		return
	}
	userd := login.Sessions[sessionCookie.Value]

	if r.Method == "POST" {
		//parse the form
		err = r.ParseForm()
		if err != nil {
			fmt.Println("Error parsing form: ", err)
			ErrorHandling(http.StatusInternalServerError, w, r)
			return
		}

		// Get the form values
		comment = r.FormValue("comment")

		post, err := SearchPost(currentPostID, &uniquePost)
		if err != nil {
			fmt.Println("Error getting post: ", err)
			ErrorHandling(http.StatusInternalServerError, w, r)
			return
		}

		comments, err := getComments(currentPostID)
		if err != nil {
			fmt.Println("Error getting comments: ", err)
			ErrorHandling(http.StatusInternalServerError, w, r)
			return
		}

		if comment == "" {
			tmpl, err := template.ParseFiles("templates/detail.html")
			if err != nil {
				log.Println(err)
				ErrorHandling(500, w, r)
				return
			}
			tmpl.Execute(w, map[string]interface{}{
				"posts":    post,
				"comments": comments,
				"user":     userd,
				"error":    "Please Stop Trying To Break The Forum :(",
			})
			return
		}

		// counting the number of visable characters in the comment
		acount := 0
		for _, v := range comment {
			if v > 32 && v < 126 {
				acount++
			}
			if v == '\t' {
				acount = 0
				break
			}
		}
		// if the comment is empty, redirect to the post page with an error
		if acount == 0 {
			tmpl, err := template.ParseFiles("templates/detail.html")
			if err != nil {
				log.Println(err)
				ErrorHandling(500, w, r)
				return
			}
			tmpl.Execute(w, map[string]interface{}{
				"posts":    post,
				"comments": comments,
				"user":     userd,
				"error":    "Please Stop Trying To Break The Forum :(",
			})
			return
		}

		// Get the current time
		timeNow = time.Now().Format("2006-01-02 15:04:05")

		// Get the user ID from the session
		userID = login.Sessions[sessionCookie.Value].UserID
		author = login.Sessions[sessionCookie.Value].Username

		// Insert the comment into the database
		err = sqlite.AddComment(currentPostID, userID, comment, timeNow, author)
		if err != nil {
			fmt.Println("Error adding comment to database: ", err)
			ErrorHandling(http.StatusInternalServerError, w, r)
			return
		}

		// Redirect to the post page
		http.Redirect(w, r, "/post-id/"+strconv.Itoa(currentPostID), http.StatusFound)
	} else {
		// Handle other request methods
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
