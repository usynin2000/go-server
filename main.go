package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
}

var db *sql.DB

func main() {
	// Инициализация базы данных
	var err error
	db, err = sql.Open("sqlite3", "./blog.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Настройка SQLite WAL режима
	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		log.Fatal(err)
	}

	// Создание таблицы posts
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Здесь происходит настройка роутера
	// 1. Создаём новый роутер с помощью chi.NewRouter(). Chi — это быстрый и удобный HTTP-роутер для Go.
	// 2. r.Use(middleware.Logger) — подключает middleware, который логирует все HTTP-запросы в консоль. Так удобно отслеживать обращения к серверу.
	// 3. r.Use(middleware.Recoverer) — этот middleware «ловит» паники (panic) в обработчиках запросов и не даёт всему серверу упасть из-за одной ошибки: он возвращает клиенту 500, а паника логируется в консоль.
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Роуты
	r.Get("/", homeHandler)
	r.Post("/posts", createPostHandler)
	r.Delete("/posts/{id}", deletePostHandler)

	// Статические файлы (если понадобятся)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Сервер запущен на http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := getAllPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := `
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Простой блог</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-100 min-h-screen">
    <div class="container mx-auto px-4 py-8 max-w-4xl">
        <h1 class="text-4xl font-bold text-gray-800 mb-8 text-center">Простой блог</h1>
        
        <!-- Форма добавления поста -->
        <div class="bg-white rounded-lg shadow-md p-6 mb-8">
            <h2 class="text-2xl font-semibold text-gray-700 mb-4">Добавить новый пост</h2>
            <form hx-post="/posts" hx-target="#posts-list" hx-swap="afterbegin" class="space-y-4">
                <div>
                    <label for="title" class="block text-sm font-medium text-gray-700 mb-2">Заголовок</label>
                    <input type="text" id="title" name="title" required 
                           class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
                </div>
                <div>
                    <label for="content" class="block text-sm font-medium text-gray-700 mb-2">Содержание</label>
                    <textarea id="content" name="content" rows="4" required
                              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"></textarea>
                </div>
                <button type="submit" 
                        class="bg-blue-500 hover:bg-blue-600 text-white font-medium py-2 px-4 rounded-md transition duration-200">
                    Добавить пост
                </button>
            </form>
        </div>

        <!-- Список постов -->
        <div class="bg-white rounded-lg shadow-md p-6">
            <h2 class="text-2xl font-semibold text-gray-700 mb-4">Все посты</h2>
            <div id="posts-list" class="space-y-4">
                {{range .}}
                <div class="border border-gray-200 rounded-lg p-4 hover:shadow-md transition duration-200">
                    <div class="flex justify-between items-start">
                        <div class="flex-1">
                            <h3 class="text-xl font-semibold text-gray-800 mb-2">{{.Title}}</h3>
                            <p class="text-gray-600 mb-2">{{.Content}}</p>
                            <p class="text-sm text-gray-500">{{.CreatedAt.Format "02.01.2006 15:04"}}</p>
                        </div>
                        <button hx-delete="/posts/{{.ID}}" hx-target="closest div" hx-swap="outerHTML"
                                class="bg-red-500 hover:bg-red-600 text-white font-medium py-1 px-3 rounded-md transition duration-200 ml-4">
                            Удалить
                        </button>
                    </div>
                </div>
                {{end}}
            </div>
        </div>
    </div>
</body>
</html>`

	t, err := template.New("home").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(w, posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" || content == "" {
		http.Error(w, "Заголовок и содержание обязательны", http.StatusBadRequest)
		return
	}

	id, err := createPost(title, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем созданный пост для отображения
	post, err := getPostByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем HTML фрагмент для HTMX
	tmpl := `
<div class="border border-gray-200 rounded-lg p-4 hover:shadow-md transition duration-200">
    <div class="flex justify-between items-start">
        <div class="flex-1">
            <h3 class="text-xl font-semibold text-gray-800 mb-2">{{.Title}}</h3>
            <p class="text-gray-600 mb-2">{{.Content}}</p>
            <p class="text-sm text-gray-500">{{.CreatedAt.Format "02.01.2006 15:04"}}</p>
        </div>
        <button hx-delete="/posts/{{.ID}}" hx-target="closest div" hx-swap="outerHTML"
                class="bg-red-500 hover:bg-red-600 text-white font-medium py-1 px-3 rounded-md transition duration-200 ml-4">
            Удалить
        </button>
    </div>
</div>`

	t, err := template.New("post").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(w, post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	err = deletePost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Функции для работы с базой данных

func getAllPosts() ([]Post, error) {
	rows, err := db.Query("SELECT id, title, content, created_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func createPost(title, content string) (int64, error) {
	result, err := db.Exec("INSERT INTO posts (title, content) VALUES (?, ?)", title, content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func getPostByID(id int64) (*Post, error) {
	var post Post
	err := db.QueryRow("SELECT id, title, content, created_at FROM posts WHERE id = ?", id).
		Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func deletePost(id int) error {
	_, err := db.Exec("DELETE FROM posts WHERE id = ?", id)
	return err
}
