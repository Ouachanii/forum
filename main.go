package main

import (
	"net/http"
	"forum/database"
	"forum/handlers"
)

func main() {
	database.Init()

	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/home", handlers.Home)
	http.HandleFunc("/create_post", handlers.CreatePost)
	http.ListenAndServe(":8080", nil)
}
