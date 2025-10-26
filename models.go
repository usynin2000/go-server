package main

import "time"

// Post представляет блог-пост
type Post struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
}
