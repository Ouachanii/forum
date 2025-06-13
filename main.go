package main

import (
	"fmt"
	"forum/database"
	"forum/handlers"
	"net/http"
)

func main() {
	database.Init()

	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/home", handlers.Home)
	http.HandleFunc("/create_post", handlers.CreatePost)
	http.HandleFunc("/post", handlers.ViewPost)
	http.HandleFunc("/comment", handlers.SubmitComment)

	http.ListenAndServe(":8080", nil)
	fmt.Println("Server started on :8080\nto view the forum visit: http://localhost:8080/home")
}
