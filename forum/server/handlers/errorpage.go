package handlers

import (
	"net/http"
	"text/template"
)

// For handling errors
//
// Usage:
//
// ErrorHandling(404, w, r) --> 404 Not Found.
func ErrorHandling(code int, w http.ResponseWriter, r *http.Request) {
	template := template.Must(template.ParseFiles("templates/errorpage.html"))
	switch code {
	case 400:
		template.ExecuteTemplate(w, "errorpage.html", map[string]string{"Error": "400 Bad request."})
	case 401:
		template.ExecuteTemplate(w, "errorpage.html", map[string]string{"Error": "401 Unauthorized."})
	case 403:
		template.ExecuteTemplate(w, "errorpage.html", map[string]string{"Error": "403 Forbidden."})
	case 404:
		template.ExecuteTemplate(w, "errorpage.html", map[string]string{"Error": "404 Not found."})
	case 405:
		template.ExecuteTemplate(w, "errorpage.html", map[string]string{"Error": "405 Method not allowed."})
	case 500:
		template.ExecuteTemplate(w, "errorpage.html", map[string]string{"Error": "500 Internal server error."})
	default:
		return
	}
}
