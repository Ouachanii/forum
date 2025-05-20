package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}

	createTables()
}

func createTables() {
	createUserTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        email TEXT NOT NULL UNIQUE,
        username TEXT NOT NULL,
        password TEXT NOT NULL
    );`

	_, err := DB.Exec(createUserTable)
	if err != nil {
		log.Fatal(err)
	}
	createPostTable := `
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    author TEXT NOT NULL,
    category TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`
_, err = DB.Exec(createPostTable)
if err != nil {
    log.Fatal(err)
}
	createCommentTable := `
CREATE TABLE IF NOT EXISTS comments (
	id INTEGER PRIMARY KEY AUTOINCREMENT,	
	post_id INTEGER NOT NULL,
	content TEXT NOT NULL,
	author TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
);`
	_, err = DB.Exec(createCommentTable)
	if err != nil {
		log.Fatal(err)
	}
	createLikesTable := `
CREATE TABLE IF NOT EXISTS likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user TEXT NOT NULL,
    target_id INTEGER NOT NULL,
    target_type TEXT NOT NULL, -- 'post' or 'comment'
    value INTEGER NOT NULL, -- 1 = like, -1 = dislike
    UNIQUE(user, target_id, target_type)
);`
_, err = DB.Exec(createLikesTable)
if err != nil {
    log.Fatal(err)
}
	// Create categories table
	createCategoryTable := `
CREATE TABLE IF NOT EXISTS categories (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE
);`
	_, err = DB.Exec(createCategoryTable)
	if err != nil {	
		log.Fatal(err)
	}
}
