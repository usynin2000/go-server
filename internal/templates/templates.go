package templates

import (
	"html/template"
)

var Tpl *template.Template

// InitTemplates загружает все шаблоны из директории templates
func InitTemplates() error {
	var err error

	Tpl, err = template.ParseFiles(
		"templates/home.html",
		"templates/post_item.html",
		"templates/comment_item.html",
	)

	return err
}
