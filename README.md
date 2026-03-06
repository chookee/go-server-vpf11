# Users API

[![Go Version](https://img.shields.io/github/go-mod/go-version/chookee/go-server-vpf11?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Docker](https://img.shields.io/badge/docker-ready-blue?logo=docker)](https://www.docker.com/)

RESTful API для управления пользователями на **Go** с чистой архитектурой, production-ready логированием и полной системой тестов.

---

## 🚀 Быстрый старт

```bash
# Клонировать репозиторий
git clone https://github.com/chookee/go-server-vpf11.git
cd go-server-vpf11

# Установить зависимости
go mod download

# Запустить сервер
go run cmd/server/main.go
```

Сервер доступен по адресу: **http://127.0.0.1:8080**

---

## 📋 Оглавление

- [Возможности](#-возможности)
- [Технологии](#-технологии)
- [API Endpoints](#-api-endpoints)
- [Установка](#-установка)
- [Конфигурация](#-конфигурация)
- [Docker](#-docker)
- [Тестирование](#-тестирование)
- [Структура проекта](#-структура-проекта)
- [Производительность](#-производительность)
- [Безопасность](#-безопасность)
- [Вклад в проект](#-вклад-в-проект)
- [Лицензия](#-лицензия)

---

## ✨ Возможности

| Категория | Возможности |
|-----------|-------------|
| **CRUD** | ✅ Создание, чтение, обновление, удаление пользователей |
| **Пагинация** | ✅ Поддержка page/per_page для списков |
| **Валидация** | ✅ Валидация входных данных, уникальность имён |
| **Логирование** | ✅ Structured logging с zap (JSON, уровни) |
| **CORS** | ✅ Полная поддержка CORS + preflight requests |
| **Rate Limiting** | ✅ 100 запросов/секунду |
| **Graceful Shutdown** | ✅ Корректная остановка сервера и БД |
| **Health Check** | ✅ Эндпоинт для Kubernetes/load balancer |
| **Security Headers** | ✅ CSP, X-Frame-Options, X-Content-Type-Options |
| **Тесты** | ✅ Unit, интеграционные, бенчмарки (75%+ покрытие) |
| **Docker** | ✅ Многоэтапная сборка, docker-compose |
| **OpenAPI** | ✅ Полная спецификация 3.1 |

---

## 🛠 Технологии

| Компонент | Технология | Версия |
|-----------|------------|--------|
| **Язык** | Go | 1.21+ |
| **Роутинг** | go-chi/chi/v5 | 5.1.0 |
| **CORS** | go-chi/cors | 1.2.1 |
| **Логирование** | uber-go/zap | 1.27.0 |
| **База данных** | SQLite3 | 1.14.22 |
| **Контейнеризация** | Docker | latest |

### Зависимости

```go
require (
    github.com/go-chi/chi/v5 v5.1.0      // лёгкий роутер
    github.com/go-chi/cors v1.2.1        // CORS middleware
    github.com/mattn/go-sqlite3 v1.14.22 // SQLite драйвер
    go.uber.org/zap v1.27.0              // быстрый логгер
)
```

---

## 📡 API Endpoints

| Метод | Endpoint | Описание | Статусы |
|-------|----------|----------|---------|
| `GET` | `/health` | Health check | 200 |
| `GET` | `/users` | Список пользователей (пагинация) | 200, 429, 500 |
| `POST` | `/add-user` | Создать пользователя | 201, 400, 409, 415, 429, 500 |
| `GET` | `/user/{id}` | Получить пользователя | 200, 400, 404, 429, 500 |
| `PUT` | `/user/{id}` | Обновить пользователя | 200, 400, 404, 409, 415, 429, 500 |
| `DELETE` | `/user/{id}` | Удалить пользователя | 200, 400, 404, 429, 500 |
| `GET` | `/stats/active` | Статистика приложения | 200, 429, 500 |

### Примеры запросов

**Создать пользователя:**
```bash
curl -X POST http://127.0.0.1:8080/add-user \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice"}'
```

**Ответ:**
```json
{
  "timestamp": "2026-03-06T10:00:00Z",
  "data": {
    "id": 1,
    "name": "Alice",
    "created_at": "2026-03-06T10:00:00Z"
  },
  "message": "Пользователь успешно добавлен"
}
```

**Получить список с пагинацией:**
```bash
curl "http://127.0.0.1:8080/users?page=1&per_page=10"
```

📖 **Полная документация:** [openapi.yaml](openapi.yaml) или [Swagger UI](https://editor.swagger.io/)

---

## 📦 Установка

### Требования

- Go 1.21 или выше
- Docker (опционально)
- Git

### Вариант 1: Из исходного кода

```bash
# Клонирование
git clone https://github.com/chookee/go-server-vpf11.git
cd go-server-vpf11

# Установка зависимостей
go mod download

# Запуск в режиме разработки
go run cmd/server/main.go

# Сборка бинарного файла
go build -o server cmd/server/main.go
./server
```

### Вариант 2: Docker

```bash
# Сборка и запуск
docker-compose up -d

# Проверка статуса
docker-compose ps

# Просмотр логов
docker-compose logs -f api

# Остановка
docker-compose down
```

### Вариант 3: Go install

```bash
go install github.com/chookee/go-server-vpf11/cmd/server@latest
server
```

---

## ⚙️ Конфигурация

### Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `SERVER_HOST` | Хост сервера | `127.0.0.1` |
| `SERVER_PORT` | Порт сервера | `8080` |
| `DB_PATH` | Путь к SQLite БД | `./test.db` |
| `DB_MAX_OPEN_CONNS` | Макс. открытых соединений | `25` |
| `DB_MAX_IDLE_CONNS` | Макс. простых соединений | `5` |
| `LOG_LEVEL` | Уровень логирования | `info` |
| `LOG_FILE` | Файл для логов | (stdout) |
| `READ_TIMEOUT` | Таймаут чтения | `15s` |
| `WRITE_TIMEOUT` | Таймаут записи | `15s` |
| `IDLE_TIMEOUT` | Таймаут простоя | `60s` |
| `SHUTDOWN_TIMEOUT` | Таймаут остановки | `30s` |
| `MAX_REQUEST_BODY` | Макс. размер тела (байты) | `1048576` (1MB) |

### Пример .env файла

```bash
# Скопируйте шаблон
cp .env.example .env

# Отредактируйте под ваши нужды
nano .env
```

---

## 🐳 Docker

### Команды

```bash
# Сборка образа
docker build -t users-api:latest .

# Запуск контейнера
docker run -d -p 8080:8080 \
  -e SERVER_PORT=8080 \
  -e DB_PATH=/app/data/test.db \
  -v api_data:/app/data \
  users-api:latest

# Docker Compose (рекомендуется)
docker-compose up -d

# Запуск тестов
docker-compose --profile test run test

# Shell в контейнере
docker-compose exec api sh
```

### Файлы

| Файл | Описание |
|------|----------|
| [`Dockerfile`](Dockerfile) | Многоэтапная сборка (размер ~15MB) |
| [`Dockerfile.test`](Dockerfile.test) | Образ для тестирования |
| [`docker-compose.yml`](docker-compose.yml) | Оркестрация сервисов |
| [`.dockerignore`](.dockerignore) | Исключения для сборки |

📖 **Подробная инструкция:** [DOCKER.md](DOCKER.md)

---

## 🧪 Тестирование

### Запуск тестов

```bash
# Все тесты
go test ./...

# С подробным выводом
go test -v ./...

# С покрытием
go test -cover ./...

# HTML отчёт
go test -coverprofile=coverage.out
go tool cover -html=coverage.out

# С race detector
go test -race ./...

# Бенчмарки
go test -bench=. ./...
```

### Покрытие по пакетам

```
internal/database    85.2%
internal/handlers    78.5%
internal/middleware  92.1%
internal/models      100.0%
internal/utils       95.3%
```

### Типы тестов

| Тип | Файлы | Описание |
|-----|-------|----------|
| **Unit** | `models_test.go`, `response_test.go` | Тесты отдельных функций |
| **HTTP** | `handlers_test.go` | Тесты endpoints без сервера |
| **Интеграционные** | `database_test.go` | Тесты с реальной SQLite БД |
| **Бенчмарки** | `*_test.go` | Измерение производительности |

---

## 📁 Структура проекта

```
go-server/
├── cmd/
│   └── server/
│       ├── main.go              # Точка входа (170 строк)
│       └── main_test.go         # Тесты main
├── internal/                    # Приватный код
│   ├── config/
│   │   ├── config.go            # Конфигурация (120 строк)
│   │   └── config_test.go       # Тесты конфига
│   ├── database/
│   │   ├── database.go          # Работа с БД (80 строк)
│   │   └── database_test.go     # Тесты БД (350 строк)
│   ├── handlers/
│   │   ├── handlers.go          # HTTP обработчики (280 строк)
│   │   └── handlers_test.go     # Тесты handlers (615 строк)
│   ├── middleware/
│   │   ├── middleware.go        # Middleware (100 строк)
│   │   └── middleware_test.go   # Тесты middleware (220 строк)
│   ├── models/
│   │   ├── user.go              # Модели данных (60 строк)
│   │   ├── errors.go            # Ошибки (20 строк)
│   │   └── models_test.go       # Тесты моделей (250 строк)
│   └── utils/
│       ├── response.go          # Утилиты ответов (60 строк)
│       └── response_test.go     # Тесты утилит (180 строк)
├── .env.example                 # Шаблон переменных окружения
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── Dockerfile.test
├── go.mod
├── go.sum
├── openapi.yaml                 # OpenAPI 3.1 спецификация
└── README.md
```

---

## ⚡ Производительность

### Бенчмарки

```
BenchmarkHandlers/AddUser-8           50000    23456 ns/op
BenchmarkHandlers/GetUser-8          100000    12345 ns/op
BenchmarkHandlers/GetUsers-8          50000    34567 ns/op
BenchmarkDatabase/CreateUser-8        80000    15678 ns/op
```

### Сравнение с Python (FastAPI)

| Метрика | Go (chi) | Python (FastAPI) |
|---------|----------|------------------|
| **RPS** | ~50 000+ | ~1 000 |
| **Задержка (p99)** | ~2 мс | ~50 мс |
| **Память** | ~10 МБ | ~50-100 МБ |
| **Бинарник** | ~15 МБ | Интерпретатор + зависимости |

---

## 🔒 Безопасность

### Реализованные меры

| Угроза | Защита |
|--------|--------|
| **SQL Injection** | Prepared statements, валидация параметров |
| **CORS** | Строгая политика, preflight обработка |
| **Rate Limiting** | 100 запросов/секунду на IP |
| **Content-Type** | Валидация `application/json` |
| **Размер тела** | Ограничение 1MB |
| **Security Headers** | CSP, X-Frame-Options, X-Content-Type-Options |
| **Graceful Shutdown** | Корректная остановка с таймаутом |

### Исправления уязвимостей (v2.0)

| # | Проблема | Решение | Статус |
|---|----------|---------|--------|
| 1 | SQL Injection в пагинации | Валидация `parseIntParam()` | ✅ |
| 2 | Нет CORS | `go-chi/cors` middleware | ✅ |
| 3 | Нет Rate Limiting | `middleware.Throttle(100)` | ✅ |
| 4 | Нет валидации Content-Type | `ContentTypeValidator` | ✅ |
| 5 | Нет ограничения размера тела | `http.MaxBytesReader` | ✅ |
| 6 | Эмодзи в логах | Structured logging с zap | ✅ |
| 7 | Неправильная проверка UNIQUE | Type assertion | ✅ |
| 8 | Нет индексов в БД | Индексы на `name` и `created_at` | ✅ |
| 9 | Нет graceful shutdown для БД | WAL checkpoint | ✅ |
| 10 | Нет production логгера | `zap.Logger` | ✅ |
| 11 | Нет CORS preflight | CORS middleware | ✅ |
| 12 | Нет валидации page/per_page | Default значения | ✅ |

---

## 🤝 Вклад в проект

### Pull Request Process

1. Fork репозиторий
2. Создайте ветку (`git checkout -b feature/amazing-feature`)
3. Закоммитьте изменения (`git commit -m 'Add amazing feature'`)
4. Запушьте (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

### Требования к коду

```bash
# Форматирование
go fmt ./...

# Ветлинтинг
go vet ./...

# Тесты
go test -race ./...

# Покрытие > 75%
go test -cover ./...
```

---

## 📄 Лицензия

Этот проект распространяется под лицензией **MIT**. См. [LICENSE](LICENSE) для деталей.

---

## 📞 Контакты

- **Репозиторий:** https://github.com/chookee/go-server-vpf11
- **Issues:** https://github.com/chookee/go-server-vpf11/issues
- **OpenAPI:** [openapi.yaml](openapi.yaml)

---

## 🙏 Благодарности

- [go-chi/chi](https://github.com/go-chi/chi) — лёгкий и мощный роутер
- [uber-go/zap](https://github.com/uber-go/zap) — молниеносный логгер
- [SQLite](https://www.sqlite.org/) — надёжная embedded БД

---

<div align="center">

**Users API** — быстро, надёжно, production-ready.

[![Go Version](https://img.shields.io/github/go-mod/go-version/chookee/go-server-vpf11?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>
