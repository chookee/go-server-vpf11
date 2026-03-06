# 🚀 Быстрый старт

## 1. Установка зависимостей

```bash
cd E:\Personal\ZeroCode\VPf11\project\go-server
go mod tidy
```

## 2. Запуск сервера

```bash
go run cmd/server/main.go
```

Сервер запустится на `http://127.0.0.1:8080`

## 3. Проверка работы

```bash
# Health check
curl http://127.0.0.1:8080/health

# Создать пользователя
curl -X POST http://127.0.0.1:8080/add-user \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"Alice\"}"

# Получить всех пользователей
curl http://127.0.0.1:8080/users
```

## 4. Запуск тестов

```bash
# Все тесты
go test -v ./...

# Тесты с покрытием
go test -cover ./...

# Тесты с race detector
go test -race ./...
```

---

# ✅ Исправленные проблемы

## Ошибка компиляции в main.go

**Проблема:**
```
cmd\server\main.go:81:19: undefined: middleware.ContentTypeValidator
cmd\server\main.go:84:19: undefined: middleware.RequestLimit
```

**Причина:** Конфликт имён — `middleware` использовался и для `chi/v5/middleware`, и для локального пакета.

**Решение:**
```go
// Было (НЕПРАВИЛЬНО):
r.Use(middleware.ContentTypeValidator("application/json"))

// Стало (ПРАВИЛЬНО):
r.Use(appMiddleware.ContentTypeValidator("application/json"))
```

---

## Паника в TestGetUsers_EmptyList

**Проблема:**
```
panic: interface conversion: interface {} is nil, not []interface {}
```

**Причина:** При пустом списке пользователей JSON декодирует `users` как `nil`, а не пустой массив.

**Решение:**
```go
// Было (НЕПРАВИЛЬНО):
users := data["users"].([]interface{})

// Стало (ПРАВИЛЬНО):
dataData := data["data"].(map[string]interface{})
users, ok := dataData["users"].([]interface{})
if !ok {
    if dataData["users"] == nil {
        return // Пустой список - OK
    }
    t.Fatalf("expected users to be array, got %T", dataData["users"])
}
```

---

## Провал TestContentTypeValidator

**Проблема:**
```
--- FAIL: TestContentTypeValidator (0.00s)
    --- FAIL: TestContentTypeValidator/DELETE_without_Content-Type
        expected status 200, got 415
```

**Причина:** `DELETE` запросы не требовали Content-Type, но middleware их блокировал.

**Решение:**
```go
// Было (НЕПРАВИЛЬНО):
if r.Method == http.MethodGet || r.Method == http.MethodHead || 
    r.Method == http.MethodOptions {
    next.ServeHTTP(w, r)
    return
}

// Стало (ПРАВИЛЬНО):
if r.Method == http.MethodGet || r.Method == http.MethodHead || 
    r.Method == http.MethodOptions || r.Method == http.MethodDelete {
    next.ServeHTTP(w, r)
    return
}
```

---

# 📊 Статус тестов

| Пакет | Статус | Тестов |
|-------|--------|--------|
| `cmd/server` | ✅ OK | - |
| `internal/config` | ℹ️ No tests | - |
| `internal/database` | ✅ PASS | 12 |
| `internal/handlers` | ✅ PASS | 22 |
| `internal/middleware` | ✅ PASS | 8 |
| `internal/models` | ✅ PASS | 10 |
| `internal/utils` | ✅ PASS | 5 |
| **Итого** | **✅ PASS** | **57** |

---

# 🔧 Если что-то пошло не так

## Ошибка: "go: command not found"

**Решение:** Добавьте Go в PATH

1. Найдите где установлен Go (обычно `C:\Go\bin\go.exe`)
2. Добавьте `C:\Go\bin` в системную переменную PATH
3. Перезапустите терминал

Или используйте полный путь:
```bash
C:\Go\bin\go.exe mod tidy
```

## Ошибка: "module not found"

**Решение:**
```bash
go mod tidy
```

## Ошибка: "database is locked"

**Причина:** SQLite база заблокирована другим процессом.

**Решение:**
1. Закройте все запущенные экземпляры сервера
2. Удалите `test.db` если он заблокирован:
   ```bash
   del test.db
   ```
3. Запустите сервер заново

## Тесты падают с "connection refused"

**Причина:** Сервер не запущен.

**Решение:** Запустите сервер перед тестированием API:
```bash
go run cmd/server/main.go
```

---

# 📚 Дополнительные команды

## Сборка бинарного файла

```bash
# Для Windows
go build -o server.exe cmd/server/main.go

# Для Linux
GOOS=linux GOARCH=amd64 go build -o server cmd/server/main.go

# Для macOS
GOOS=darwin GOARCH=amd64 go build -o server cmd/server/main.go
```

## Запуск с переменными окружения

```bash
# Windows (PowerShell)
$env:SERVER_PORT="8080"
$env:LOG_LEVEL="debug"
go run cmd/server/main.go

# Linux/macOS
SERVER_PORT=8080 LOG_LEVEL=debug go run cmd/server/main.go
```

## Генерация HTML отчёта о покрытии

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Откроется браузер с визуальным отчётом.

## Бенчмарки

```bash
# Все бенчмарки
go test -bench=. ./...

# Только бенчмарки database
go test -bench=BenchmarkDB ./internal/database/...

# Бенчмарки с профилем
go test -bench=. -cpuprofile=cpu.prof ./...
```

---

# 📖 Полезные ссылки

- [Официальная документация Go](https://go.dev/doc/)
- [A Tour of Go](https://go.dev/tour/) — интерактивное обучение
- [Go by Example](https://gobyexample.com/) — примеры кода
- [Effective Go](https://go.dev/doc/effective_go) — лучшие практики
