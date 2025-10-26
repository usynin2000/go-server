package main

import (
	"database/sql"
)

// db - глобальная переменная для подключения к БД
var db *sql.DB

// InitDB инициализирует подключение к базе данных
func InitDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./blog.db")
	if err != nil {
		return err
	}

	// Настройка SQLite WAL режима
	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		return err
	}

	// Создание таблицы posts
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

// CloseDB закрывает подключение к базе данных
func CloseDB() {
	if db != nil {
		db.Close()
	}
}

// Функции для работы с базой данных

func getAllPosts() ([]Post, error) {
	rows, err := db.Query("SELECT id, title, content, created_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func createPost(title, content string) (int64, error) {
	result, err := db.Exec("INSERT INTO posts (title, content) VALUES (?, ?)", title, content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func getPostByID(id int64) (*Post, error) {
	var post Post
	err := db.QueryRow("SELECT id, title, content, created_at FROM posts WHERE id = ?", id).
		Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func deletePost(id int) error {
	_, err := db.Exec("DELETE FROM posts WHERE id = ?", id)
	return err
}
