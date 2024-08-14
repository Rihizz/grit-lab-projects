package main

import (
	"fmt"
	"forum/server/handlers"
	"log"
	"net/http"
)

func main() {
	// encpsw, err := util.HashPassword("1")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// sqlite.CreateDB()
	// sqlite.AddCategory("Sports")
	// sqlite.AddCategory("Politics")
	// sqlite.AddCategory("Technology")
	// sqlite.RegisterUser("admin@localhost", "admin", encpsw, "admin")
	// sqlite.AddPost("This is a test post", "test post connteennnntttntntntn", "2023-01-01", 1, "admin", "Sports")
	// sqlite.AddPost("This is a test post", "test post connteennnntttntntntn", "2023-01-01", 1, "admin", "Politics")
	// sqlite.AddPost("This is a test post", "test post connteennnntttntntntn", "2023-01-01", 1, "admin", "Technology")

	css := http.FileServer(http.Dir("templates/css"))
	js := http.FileServer(http.Dir("templates/js"))
	http.HandleFunc("/", handlers.Handler)
	http.Handle("/css/", http.StripPrefix("/css/", css))
	http.Handle("/js/", http.StripPrefix("/js/", js))
	fmt.Println("Starting server on port 80")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
