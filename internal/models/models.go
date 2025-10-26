package models

import "time"

// Post представляет блог-пост
type Post struct {
	ID         int       `json:"id" db:"id"`
	Title      string    `json:"title" db:"title"`
	Content    string    `json:"content" db:"content"`
	CategoryID int       `json:"category_id" db:"category_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`

	// Связи (для загрузки)
	Category      *Category `json:"category,omitempty"`
	Comments      []Comment `json:"comments,omitempty"`
	CommentsCount int       `json:"comments_count"`
	LikesCount    int       `json:"likes_count"`
}

// Category представляет категорию поста
type Category struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Comment представляет комментарий к посту
type Comment struct {
	ID        int       `json:"id" db:"id"`
	PostID    int       `json:"post_id" db:"post_id"`
	Author    string    `json:"author" db:"author"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Like представляет лайк к посту
type Like struct {
	ID        int       `json:"id" db:"id"`
	PostID    int       `json:"post_id" db:"post_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
