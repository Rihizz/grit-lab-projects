package main

import (
	"bomberman-dom/server"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting server...")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Register the server.Home handler for the root path
	http.HandleFunc("/", server.Home)

	// setup websocket route
	server.SetupRoutes()

	// Start a server on port 80
	log.Fatal(http.ListenAndServe(":80", nil))
}
