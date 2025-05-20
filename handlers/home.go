package handlers

import (
	"html/template"
	"net/http"

	"forum/database"
	"forum/sessions"
)

func Home(w http.ResponseWriter, r *http.Request) {
	// tmpl, _ := template.New("home").Parse(`
	//     <h1>Welcome to the Forum</h1>
	//     {{if .Email}}<p>Logged in as {{.Email}}</p>{{else}}<p>Not logged in</p>{{end}}
	//     <a href="/login">Login</a> | <a href="/register">Register</a>
	// `)

	// email, ok := sessions.GetUserEmail(r)
	// data := struct {
	// 	Email string
	// }{}
	// if ok {
	// 	data.Email = email
	// }

	// tmpl.Execute(w, data)

	tmpl, _ := template.ParseFiles("templates/home.html")

	rows, err := database.DB.Query("SELECT id, title, content, author, category, created_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Post struct {
		ID       int
		Title    string
		Content  string
		Author   string
		Category string
		Created  string
		Likes    int
		Dislikes int
	}

	

	var posts []Post
	for rows.Next() {
		var p Post
		rows.Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.Category, &p.Created)

		// count likes/dislikes
		database.DB.QueryRow(`SELECT COUNT(*) FROM likes WHERE target_id = ? AND target_type = 'post' AND value = 1`, p.ID).Scan(&p.Likes)
		database.DB.QueryRow(`SELECT COUNT(*) FROM likes WHERE target_id = ? AND target_type = 'post' AND value = -1`, p.ID).Scan(&p.Dislikes)

		posts = append(posts, p)
	}

	email, _ := sessions.GetUserEmail(r)
	data := struct {
		Email string
		Posts []Post
	}{email, posts}

	tmpl.Execute(w, data)
}
