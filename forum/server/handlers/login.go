package handlers

import (
	"forum/server/login"
	"forum/server/sqlite"
	util "forum/server/utils"
	"log"
	"net/http"
	"text/template"
)

// processing form and checking login credentials
func Login(w http.ResponseWriter, r *http.Request) {
	var username, psw string
	pageTemplate := template.Must(template.ParseFiles("templates/index.html"))

	// Check if the user is already logged in
	if login.CheckSession(w, r) {
		pageTemplate.ExecuteTemplate(w, "index.html", M{
			"cats":  util.Fcategories,
			"posts": util.Fposts,
			"msg":   "You are already logged in.",
		})
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	switch r.Method {
	case "GET":
		// if method is GET, error. Only POST should happen
		if err := pageTemplate.Execute(w, nil); err != nil {
			log.Println("Login handler error: ", err)
			ErrorHandling(405, w, r)
			return
		}
	case "POST":
		r.ParseForm()
		username, psw = r.FormValue("username"), r.FormValue("psw")
		if userID, err := sqlite.VerifyLogin(username, string(psw)); err == nil && userID > 0 {
			login.Signin(w, r)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			pageTemplate.ExecuteTemplate(w, "index.html", M{
				"cats":  util.Fcategories,
				"posts": util.Fposts,
				"msg":   "Wrong username or password.",
			})

		}
	default:
		ErrorHandling(405, w, r)
		return
	}
}
