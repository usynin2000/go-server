package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"

	"github.com/s.usynin/testing/go-server/internal/database"
	"github.com/s.usynin/testing/go-server/internal/handlers"
	"github.com/s.usynin/testing/go-server/internal/repository"
	"github.com/s.usynin/testing/go-server/internal/service"
	templatesPkg "github.com/s.usynin/testing/go-server/internal/templates"  // только чтобы избежать конфликт имен
)

func main() {
	// Инициализация базы данных
	db, err := database.InitDB("./blog.db")
	if err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}
	defer db.Close()

	// Заполняем базу начальными данными
	if err := database.SeedDatabase(db); err != nil {
		log.Fatal("Ошибка заполнения БД:", err)
	}

	// Инициализация шаблонов
	if err := templatesPkg.InitTemplates(); err != nil {
		log.Fatal("Ошибка инициализации шаблонов:", err)
	}

	// Создаём репозитории
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	likeRepo := repository.NewLikeRepository(db)

	// Создаём сервисы
	postService := service.NewPostService(postRepo, commentRepo, categoryRepo, likeRepo)

	// Создаём handlers
	postHandler := handlers.NewPostHandler(postService, templatesPkg.Tpl)

	// Настройка роутера
	r := setupRoutes(postHandler)

	fmt.Println("Сервер запущен на http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}

func setupRoutes(postHandler *handlers.PostHandler) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Роуты
	r.Get("/", postHandler.Home)
	r.Post("/posts", postHandler.CreatePost)
	r.Delete("/posts/{id}", postHandler.DeletePost)
	r.Post("/comments", postHandler.AddComment)
	r.Post("/likes", postHandler.AddLike)

	// Статические файлы
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return r
}
