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

func ViewPost(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var post struct {
		ID       int
		Title    string
		Content  string
		Author   string
		Category string
		Created  string
	}

	err := database.DB.QueryRow(`SELECT id, title, content, author, category, created_at FROM posts WHERE id = ?`, id).
		Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.Category, &post.Created)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	rows, err := database.DB.Query(`SELECT content, author, created_at FROM comments WHERE post_id = ? ORDER BY created_at ASC`, id)
	if err != nil {
		http.Error(w, "Error fetching comments", http.StatusInternalServerError)
		return
	}

	var comments []struct {
		ID       int
		Content  string
		Author   string
		Created  string
		Likes    int
		Dislikes int
	}

	for rows.Next() {
		var c struct {
			ID       int
			Content  string
			Author   string
			Created  string
			Likes    int
			Dislikes int
		}
		rows.Scan(&c.Content, &c.Author, &c.Created)
		comments = append(comments, c)
	}

	email, _ := sessions.GetUserEmail(r)

	tmpl, _ := template.ParseFiles("templates/post_detail.html")
	tmpl.Execute(w, struct {
		Email    string
		Post     any
		Comments []struct {
			ID       int
			Content  string
			Author   string
			Created  string
			Likes    int
			Dislikes int
		}
	}{
		Email:    email,
		Post:     post,
		Comments: comments,
	})
}

func SubmitComment(w http.ResponseWriter, r *http.Request) {
	email, ok := sessions.GetUserEmail(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	postID := r.FormValue("post_id")
	content := r.FormValue("content")

	_, err := database.DB.Exec(`INSERT INTO comments (post_id, author, content) VALUES (?, ?, ?)`, postID, email, content)
	if err != nil {
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post?id="+postID, http.StatusSeeOther)
}
