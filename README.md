# go-server

# Testing project with stack

Backend: Go 1.23+ с Chi router
Database: SQLite с WAL режимом
Frontend: HTMX + Tailwind CSS (server-side rendering)

## Структура проекта

```
go-server/
├── main.go          # Точка входа, инициализация и запуск сервера
├── models.go        # Модели данных (структуры)
├── database.go      # Подключение к БД и функции работы с данными
├── templates.go     # Загрузка HTML шаблонов
├── handlers.go      # HTTP handlers (обработчики запросов)
├── routes.go        # Настройка роутинга
└── templates/       # HTML шаблоны
    ├── home.html           # Главная страница
    └── post_item.html      # Фрагмент поста для HTMX
```

## Объяснение структуры

### main.go
- Инициализирует базу данных (`InitDB()`)
- Загружает шаблоны (`InitTemplates()`)
- Настраивает роуты (`setupRoutes()`)
- Запускает HTTP сервер

### models.go
- Содержит структуру `Post` для представления поста в блоге

### database.go
- `var db *sql.DB` - глобальное подключение к БД
- `InitDB()` - создаёт подключение и таблицы
- `CloseDB()` - закрывает подключение
- Функции для работы с данными: `getAllPosts()`, `createPost()`, `getPostByID()`, `deletePost()`

### templates.go
- Загружает HTML шаблоны из директории `templates/`
- Используется в handlers для отображения

### handlers.go
- `homeHandler` - отображает главную страницу со списком постов
- `createPostHandler` - создаёт новый пост и возвращает HTMX фрагмент
- `deletePostHandler` - удаляет пост

### routes.go
- Настраивает middleware (Logger, Recoverer)
- Определяет маршруты API
- Возвращает готовый роутер

## Как это работает

1. **main.go** запускается и инициализирует всё необходимое
2. Запрос приходит на определённый роут
3. **routes.go** направляет запрос в нужный handler
4. **handlers.go** обрабатывает запрос, используя функции из **database.go**
5. Результат рендерится через шаблоны из **templates/** и возвращается клиенту

## Команды

### 1. Инициализация go.mod
```shell
go mod init go-server
# Определяет имя модуля (в нашем случае go-server)
```

### 2. Установка зависимостей
```shell
go mod tidy
```

### Зачем это нужно:
- Автоматически скачивает все пакеты, которые импортируются в коде
- Создаёт файл go.sum с хешами для проверки целостности
- Удаляет неиспользуемые зависимости
В нашем случае скачает github.com/go-chi/chi/v5 и github.com/mattn/go-sqlite3

### 3. Запуск сервера
```shell
go run .
```

Сервер будет доступен по адресу http://localhost:3000