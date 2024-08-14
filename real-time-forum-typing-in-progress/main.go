package main

import (
	"log"
	"net/http"
	"pillu/server"
	"pillu/sqlite"
)

func main() {

	log.Println("Starting application...")
	sqlite.CreateDB()

	// Create a file server for the "templates" directory
	templates := http.FileServer(http.Dir("templates"))

	// Handle requests for the files
	http.Handle("/templates/", http.StripPrefix("/templates/", templates))

	// Handle requests for the home page
	http.HandleFunc("/", server.Home)
	http.HandleFunc("/api/online-users", server.OnlineUsersAPI)
	http.HandleFunc("/api/comments", server.CommentsAPI)
	http.HandleFunc("/api/isLoggedIn", server.LoggedInHandler)
	http.HandleFunc("/api/get-messages", server.GetHistoryAPI)
	http.HandleFunc("/api/current-user", server.CurrentUserAPI)
	http.HandleFunc("/api/posts", server.PostsAPI)
	http.HandleFunc("/api/users", server.UsersAPI)

	server.SetupRoutes()

	// start the server
	log.Println("Starting server on port 80")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
