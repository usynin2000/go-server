package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Инициализация базы данных
	if err := InitDB(); err != nil {
		log.Fatal(err)
	}
	defer CloseDB()

	// Инициализация шаблонов
	if err := InitTemplates(); err != nil {
		log.Fatal(err)
	}

	// Настройка роутера
	r := setupRoutes()

	fmt.Println("Сервер запущен на http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
