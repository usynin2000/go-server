package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/s.usynin/testing/go-server/internal/models"
	"github.com/s.usynin/testing/go-server/internal/service"
)

type PostHandler struct {
	postService *service.PostService
	templates   *template.Template
}

func NewPostHandler(postService *service.PostService, templates *template.Template) *PostHandler {
	return &PostHandler{
		postService: postService,
		templates:   templates,
	}
}

func (h *PostHandler) Home(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.GetAllPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categories, _ := h.postService.GetCategories()

	data := struct {
		Posts      []models.Post
		Categories []models.Category
	}{
		Posts:      posts,
		Categories: categories,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = h.templates.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	categoryIDStr := r.FormValue("category_id")

	if title == "" || content == "" {
		http.Error(w, "Заголовок и содержание обязательны", http.StatusBadRequest)
		return
	}

	categoryID := 1 // Дефолтная категория
	if categoryIDStr != "" {
		categoryID, _ = strconv.Atoi(categoryIDStr)
	}

	post, err := h.postService.CreatePost(title, content, categoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = h.templates.ExecuteTemplate(w, "post_item.html", post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	err = h.postService.DeletePost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PostHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.FormValue("post_id")
	author := r.FormValue("author")
	content := r.FormValue("content")

	if author == "" || content == "" || postIDStr == "" {
		http.Error(w, "Все поля обязательны", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	comment, err := h.postService.AddComment(postID, author, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = h.templates.ExecuteTemplate(w, "comment_item.html", comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) AddLike(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.FormValue("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	_, err = h.postService.AddLike(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
