package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// setupRoutes настраивает роутер и возвращает его
func setupRoutes() http.Handler {
	// Создаём новый роутер с помощью chi.NewRouter()
	r := chi.NewRouter()

	// Подключаем middleware для логирования запросов
	r.Use(middleware.Logger)

	// Middleware для обработки паник
	r.Use(middleware.Recoverer)

	// Роуты
	r.Get("/", homeHandler)
	r.Post("/posts", createPostHandler)
	r.Delete("/posts/{id}", deletePostHandler)

	// Статические файлы (если понадобятся)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return r
}
