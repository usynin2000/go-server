package database

import (
	"database/sql"
	"log"
)

// SeedDatabase заполняет базу данных начальными данными
func SeedDatabase(db *sql.DB) error {
	// Проверяем, есть ли уже категории
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count)
	if err != nil {
		return err
	}

	// Если категории уже есть, не делаем ничего
	if count > 0 {
		log.Println("Database already seeded")
		return nil
	}

	// Создаём начальные категории
	categories := []struct {
		name string
		slug string
	}{
		{"Общие", "general"},
		{"Технологии", "technology"},
		{"Путешествия", "travel"},
		{"Кулинария", "cooking"},
	}

	for _, cat := range categories {
		_, err := db.Exec("INSERT INTO categories (name, slug) VALUES (?, ?)", cat.name, cat.slug)
		if err != nil {
			return err
		}
	}

	log.Println("Database seeded successfully")
	return nil
}
