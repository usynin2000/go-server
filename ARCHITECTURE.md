# Архитектурные изменения проекта

## Что было изменено

Проект был переорганизован из монолитного кода в **чистую архитектуру** с разделением на пакеты.

## Было → Стало

### Было:
```
go-server/
├── main.go           # Все в одном файле
├── handlers.go       # Handlers в корне
├── database.go       # БД в корне
├── models.go         # Модели в корне
├── templates.go      # Шаблоны в корне
└── routes.go         # Роуты в корне
```

### Стало:
```
go-server/
├── cmd/server/
│   └── main.go           # Точка входа
├── internal/
│   ├── models/           # Модели данных
│   ├── repository/       # Доступ к данным (Repository Pattern)
│   ├── service/          # Бизнес-логика
│   ├── handlers/         # HTTP handlers
│   ├── middleware/       # Middleware компоненты
│   ├── database/         # Миграции и работа с БД
│   └── templates/         # Управление шаблонами
└── templates/            # HTML шаблоны
```

## Слои архитектуры

### 1. Models (`internal/models/`)
**Ответственность**: Определение структуры данных

```go
type Post struct {
    ID          int
    Title       string
    Content     string
    CategoryID  int
    // ...
}
```

### 2. Repository (`internal/repository/`)
**Ответственность**: Работа с базой данных

```go
type PostRepository struct {
    db *sql.DB
}

func (r *PostRepository) GetAll() ([]models.Post, error) {
    // SQL запросы здесь
}
```

**Преимущества**:
- Инкапсуляция SQL логики
- Легко мокировать для тестирования
- Один репозиторий = одна сущность

### 3. Service (`internal/service/`)
**Ответственность**: Бизнес-логика и координация

```go
type PostService struct {
    postRepo     *repository.PostRepository
    commentRepo  *repository.CommentRepository
    categoryRepo *repository.CategoryRepository
    likeRepo     *repository.LikeRepository
}

func (s *PostService) CreatePost(title, content string, categoryID int) (*models.Post, error) {
    // Координирует работу нескольких репозиториев
    // Может выполнять валидацию, трансформацию данных
}
```

**Преимущества**:
- Изоляция бизнес-логики от HTTP слоя
- Можно переиспользовать в других контекстах (CLI, gRPC)
- Легко тестировать

### 4. Handlers (`internal/handlers/`)
**Ответственность**: HTTP запросы/ответы

```go
type PostHandler struct {
    postService *service.PostService
    templates   *template.Template
}

func (h *PostHandler) Home(w http.ResponseWriter, r *http.Request) {
    posts, err := h.postService.GetAllPosts()
    // Отправляет данные в шаблон
}
```

**Преимущества**:
- Тонкий слой, только HTTP специфика
- Не зависит от БД напрямую
- Легко добавлять новые endpoints

### 5. Database (`internal/database/`)
**Ответственность**: Инициализация БД и миграции

```go
func InitDB(dbPath string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", dbPath)
    runMigrations(db)  // Все миграции в одном месте
    return db, nil
}
```

**Преимущества**:
- Версионирование схемы БД
- Легко добавлять новые таблицы
- Семплинг начальных данных

## Почему это важно?

### Масштабируемость
При добавлении новых функций (пользователи, авторизация, etc.) структура уже готова:
- Добавил модель → создал репозиторий → добавил в service → создал handler

### Тестируемость
```go
// Можно легко мокировать репозиторий
mockRepo := &MockPostRepository{}
service := NewPostService(mockRepo, ...)
// Тестируем бизнес-логику без БД
```

### Разделение ответственности
- Handlers знают только про HTTP
- Service знает только про бизнес-логику
- Repository знает только про БД

### Готовность к росту
Можно легко добавить:
- Authentication/Authorization слой
- Validation слой
- Caching слой
- Event Bus
- gRPC handlers рядом с HTTP

## Добавленный функционал

1. **Категории постов** - структурирование контента
2. **Комментарии** - интерактивность
3. **Лайки** - социальные функции
4. **Счётчики** - статистика

## Сравнение подходов

### Старый подход (всё в одном месте)
```go
// main.go
db := InitDB()
posts := getAllPosts(db)
// обработка запроса
```

Проблемы:
- Всё связано
- Сложно тестировать
- Сложно масштабировать

### Новый подход (слои)
```go
// cmd/server/main.go
db := database.InitDB()
repo := repository.NewPostRepository(db)
service := service.NewPostService(repo, ...)
handler := handlers.NewPostHandler(service, templates)
```

Преимущества:
- Понятная структура
- Легко тестировать
- Легко расширять

## Как запустить

```bash
# Старый способ (больше не работает)
go run main.go

# Новый способ
make run
# или
go run ./cmd/server
```

## Следующие шаги

1. **Добавить тесты** - для каждого слоя
2. **Добавить валидацию** - входных данных
3. **Добавить логирование** - структурированные логи
4. **Добавить кэширование** - для производительности
5. **Добавить миграции** - версионирование схемы БД

