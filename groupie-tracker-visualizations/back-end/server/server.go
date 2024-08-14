package server

import (
	"groupie-tracker-visualizations/back-end/server/getData.go"
	"groupie-tracker-visualizations/back-end/server/pages"
	"log"
	"net/http"
)

// startServer starts the server
func StartServer() {
	http.HandleFunc("/", handler)
	css := http.FileServer(http.Dir("front-end/css"))
	http.Handle("/css/", http.StripPrefix("/css/", css))
	log.Fatal(http.ListenAndServe(":80", nil))
}

// handler handles the request
func handler(w http.ResponseWriter, r *http.Request) {
	index := pages.IndexData{}
	index.Data = getData.GetArtists()
	pages.Home(w, r, index)
	pages.Artist(w, r, index.Data)
	pages.About(w, r, index)
	pages.Donate(w, r, index)
	if pages.Error(w, r, index) {
		return
	}
}
