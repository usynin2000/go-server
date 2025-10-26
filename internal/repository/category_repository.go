package repository

import (
	"database/sql"
	"time"

	"github.com/s.usynin/testing/go-server/internal/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	query := `SELECT id, name, slug, created_at FROM categories ORDER BY name ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		var createdAt string
		err := rows.Scan(&category.ID, &category.Name, &category.Slug, &createdAt)
		if err != nil {
			return nil, err
		}

		category.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	query := `SELECT id, name, slug, created_at FROM categories WHERE id = ?`

	var category models.Category
	var createdAt string
	err := r.db.QueryRow(query, id).Scan(&category.ID, &category.Name, &category.Slug, &createdAt)
	if err != nil {
		return nil, err
	}

	category.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	return &category, nil
}

func (r *CategoryRepository) Create(name, slug string) (int64, error) {
	query := `
		INSERT INTO categories (name, slug, created_at) 
		VALUES (?, ?, datetime('now'))
	`

	result, err := r.db.Exec(query, name, slug)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
