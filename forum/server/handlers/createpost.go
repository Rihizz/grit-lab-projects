package handlers

import (
	"forum/server/login"
	"forum/server/sqlite"
	util "forum/server/utils"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
)

var AllCats = map[int]util.Categories{}

type FormData struct {
	Title      string
	Content    string
	Categories []string
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var (
		userID int
		date   string
		author string
	)

	if r.URL.Path == "/createpost/" {
		// Check if the user is logged in
		if !login.CheckSession(w, r) {
			ErrorHandling(401, w, r)
			return
		} else {
			// rendering the create post page
			tmpl, err := template.ParseFiles("templates/createpost.html")
			if err != nil {
				// handle error
				return
			}
			tmpl.Execute(w, map[string]interface{}{
				"categories": util.Fcategories,
			})
		}
	}

	if r.URL.Path == "/postpost/" {
		// Check if the user is logged in
		if !login.CheckSession(w, r) {
			log.Println("Not logged in, redirecting to /")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		// parse the form data
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			ErrorHandling(500, w, r)
			return
		}

		// get the form values
		title := r.Form.Get("title")
		content := r.Form.Get("content")

		for _, v := range title {
			if v == '<' || v == '>' {
				tmpl, err := template.ParseFiles("templates/createpost.html")
				if err != nil {
					log.Println(err)
					ErrorHandling(500, w, r)
					return
				}
				tmpl.Execute(w, map[string]interface{}{
					"categories": util.Fcategories,
					"error":      "Please don't use < or > in the title",
				})
				return
			}
		}

		if title == "" || content == "" {
			tmpl, err := template.ParseFiles("templates/createpost.html")
			if err != nil {
				log.Println(err)
				ErrorHandling(500, w, r)
				return
			}
			tmpl.Execute(w, map[string]interface{}{
				"categories": util.Fcategories,
				"error":      "Please fill out all fields",
			})
			return
		}

		for _, v := range content {
			if v == '<' || v == '>' {
				tmpl, err := template.ParseFiles("templates/createpost.html")
				if err != nil {
					log.Println(err)
					ErrorHandling(500, w, r)
					return
				}
				tmpl.Execute(w, map[string]interface{}{
					"categories": util.Fcategories,
					"error":      "Please don't use < or > in the content",
				})
				return
			}
		}

		countForTitle := 0
		for _, v := range title {
			if v > 32 && v < 126 {
				countForTitle++
			}
			if v == '\t' {
				countForTitle = 0
				break
			}
		}
		if countForTitle == 0 {
			tmpl, err := template.ParseFiles("templates/createpost.html")
			if err != nil {
				log.Println(err)
				ErrorHandling(500, w, r)
				return
			}
			tmpl.Execute(w, map[string]interface{}{
				"categories": util.Fcategories,
				"error":      "KYS RETARD STOP TRYING TO BREAK OUR FORUM :(",
			})
			return
		}

		acount := 0
		for _, v := range content {
			if v > 32 && v < 126 {
				acount++
			}
			if v == '\t' {
				acount = 0
				break
			}
		}
		if acount == 0 {
			tmpl, err := template.ParseFiles("templates/createpost.html")
			if err != nil {
				log.Println(err)
				ErrorHandling(500, w, r)
				return
			}
			tmpl.Execute(w, map[string]interface{}{
				"categories": util.Fcategories,
				"error":      "KYS RETARD STOP TRYING TO BREAK OUR FORUM :(",
			})
			return
		}

		// get the categories
		var categories []string
		for key, values := range r.Form {
			if key == "category" {
				categories = append(categories, values...)
			}
		}

		// check if the user selected at least one category
		if len(categories) == 0 {
			tmpl, err := template.ParseFiles("templates/createpost.html")
			if err != nil {
				log.Println(err)
				ErrorHandling(500, w, r)
				return
			}
			tmpl.Execute(w, map[string]interface{}{
				"categories": util.Fcategories,
				"error":      "Please select at least one category",
			})
			return
		}

		// create a struct with the form data
		formData := FormData{
			Title:      title,
			Content:    content,
			Categories: categories,
		}
		// join the categories into a string
		cat := strings.Join(formData.Categories, " ")
		// get the user id from the session
		c, err := r.Cookie("session_token")
		if err != nil {
			// handle error
			return
		}
		userID = login.Sessions[c.Value].UserID
		author = login.Sessions[c.Value].Username

		// format the date and time in the correct format
		date = time.Now().Format("2006-01-02 15:04:05")

		// add the new post to the database
		sqlite.AddPost(title, content, date, userID, author, cat)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
