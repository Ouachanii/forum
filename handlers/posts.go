package handlers

import (
	"html/template"
	"net/http"

	"forum/database"
	"forum/sessions"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	email, ok := sessions.GetUserEmail(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == "GET" {
		tmpl, _ := template.ParseFiles("templates/create_post.html")
		tmpl.Execute(w, nil)
		return
	}

	r.ParseForm()
	title := r.FormValue("title")
	content := r.FormValue("content")
	category := r.FormValue("category")

	_, err := database.DB.Exec("INSERT INTO posts (title, content, author, category) VALUES (?, ?, ?, ?)",
		title, content, email, category)
	if err != nil {
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
