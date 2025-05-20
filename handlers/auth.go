package handlers

import (
	"html/template"
	"net/http"

	"forum/database"
	"forum/sessions"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, _ := template.ParseFiles("templates/register.html")
		tmpl.Execute(w, nil)
		return
	}

	r.ParseForm()
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Check if email exists
	row := database.DB.QueryRow("SELECT email FROM users WHERE email = ?", email)
	var existing string
	row.Scan(&existing)
	if existing != "" {
		http.Error(w, "Email already in use", http.StatusConflict)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err := database.DB.Exec("INSERT INTO users(email, username, password) VALUES (?, ?, ?)", email, username, hashedPassword)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, _ := template.ParseFiles("templates/login.html")
		tmpl.Execute(w, nil)
		return
	}

	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	var hashedPassword string
	err := database.DB.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&hashedPassword)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create session
	sessions.CreateSession(w, email)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
