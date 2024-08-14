package handlers

import (
	"forum/server/login"
	util "forum/server/utils"
	"log"
	"net/http"
	"text/template"
)

type M map[string]interface{}

// If the request is for the home page, render the home page
func Home(w http.ResponseWriter, r *http.Request) {
	util.GetPosts()
	util.GetCategories()
	if r.URL.Path != "/" {
		ErrorHandling(404, w, r)
		return
	}
	pageTemplate, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println("Home - 500")
		ErrorHandling(500, w, r)
		return
	}

	if r.Method != "GET" {
		log.Println("Home - 405")
		ErrorHandling(405, w, r)
		return
	}

	if login.CheckSession(w, r) {
		c, err := r.Cookie("session_token")
		if err != nil {
			log.Println("Home - 500. Can't get cookie.")
			ErrorHandling(500, w, r)
			return
		}

		userd := login.Sessions[c.Value]
		pageTemplate.Execute(w, M{
			"posts": util.Fposts,
			"cats":  util.Fcategories,
			"user":  userd,
			"msg":   "Welcome back, " + userd.Username + "!",
		})
		return
	} else {
		pageTemplate.Execute(w, M{
			"cats":  util.Fcategories,
			"posts": util.Fposts,
			"user":  nil,
		})
		return
	}
}
