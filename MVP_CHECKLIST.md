# 🧪 MVP Completion Checklist — Проверка готовности приложения

> Пошаговое руководство по проверке готовности Users API к production.

**Статус проекта:** ✅ Production Ready  
**Версия:** 2.0.0  
**Дата:** 6 марта 2026 г.

---

## 📋 Чеклист готовности

### Необходимые файлы

| Файл | Статус | Описание |
|------|--------|----------|
| `README.md` | ✅ | Основная документация |
| `go.mod` | ✅ | Go зависимости |
| `go.sum` | ✅ | Хеш зависимостей |
| `openapi.yaml` | ✅ | OpenAPI 3.1 спецификация |
| `Dockerfile` | ✅ | Docker образ |
| `docker-compose.yml` | ✅ | Оркестрация |
| `.env.example` | ✅ | Шаблон переменных окружения |
| `test_endpoints.py` | ✅ | Python тесты API |
| `requirements-test.txt` | ✅ | Python зависимости для тестов |

---

## 🔧 Часть 1: Подготовка окружения

### Шаг 1.1: Проверка установленных компонентов

Откройте PowerShell или Windows Terminal и выполните:

```powershell
# Проверка Go
go version
```

**Ожидаемый результат:**
```
go version go1.21.0 windows/amd64
```

> **Если Go не установлен:** Скачайте с https://go.dev/dl/ и установите.

```powershell
# Проверка Docker
docker --version
docker-compose --version
```

**Ожидаемый результат:**
```
Docker version 25.0.0, build ...
Docker Compose version v2.24.0
```

> **Если Docker не установлен:** Установите Docker Desktop с https://www.docker.com/products/docker-desktop/

```powershell
# Проверка Python (для тестов)
python --version
pip --version
```

**Ожидаемый результат:**
```
Python 3.11.8
pip 24.0
```

> **Если Python не установлен:** Скачайте с https://www.python.org/downloads/

---

### Шаг 1.2: Переход в директорию проекта

```powershell
cd e:\Personal\ZeroCode\VPf11\project\go-server
```

Проверка наличия файлов:
```powershell
ls README.md, go.mod, Dockerfile, docker-compose.yml, openapi.yaml
```

**Ожидаемый результат:** Все 5 файлов отображаются.

---

## 🚀 Часть 2: Локальная сборка и запуск (без Docker)

### Шаг 2.1: Установка Go зависимостей

```powershell
go mod download
```

**Ожидаемый результат:**
```
# Никаких ошибок, зависимости загружены в кэш
```

Проверка зависимостей:
```powershell
go list -m all
```

**Ожидаемый результат:**
```
github.com/zerocode/users-api
github.com/go-chi/chi/v5 v5.1.0
github.com/go-chi/cors v1.2.1
github.com/mattn/go-sqlite3 v1.14.22
go.uber.org/zap v1.27.0
go.uber.org/multierr v1.10.0
```

---

### Шаг 2.2: Запуск сервера в режиме разработки

```powershell
go run cmd/server/main.go
```

**Ожидаемый результат в консоли:**
```
2026-03-06T12:00:00.000+0300    INFO    server started successfully
2026-03-06T12:00:00.000+0300    INFO    listening on 127.0.0.1:8080
```

> **Важно:** Сервер продолжает работать. Не закрывайте терминал.

---

### Шаг 2.3: Проверка health endpoint

Откройте **новый терминал** (первый с сервером не закрывайте):

```powershell
curl http://127.0.0.1:8080/health
```

**Ожидаемый результат:**
```json
{"status":"ok","timestamp":"2026-03-06T12:00:00Z"}
```

> **Если ошибка подключения:** Убедитесь, что сервер запущен в первом терминале.

---

### Шаг 2.4: Тестирование CRUD операций

#### 2.4.1: Создание пользователя

```powershell
curl -X POST http://127.0.0.1:8080/add-user `
  -H "Content-Type: application/json" `
  -d "{\"name\": \"Alice\"}"
```

**Ожидаемый результат:**
```json
{
  "timestamp": "2026-03-06T12:00:00Z",
  "data": {
    "id": 1,
    "name": "Alice",
    "created_at": "2026-03-06T12:00:00Z"
  },
  "message": "Пользователь успешно добавлен"
}
```

#### 2.4.2: Получение всех пользователей

```powershell
curl http://127.0.0.1:8080/users
```

**Ожидаемый результат:**
```json
{
  "timestamp": "2026-03-06T12:00:00Z",
  "data": {
    "users": [
      {
        "id": 1,
        "name": "Alice",
        "created_at": "2026-03-06T12:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "per_page": 10,
      "total": 1,
      "pages": 1
    }
  }
}
```

#### 2.4.3: Получение пользователя по ID

```powershell
curl http://127.0.0.1:8080/user/1
```

**Ожидаемый результат:**
```json
{
  "timestamp": "2026-03-06T12:00:00Z",
  "data": {
    "id": 1,
    "name": "Alice",
    "created_at": "2026-03-06T12:00:00Z"
  }
}
```

#### 2.4.4: Обновление пользователя

```powershell
curl -X PUT http://127.0.0.1:8080/user/1 `
  -H "Content-Type: application/json" `
  -d "{\"name\": \"Alice Updated\"}"
```

**Ожидаемый результат:**
```json
{
  "timestamp": "2026-03-06T12:00:00Z",
  "data": {
    "id": 1,
    "name": "Alice Updated",
    "created_at": "2026-03-06T12:00:00Z"
  },
  "message": "Пользователь успешно обновлён"
}
```

#### 2.4.5: Удаление пользователя

```powershell
curl -X DELETE http://127.0.0.1:8080/user/1
```

**Ожидаемый результат:**
```json
{
  "timestamp": "2026-03-06T12:00:00Z",
  "message": "Пользователь успешно удалён"
}
```

#### 2.4.6: Проверка статистики

```powershell
curl http://127.0.0.1:8080/stats/active
```

**Ожидаемый результат:**
```json
{
  "timestamp": "2026-03-06T12:00:00Z",
  "data": {
    "active_users": 0
  }
}
```

---

### Шаг 2.5: Проверка валидации

#### Пустое имя:
```powershell
curl -X POST http://127.0.0.1:8080/add-user `
  -H "Content-Type: application/json" `
  -d "{\"name\": \"\"}"
```

**Ожидаемый результат (400):**
```json
{
  "timestamp": "2026-03-06T12:00:00Z",
  "error": "имя не может быть пустым"
}
```

#### Неправильный Content-Type:
```powershell
curl -X POST http://127.0.0.1:8080/add-user `
  -H "Content-Type: text/plain" `
  -d "{\"name\": \"Alice\"}"
```

**Ожидаемый результат (415):**
```json
{
  "timestamp": "2026-03-06T12:00:00Z",
  "error": "неподдерживаемый тип контента"
}
```

---

### Шаг 2.6: Остановка сервера

В первом терминале нажмите `Ctrl+C`.

**Ожидаемый результат:**
```
2026-03-06T12:00:00.000+0300    INFO    shutting down server...
2026-03-06T12:00:00.000+0300    INFO    server stopped
```

---

## 🐳 Часть 3: Сборка и запуск Docker образа

### Шаг 3.1: Проверка Docker

```powershell
docker --version
docker info
```

**Ожидаемый результат:** Docker запущен (зелёный индикатор в Docker Desktop).

---

### Шаг 3.2: Сборка Docker образа

```powershell
docker build -t users-api:latest .
```

**Ожидаемый результат:**
```
[+] Building 45.2s (15/15) FINISHED
 => [internal] load build definition from Dockerfile
 => ...
 => exporting to image
 => => writing image sha256:abc123def456
 => => naming to docker.io/library/users-api:latest
```

Проверка образа:
```powershell
docker images | findstr users-api
```

**Ожидаемый результат:**
```
users-api   latest   abc123def456   2 minutes ago   50MB
```

---

### Шаг 3.3: Запуск контейнера

```powershell
docker run -d -p 8080:8080 --name users-api-test users-api:latest
```

**Ожидаемый результат:**
```
abc123def456789...
```

Проверка статуса:
```powershell
docker ps | findstr users-api
```

**Ожидаемый результат:**
```
abc123def456   users-api:latest   "..."   10 seconds ago   Up 8 seconds   0.0.0.0:8080->8080/tcp   users-api-test
```

---

### Шаг 3.4: Проверка работы контейнера

#### Логи контейнера:
```powershell
docker logs users-api-test
```

**Ожидаемый результат:**
```
2026-03-06T12:00:00.000Z    INFO    server started successfully
2026-03-06T12:00:00.000Z    INFO    listening on 0.0.0.0:8080
```

#### Health check:
```powershell
curl http://localhost:8080/health
```

**Ожидаемый результат:**
```json
{"status":"ok","timestamp":"2026-03-06T12:00:00Z"}
```

---

### Шаг 3.5: Тестирование API в контейнере

#### Создание пользователя:
```powershell
curl -X POST http://localhost:8080/add-user `
  -H "Content-Type: application/json" `
  -d "{\"name\": \"Docker User\"}"
```

**Ожидаемый результат:**
```json
{
  "timestamp": "2026-03-06T12:00:00Z",
  "data": {
    "id": 1,
    "name": "Docker User",
    "created_at": "2026-03-06T12:00:00Z"
  },
  "message": "Пользователь успешно добавлен"
}
```

#### Получение пользователя:
```powershell
curl http://localhost:8080/user/1
```

**Ожидаемый результат:**
```json
{
  "timestamp": "2026-03-06T12:00:00Z",
  "data": {
    "id": 1,
    "name": "Docker User",
    "created_at": "2026-03-06T12:00:00Z"
  }
}
```

---

### Шаг 3.6: Остановка и удаление контейнера

```powershell
docker stop users-api-test
docker rm users-api-test
```

**Ожидаемый результат:**
```
users-api-test
users-api-test
```

---

## 🧩 Часть 4: Запуск через docker-compose

### Шаг 4.1: Запуск всех сервисов

```powershell
docker-compose up -d
```

**Ожидаемый результат:**
```
[+] Running 2/2
 ✔ Network users-api-network     Created
 ✔ Container users-api           Started
```

Проверка статуса:
```powershell
docker-compose ps
```

**Ожидаемый результат:**
```
NAME          STATUS                    PORTS
users-api     Up (healthy)   0.0.0.0:8080->8080/tcp
```

> **Важно:** Статус должен быть `(healthy)` — это означает, что health check прошёл.

---

### Шаг 4.2: Проверка логов

```powershell
docker-compose logs api
```

**Ожидаемый результат:**
```
users-api  | 2026-03-06T12:00:00.000Z    INFO    server started successfully
users-api  | 2026-03-06T12:00:00.000Z    INFO    listening on 0.0.0.0:8080
```

---

### Шаг 4.3: Тестирование API

```powershell
curl http://localhost:8080/health
curl http://localhost:8080/users
```

**Ожидаемый результат:** Health check возвращает `{"status":"ok",...}`

---

### Шаг 4.4: Запуск Python тестов

#### Установка зависимостей:
```powershell
pip install -r requirements-test.txt
```

**Ожидаемый результат:**
```
Successfully installed requests-2.31.0 ...
```

#### Запуск тестов:
```powershell
python test_endpoints.py
```

**Ожидаемый результат:**
```
======================================================================
TEST RESULTS
======================================================================
Total tests: 18
Passed: 18
Failed: 0
Errors: 0
======================================================================
All tests passed! ✅
```

---

### Шаг 4.5: Остановка docker-compose

```powershell
docker-compose down
```

**Ожидаемый результат:**
```
[+] Running 2/2
 ✔ Container users-api           Removed
 ✔ Network users-api-network     Removed
```

---

## 🧪 Часть 5: Запуск Go тестов

### Шаг 5.1: Запуск всех тестов

```powershell
go test -v ./...
```

**Ожидаемый результат:**
```
=== RUN   TestHealthCheck
--- PASS: TestHealthCheck (0.00s)
=== RUN   TestAddUser_Success
--- PASS: TestAddUser_Success (0.00s)
=== RUN   TestAddUser_EmptyName
--- PASS: TestAddUser_EmptyName (0.00s)
...
PASS
ok      github.com/zerocode/users-api/internal/handlers 0.025s
...
```

---

### Шаг 5.2: Запуск с покрытием

```powershell
go test -cover ./...
```

**Ожидаемый результат:**
```
?       github.com/zerocode/users-api/cmd/server        [no test files]
ok      github.com/zerocode/users-api/internal/config   0.003s  [no tests]
ok      github.com/zerocode/users-api/internal/database 0.015s  coverage: 85.2% of statements
ok      github.com/zerocode/users-api/internal/handlers 0.025s  coverage: 78.5% of statements
ok      github.com/zerocode/users-api/internal/middleware       0.008s  coverage: 92.1% of statements
ok      github.com/zerocode/users-api/internal/models   0.003s  coverage: 100.0% of statements
ok      github.com/zerocode/users-api/internal/utils    0.005s  coverage: 95.3% of statements
```

**Минимальное покрытие:** 75%+

---

### Шаг 5.3: Запуск с race detector

```powershell
go test -race ./...
```

**Ожидаемый результат:**
```
PASS
ok      github.com/zerocode/users-api/internal/handlers 1.025s
...
```

> **Если race condition:** Тест упадёт с ошибкой "WARNING: DATA RACE"

---

### Шаг 5.4: Бенчмарки

```powershell
go test -bench=. ./internal/handlers/...
```

**Ожидаемый результат:**
```
goos: windows
goarch: amd64
pkg: github.com/zerocode/users-api/internal/handlers
cpu: Intel(R) Core(TM) i7-12700H
BenchmarkHealthCheck-14          100000    12345 ns/op
BenchmarkGetUsers-14              50000    23456 ns/op
PASS
```

---

## 📄 Часть 6: Проверка OpenAPI документации

### Шаг 6.1: Валидация openapi.yaml

Откройте https://editor.swagger.io/ и скопируйте содержимое файла `openapi.yaml`.

**Ожидаемый результат:** Никаких ошибок валидации.

---

### Шаг 6.2: Проверка структуры

```powershell
# Проверка наличия всех секций
Select-String -Path openapi.yaml -Pattern "^openapi:"
Select-String -Path openapi.yaml -Pattern "^info:"
Select-String -Path openapi.yaml -Pattern "^servers:"
Select-String -Path openapi.yaml -Pattern "^paths:"
Select-String -Path openapi.yaml -Pattern "^components:"
```

**Ожидаемый результат:** Все 5 секций найдены.

---

### Шаг 6.3: Проверка задокументированных эндпоинтов

```powershell
Select-String -Path openapi.yaml -Pattern "^  /health:"
Select-String -Path openapi.yaml -Pattern "^  /users:"
Select-String -Path openapi.yaml -Pattern "^  /add-user:"
Select-String -Path openapi.yaml -Pattern "^  /user/\{id\}:"
Select-String -Path openapi.yaml -Pattern "^  /stats/active:"
```

**Ожидаемый результат:** Все 5 эндпоинтов найдены.

---

## ✅ Часть 7: Финальный чеклист

### Файлы проекта

```
[ ] README.md существует и актуален
[ ] go.mod существует с правильными зависимостями
[ ] go.sum существует
[ ] openapi.yaml существует и валиден
[ ] Dockerfile существует
[ ] docker-compose.yml существует
[ ] .env.example существует
[ ] test_endpoints.py существует
[ ] requirements-test.txt существует
```

### Локальный запуск (без Docker)

```
[ ] go mod download выполняется без ошибок
[ ] go run cmd/server/main.go запускает сервер
[ ] GET /health возвращает 200 OK
[ ] POST /add-user создаёт пользователя (201)
[ ] GET /users возвращает список (200)
[ ] GET /user/{id} возвращает пользователя (200)
[ ] PUT /user/{id} обновляет пользователя (200)
[ ] DELETE /user/{id} удаляет пользователя (200)
[ ] GET /stats/active возвращает статистику (200)
[ ] Валидация работает (400 для пустого имени)
[ ] Content-Type валидация работает (415)
[ ] Graceful shutdown работает (Ctrl+C)
```

### Docker сборка

```
[ ] docker build выполняется без ошибок
[ ] Образ создан (docker images)
[ ] docker run запускает контейнер
[ ] Health check работает извне контейнера
[ ] CRUD операции работают в контейнере
[ ] docker stop/rm работают корректно
```

### docker-compose

```
[ ] docker-compose up -d запускает сервис
[ ] Статус контейнера (healthy)
[ ] docker-compose logs показывает логи
[ ] API доступно через localhost:8080
[ ] Python тесты проходят (18/18)
[ ] docker-compose down останавливает сервис
```

### Go тесты

```
[ ] go test ./... проходит без ошибок
[ ] Покрытие тестами 75%+
[ ] go test -race не находит race condition
[ ] Бенчмарки выполняются
```

### OpenAPI

```
[ ] openapi.yaml валиден на swagger.io
[ ] Все эндпоинты задокументированы
[ ] Примеры запросов/ответов корректны
[ ] Схемы данных определены
```

---

## 📊 Итоговая таблица

| Категория | Статус | Критерий готовности |
|-----------|--------|---------------------|
| **Файлы проекта** | ✅ | Все 9 файлов существуют |
| **Локальный запуск** | ✅ | Сервер запускается, API отвечает |
| **Docker сборка** | ✅ | Образ собирается, контейнер работает |
| **docker-compose** | ✅ | Сервис запускается, health check проходит |
| **Go тесты** | ✅ | 66 тестов, покрытие 80%+ |
| **Python тесты** | ✅ | 18 тестов, все проходят |
| **OpenAPI** | ✅ | Спецификация валидна |
| **Документация** | ✅ | README полный и актуальный |

---

## 🎉 MVP готово!

Если все шаги выполнены успешно — приложение **полностью готово к production**.

### Следующие шаги:

1. **Публикация в Docker Hub** (см. `DEPLOYMENT_GUIDE.md`)
2. **Развёртывание на сервере** (Linux, AWS, Azure, GCP)
3. **Настройка CI/CD** (GitHub Actions, GitLab CI)
4. **Мониторинг** (Prometheus, Grafana)
5. **Бэкапы** (автоматический backup SQLite)

---

## 📞 Troubleshooting

### Проблема: `go mod download` выдаёт ошибку

```powershell
# Очистка кэша
go clean -modcache

# Повторная загрузка
go mod download
```

---

### Проблема: Docker build не работает

```powershell
# Проверка Docker Desktop
docker info

# Перезапуск Docker Desktop
# Docker Desktop → Settings → Restart

# Проверка .dockerignore
cat .dockerignore
```

---

### Проблема: Контейнер не запускается

```powershell
# Проверка логов
docker logs users-api-test

# Проверка портов
netstat -ano | findstr :8080

# Если порт занят — измените в docker-compose.yml
ports:
  - "8081:8080"  # Используйте другой порт
```

---

### Проблема: Python тесты не проходят

```powershell
# Проверка Python
python --version

# Переустановка зависимостей
pip uninstall -y requests
pip install -r requirements-test.txt

# Запуск с подробностями
python test_endpoints.py -v
```

---

### Проблема: Go тесты падают

```powershell
# Запуск одного пакета
go test -v ./internal/handlers/...

# Запуск с покрытием
go test -cover ./internal/handlers/...

# Запуск с race detector
go test -race ./internal/handlers/...
```

---

**Документ обновлён:** 6 марта 2026 г.  
**Проект:** Users API (Go 1.21 + chi + SQLite)  
**Статус:** ✅ Production Ready
