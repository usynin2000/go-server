package repository

import (
	"database/sql"
	"time"

	"github.com/s.usynin/testing/go-server/internal/models"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) GetByPostID(postID int) ([]models.Comment, error) {
	query := `
		SELECT id, post_id, author, content, created_at 
		FROM comments 
		WHERE post_id = ? 
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		var createdAt string
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.Author,
			&comment.Content, &createdAt)
		if err != nil {
			return nil, err
		}

		comment.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *CommentRepository) Create(postID int, author, content string) (int64, error) {
	query := `
		INSERT INTO comments (post_id, author, content, created_at) 
		VALUES (?, ?, ?, datetime('now'))
	`

	result, err := r.db.Exec(query, postID, author, content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *CommentRepository) GetByID(id int) (*models.Comment, error) {
	query := `SELECT id, post_id, author, content, created_at FROM comments WHERE id = ?`

	var comment models.Comment
	var createdAt string
	err := r.db.QueryRow(query, id).Scan(&comment.ID, &comment.PostID,
		&comment.Author, &comment.Content, &createdAt)
	if err != nil {
		return nil, err
	}

	comment.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	return &comment, nil
}
