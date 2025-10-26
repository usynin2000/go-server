package repository

import (
	"database/sql"
	"time"

	"github.com/s.usynin/testing/go-server/internal/models"
)

type LikeRepository struct {
	db *sql.DB
}

func NewLikeRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

func (r *LikeRepository) Add(postID int) (int64, error) {
	query := `INSERT INTO likes (post_id, created_at) VALUES (?, datetime('now'))`

	result, err := r.db.Exec(query, postID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *LikeRepository) GetByID(id int) (*models.Like, error) {
	query := `SELECT id, post_id, created_at FROM likes WHERE id = ?`

	var like models.Like
	var createdAt string
	err := r.db.QueryRow(query, id).Scan(&like.ID, &like.PostID, &createdAt)
	if err != nil {
		return nil, err
	}

	like.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	return &like, nil
}
