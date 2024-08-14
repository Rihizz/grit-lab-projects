package pages

import (
	"groupie-tracker-visualizations/back-end/server/getData.go"
	"log"
	"net/http"
	"text/template"
)

func allowedURLs(s string, data []getData.Artist) bool {
	urls := []string{}
	for _, v := range data {
		urls = append(urls, "/"+v.Name)
	}
	for _, v := range urls {
		if s == v {
			return true
		}
	}
	return false
}

func Error(w http.ResponseWriter, r *http.Request, index IndexData) bool {
	if r.Method != "GET" {
		// If the request is not for the home page, render the 404 page
		tmpl, err := template.ParseFiles("./front-end/templates/405.html")
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(w, "405.html")
	}
	if r.URL.Path != "/" && r.URL.Path != "/artist" && r.URL.Path != "/about" && r.URL.Path != "/css/style.css" && r.URL.Path != "/donate" && !allowedURLs(r.URL.Path, index.Data) {
		// If the request is not for the home page, render the 404 page
		tmpl, err := template.ParseFiles("./front-end/templates/404.html")
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(w, "404.html")
	}
	return false
}
