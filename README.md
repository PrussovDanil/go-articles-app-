# Go Articles App

Приложение для управления статьями и пользователями на Go с использованием PostgreSQL.

## Описание

Это учебный проект, демонстрирующий работу с базой данных PostgreSQL в Go. Приложение реализует CRUD операции для пользователей и статей, используя паттерн Repository.

## Возможности

- ✅ Управление пользователями (создание, чтение, обновление, удаление)
- ✅ Управление статьями (создание, чтение, обновление, удаление)
- ✅ Публикация статей
- ✅ Подсчет просмотров статей
- ✅ Получение статей по автору
- ✅ Получение всех опубликованных статей
- ✅ Транзакции (создание статьи вместе с автором)
- ✅ Каскадное удаление (при удалении пользователя удаляются его статьи)

## Технологии

- **Go** 1.21+
- **PostgreSQL** 14+
- **database/sql** - стандартная библиотека для работы с БД
- **pq** - драйвер PostgreSQL для Go

## Структура проекта

```
go-articles-app/
├── db/                      # Пакет для работы с базой данных
│   └── connection.go        # Подключение к PostgreSQL
├── models/                  # Модели данных
│   ├── user.go             # Модель пользователя
│   └── article.go          # Модель статьи
├── repository/             # Репозитории для работы с БД
│   ├── user_repository.go   # Репозиторий пользователей
│   └── article_repository.go # Репозиторий статей
├── migrations/             # SQL миграции
│   ├── 001_create_users.sql
│   └── 002_create_articles.sql
├── main.go                 # Точка входа
├── Makefile               # Команды для сборки и запуска
└── README.md              # Этот файл
```

## Установка и запуск

### Предварительные требования

- Go 1.21 или выше
- PostgreSQL 14 или выше
- Make (опционально)

### 1. Клонировать репозиторий

```bash
git clone https://github.com/YOUR_USERNAME/go-articles-app.git
cd go-articles-app
```

### 2. Установить зависимости

```bash
go mod download
```

### 3. Настроить PostgreSQL

```bash
# Создать базу данных и пользователя
psql -U postgres

CREATE DATABASE go_article_app;
CREATE USER gouser WITH PASSWORD 'gopass';
GRANT ALL PRIVILEGES ON DATABASE go_article_app TO gouser;
\q
```

### 4. Применить миграции

```bash
# Выполнить миграции
psql -U gouser -d go_article_app -f migrations/001_create_users.sql
psql -U gouser -d go_article_app -f migrations/002_create_articles.sql
```

### 5. Настроить подключение к БД

Отредактируйте `main.go`, если нужно изменить параметры подключения:

```go
cfg := db.Config{
    Host:     "localhost",
    Port:     5432,
    User:     "gouser",
    Password: "gopass",
    DBName:   "go_article_app",
    SSLMode:  "disable",
}
```

### 6. Запустить приложение

```bash
# С помощью make
make run

# Или напрямую
go run main.go
```

## Использование Make

```bash
# Запустить приложение
make run

# Собрать бинарник
make build

# Запустить тесты
make test

# Очистить собранные файлы
make clean
```

## Примеры использования

### Создание пользователя

```go
user := &models.User{
    Email: "alice@example.com",
    Name:  "Alice",
}
err := userRepo.Create(ctx, user)
```

### Создание статьи

```go
article := &models.Article{
    Title:    "Introduction to Go",
    Content:  "Go is a statically typed language...",
    AuthorID: user.ID,
}
err := articleRepo.Create(ctx, article)
```

### Публикация статьи

```go
err := articleRepo.Publish(ctx, article.ID)
```

### Получение статей по автору

```go
articles, err := articleRepo.GetByAuthorID(ctx, authorID)
```

### Получение всех опубликованных статей

```go
publishedArticles, err := articleRepo.GetPublished(ctx)
```

## Структура базы данных

### Таблица `users`

| Поле       | Тип       | Описание                |
|------------|-----------|-------------------------|
| id         | SERIAL    | Первичный ключ          |
| email      | VARCHAR   | Email (уникальный)      |
| name       | VARCHAR   | Имя пользователя        |
| created_at | TIMESTAMP | Дата создания           |
| updated_at | TIMESTAMP | Дата последнего обновления |

### Таблица `articles`

| Поле       | Тип       | Описание                |
|------------|-----------|-------------------------|
| id         | SERIAL    | Первичный ключ          |
| title      | VARCHAR   | Заголовок статьи        |
| content    | TEXT      | Содержание статьи       |
| author_id  | INTEGER   | ID автора (FK на users) |
| published  | BOOLEAN   | Опубликована ли статья  |
| views      | INTEGER   | Количество просмотров   |
| created_at | TIMESTAMP | Дата создания           |
| updated_at | TIMESTAMP | Дата последнего обновления |

## API репозиториев

### UserRepository

- `Create(ctx, user)` - создать пользователя
- `GetByID(ctx, id)` - получить пользователя по ID
- `GetByEmail(ctx, email)` - получить пользователя по email
- `GetAll(ctx)` - получить всех пользователей
- `Update(ctx, user)` - обновить пользователя
- `Delete(ctx, id)` - удалить пользователя

### ArticleRepository

- `Create(ctx, article)` - создать статью
- `GetByID(ctx, id)` - получить статью по ID
- `GetByAuthorID(ctx, authorID)` - получить статьи автора
- `GetPublished(ctx)` - получить все опубликованные статьи
- `Update(ctx, article)` - обновить статью
- `Delete(ctx, id)` - удалить статью
- `Publish(ctx, id)` - опубликовать статью
- `IncrementViews(ctx, id)` - увеличить счетчик просмотров
- `CreateArticleWithAuthor(ctx, userName, userEmail, title, content)` - создать статью с автором в транзакции

## Примеры вывода

```
✅ Connected to PostgreSQL

📝 Creating users...
✅ Created user: Alice (alice@example.com)
✅ Created user: Bob (bob@example.com)
✅ Created user: Charlie (charlie@example.com)

📰 Creating articles with CreateArticleWithAuthor...
✅ Created article "Introduction to Go" by Alice
✅ Created article "PostgreSQL Basics" by Alice
✅ Created article "Web Development in Go" by Bob
✅ Created article "Docker for Beginners" by Bob
✅ Created article "Microservices Architecture" by Diana (new user created)

📢 Publishing articles...
✅ Published: "Introduction to Go"
✅ Published: "Web Development in Go"
✅ Published: "Docker for Beginners"

👁️  Incrementing views...
✅ "Introduction to Go" views: 0 → 5

📊 Statistics:
  - Total users: 4
  - Total articles: 5
  - Published articles: 3

🎉 All operations completed successfully!
```

## Особенности реализации

### Проверка на дубликаты
При создании пользователя проверяется уникальность email:
```go
if strings.Contains(err.Error(), "duplicate key") {
    return fmt.Errorf("user with email %s already exists", user.Email)
}
```

### Транзакции
Метод `CreateArticleWithAuthor` использует транзакцию для атомарного создания пользователя и статьи:
```go
tx, err := r.db.BeginTx(ctx, nil)
defer func() {
    if err != nil {
        tx.Rollback()
    }
}()
// ... операции
tx.Commit()
```

### Каскадное удаление
При удалении пользователя автоматически удаляются все его статьи благодаря `ON DELETE CASCADE` в БД.

## Лицензия

MIT

## Автор

Учебный проект для изучения Go и PostgreSQL
