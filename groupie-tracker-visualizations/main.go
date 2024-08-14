package main

import (
	"fmt"
	"groupie-tracker-visualizations/back-end/server"
)

func main() {
	fmt.Println("Server started at port 80")
	server.StartServer()
}
