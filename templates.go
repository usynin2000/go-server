package main

import (
	"html/template"
)

var templates *template.Template

// InitTemplates загружает все шаблоны из директории templates
func InitTemplates() error {
	var err error

	// Используем ParseFiles для загрузки конкретных шаблонов
	templates, err = template.ParseFiles(
		"templates/home.html",
		"templates/post_item.html",
	)

	return err
}
