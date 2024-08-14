package handlers

import (
	"fmt"
	"forum/server/login"
	util "forum/server/utils"
	"html/template"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	pageTemplate, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println("Logout - 500")
		ErrorHandling(500, w, r)
	}

	if r.URL.Path != "/logout/" {
		ErrorHandling(404, w, r)
		return
	} else if !login.CheckSession(w, r) {
		pageTemplate.Execute(w, M{
			"msg":   "You are not logged in!",
			"posts": util.Fposts,
			"cats":  util.Fcategories,
		})
	} else {
		c, err := r.Cookie("session_token")
		if err != nil {
			fmt.Println("Logout - 500. Can't get cookie.")
			ErrorHandling(500, w, r)
			return
		}
		delete(login.Sessions, c.Value)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
