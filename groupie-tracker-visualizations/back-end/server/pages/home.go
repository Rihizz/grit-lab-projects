package pages

import (
	"groupie-tracker-visualizations/back-end/server/getData.go"
	"log"
	"net/http"
	"text/template"
)

type IndexData struct {
	Data []getData.Artist
}

func Home(w http.ResponseWriter, r *http.Request, index IndexData) {
	if r.URL.Path == "/" && r.Method == "GET" { // If the request is for the home page, render the home page
		tmpl, err := template.ParseFiles("./front-end/templates/index.html")
		if err != nil {
			tmpl, err := template.ParseFiles("./front-end/templates/500.html")
			if err != nil {
				log.Fatal(err)
			}
			tmpl.Execute(w, "500.html")
		}
		tmpl.Execute(w, index)
	}
}
