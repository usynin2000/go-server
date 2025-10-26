package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB инициализирует подключение к базе данных и выполняет миграции
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Настройка SQLite WAL режима
	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		return nil, err
	}

	// Выполняем миграции
	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

// runMigrations выполняет все миграции базы данных
func runMigrations(db *sql.DB) error {
	migrations := []string{
		createPostsTable,
		createCategoriesTable,
		createCommentsTable,
		createLikesTable,
	}

	for i, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			log.Printf("Migration %d failed: %v", i+1, err)
			return err
		}
		log.Printf("Migration %d completed successfully", i+1)
	}

	return nil
}

// SQL миграции
const createPostsTable = `
CREATE TABLE IF NOT EXISTS posts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	category_id INTEGER NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (category_id) REFERENCES categories(id)
)`

const createCategoriesTable = `
CREATE TABLE IF NOT EXISTS categories (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	slug TEXT NOT NULL UNIQUE,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
)`

const createCommentsTable = `
CREATE TABLE IF NOT EXISTS comments (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id INTEGER NOT NULL,
	author TEXT NOT NULL,
	content TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
)`

const createLikesTable = `
CREATE TABLE IF NOT EXISTS likes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id INTEGER NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
)`
