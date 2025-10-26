package repository

import (
	"database/sql"
	"time"

	"github.com/s.usynin/testing/go-server/internal/models"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) GetAll() ([]models.Post, error) {
	query := `
		SELECT p.id, p.title, p.content, p.category_id, p.created_at, p.updated_at,
		       COUNT(DISTINCT c.id) as comments_count,
		       COUNT(DISTINCT l.id) as likes_count
		FROM posts p
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN likes l ON p.id = l.post_id
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var createdAt, updatedAt string
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CategoryID,
			&createdAt, &updatedAt, &post.CommentsCount, &post.LikesCount)
		if err != nil {
			return nil, err
		}

		post.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		post.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) GetByID(id int) (*models.Post, error) {
	query := `
		SELECT p.id, p.title, p.content, p.category_id, p.created_at, p.updated_at,
		       COUNT(DISTINCT c.id) as comments_count,
		       COUNT(DISTINCT l.id) as likes_count
		FROM posts p
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN likes l ON p.id = l.post_id
		WHERE p.id = ?
		GROUP BY p.id
	`

	var post models.Post
	var createdAt, updatedAt string
	err := r.db.QueryRow(query, id).Scan(&post.ID, &post.Title, &post.Content,
		&post.CategoryID, &createdAt, &updatedAt, &post.CommentsCount, &post.LikesCount)
	if err != nil {
		return nil, err
	}

	post.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	post.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	return &post, nil
}

func (r *PostRepository) Create(title, content string, categoryID int) (int64, error) {
	query := `
		INSERT INTO posts (title, content, category_id, created_at, updated_at) 
		VALUES (?, ?, ?, datetime('now'), datetime('now'))
	`

	result, err := r.db.Exec(query, title, content, categoryID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostRepository) Delete(id int) error {
	query := `DELETE FROM posts WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *PostRepository) GetByCategoryID(categoryID int) ([]models.Post, error) {
	query := `
		SELECT p.id, p.title, p.content, p.category_id, p.created_at, p.updated_at,
		       COUNT(DISTINCT c.id) as comments_count,
		       COUNT(DISTINCT l.id) as likes_count
		FROM posts p
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN likes l ON p.id = l.post_id
		WHERE p.category_id = ?
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`

	rows, err := r.db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var createdAt, updatedAt string
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CategoryID,
			&createdAt, &updatedAt, &post.CommentsCount, &post.LikesCount)
		if err != nil {
			return nil, err
		}

		post.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		post.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

		posts = append(posts, post)
	}

	return posts, nil
}
