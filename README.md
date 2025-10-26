# go-server

Простой блог на Go с HTMX и чистой архитектурой

## Stack

- **Backend**: Go 1.22+ с Chi router
- **Database**: SQLite с WAL режимом
- **Frontend**: HTMX + Tailwind CSS (server-side rendering)

## 🏗️ Архитектура проекта

Проект разделён на пакеты для масштабируемости:

```
go-server/
├── cmd/
│   └── server/
│       └── main.go           # Точка входа приложения
├── internal/
│   ├── models/               # Модели данных
│   │   └── models.go
│   ├── repository/           # Слой доступа к данным (Repository pattern)
│   │   ├── post_repository.go
│   │   ├── comment_repository.go
│   │   ├── category_repository.go
│   │   └── like_repository.go
│   ├── service/              # Бизнес-логика (Service layer)
│   │   └── post_service.go
│   ├── handlers/             # HTTP handlers
│   │   └── post_handler.go
│   ├── middleware/           # Middleware
│   │   ├── logging.go
│   │   └── recovery.go
│   ├── database/             # Работа с БД и миграции
│   │   ├── database.go
│   │   └── seed.go
│   └── templates/            # Управление шаблонами
│       └── templates.go
├── templates/                # HTML шаблоны
│   ├── home.html
│   ├── post_item.html
│   └── comment_item.html
├── bin/                      # Собранный бинарный файл
├── go.mod                    # Go модуль
├── go.sum                    # Checksums
├── Makefile                  # Команды для работы
└── README.md
```

## 📦 Структура пакетов

### cmd/server/
Точка входа приложения:
- Инициализирует БД
- Создаёт репозитории и сервисы
- Настраивает handlers
- Запускает HTTP сервер

### internal/models/
Модели данных:
- `Post` - посты блога
- `Comment` - комментарии
- `Category` - категории
- `Like` - лайки

### internal/repository/
Repository pattern - работа с БД:
- `PostRepository` - CRUD постов
- `CommentRepository` - управление комментариями
- `CategoryRepository` - управление категориями
- `LikeRepository` - управление лайками

### internal/service/
Service layer - бизнес-логика:
- `PostService` - координация работы с постами
- Объединяет несколько репозиториев
- Выполняет валидацию и трансформацию данных

### internal/handlers/
HTTP handlers:
- `PostHandler` - обработка HTTP запросов
- Использует service слой
- Рендерит шаблоны

### internal/middleware/
Middleware компоненты:
- `LoggingMiddleware` - логирование запросов
- `RecoveryMiddleware` - обработка паник

### internal/database/
Работа с БД:
- `InitDB()` - подключение и миграции
- `SeedDatabase()` - заполнение начальными данными
- Все SQL миграции в одном месте

## 🎯 Преимущества новой архитектуры

1. **Разделение ответственности** - каждый пакет отвечает за свою область
2. **Масштабируемость** - легко добавлять новые функции
3. **Тестируемость** - можно мокировать репозитории и тестировать сервисы
4. **Поддерживаемость** - понятная структура для новых разработчиков
5. **Готовность к росту** - легко добавлять новые слои (auth, validation, etc.)

## 🚀 Запуск проекта

### Базовый запуск
```bash
make run
# или
go run ./cmd/server
```

### Сборка бинарного файла
```bash
make build
./bin/server
```

### Очистка БД и перезапуск
```bash
make clean
```

### Очистить порт
```bash
make clear-port
```

## 📋 API Endpoints

- `GET /` - Главная страница со всеми постами
- `POST /posts` - Создать новый пост
- `DELETE /posts/{id}` - Удалить пост
- `POST /comments` - Добавить комментарий к посту
- `POST /likes` - Добавить лайк к посту

## 🎨 Функционал

- ✅ Создание, просмотр и удаление постов
- ✅ Категории постов
- ✅ Комментарии к постам
- ✅ Счётчики комментариев и лайков
- ✅ Красивый UI с Tailwind CSS
- ✅ HTMX для интерактивности без JavaScript
- ✅ Чистая архитектура для масштабирования

## 🏃‍♂️ Разработка

### Добавление новой функциональности

1. Добавьте модель в `internal/models/models.go`
2. Создайте репозиторий в `internal/repository/`
3. Добавьте методы в service в `internal/service/post_service.go`
4. Создайте handler в `internal/handlers/post_handler.go`
5. Добавьте роуты в `cmd/server/main.go`
6. Создайте шаблон в `templates/`

### Добавление миграций

Добавьте SQL в массив `migrations` в `internal/database/database.go`:

```go
const createMyTable = `
CREATE TABLE IF NOT EXISTS my_table (...)
`
```

## 📚 Технологии

- [Chi](https://github.com/go-chi/chi) - HTTP router
- [HTMX](https://htmx.org/) - HTML over the wire
- [Tailwind CSS](https://tailwindcss.com/) - Utility-first CSS
- [SQLite](https://www.sqlite.org/) - Embedded database
