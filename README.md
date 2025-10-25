# go-server

# Testing project with stack

Backend: Go 1.23+ с Chi router
Database: SQLite с WAL режимом
Frontend: HTMX + Tailwind CSS (server-side rendering)


1. Инициализация go.mod
```shell
go mod init go-server
# Определяет имя модуля (в нашем случае go-server)
```


2. Установка зависимостей
```shell
go mod tidy
```

### Зачем это нужно:
- Автоматически скачивает все пакеты, которые импортируются в коде
- Создаёт файл go.sum с хешами для проверки целостности
- Удаляет неиспользуемые зависимости
В нашем случае скачает github.com/go-chi/chi/v5 и github.com/mattn/go-sqlite3