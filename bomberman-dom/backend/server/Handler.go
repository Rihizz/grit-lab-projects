package server

import (
	"html/template"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	// Serve the index.html file
	tmpl := template.Must(template.ParseFiles("../frontend/index.html"))
	tmpl.Execute(w, nil)
}