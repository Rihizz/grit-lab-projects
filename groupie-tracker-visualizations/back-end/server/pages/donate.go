package pages

import (
	"log"
	"net/http"
	"text/template"
)

func Donate(w http.ResponseWriter, r *http.Request, index IndexData) {
	if r.URL.Path == "/donate" {
		tmpl, err := template.ParseFiles("./front-end/templates/donate.html")
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(w, index)
	}
}
