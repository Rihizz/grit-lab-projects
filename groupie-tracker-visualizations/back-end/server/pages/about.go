package pages

import (
	"log"
	"net/http"
	"text/template"
)

func About(w http.ResponseWriter, r *http.Request, index IndexData) {
	if r.URL.Path == "/about" { // If the request is for the about page, render the about page
		tmpl, err := template.ParseFiles("./front-end/templates/about.html")
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(w, index)
	}
}
