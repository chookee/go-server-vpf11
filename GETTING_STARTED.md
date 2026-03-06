# 🚀 Быстрый старт — Users API

## 📋 Варианты запуска

| Способ | Сложность | Время | Для чего |
|--------|-----------|-------|----------|
| **Docker Compose** | ⭐ Легко | 2 мин | Production, тесты |
| **Python скрипт** | ⭐⭐ Средне | 5 мин | Локальная разработка |
| **Go сборка** | ⭐⭐⭐ Сложно | 10 мин | Отладка, разработка |

---

## 🐳 Способ 1: Docker Compose (Рекомендуется)

### Требования
- Docker Desktop (Windows/macOS) или docker.io (Linux)
- Docker Compose

### Запуск

```bash
# 1. Перейдите в директорию проекта
cd E:\Personal\ZeroCode\VPf11\project\go-server

# 2. Запустите сервер
docker-compose up -d --build

# 3. Проверьте что работает
docker-compose ps

# 4. Откройте в браузере
http://localhost:8080/health
```

### Тестирование

```bash
# Вариант A: Python скрипт
pip install -r requirements-test.txt
python test_endpoints.py

# Вариант B: Docker тесты
docker-compose --profile test run test

# Вариант C: Ручные запросы
curl http://localhost:8080/health
curl http://localhost:8080/users
```

### Остановка

```bash
docker-compose down
```

---

## 🐍 Способ 2: Python скрипт для тестов

### Требования
- Python 3.8+
- Запущенный сервер (Docker или Go)

### Установка

```bash
# Установка зависимостей
pip install -r requirements-test.txt
```

### Запуск тестов

```bash
# Сервер должен быть запущен!
python test_endpoints.py
```

### Пример вывода

```
══════════════════════════════════════════════════════════════════
                    🧪 ТЕСТИРОВАНИЕ API ENDPOINTS                 
══════════════════════════════════════════════════════════════════

Базовый URL: http://127.0.0.1:8080
Таймаут:     5s

✓ API доступен

Запланировано тестов: 18

Тест 1/18: Проверка работоспособности сервиса
[✓ PASS] Health Check
       GET /health
       Статус: 200 (ожидался 200)
       Время: 12.45ms

...

══════════════════════════════════════════════════════════════════
                           СВОДКА                                 
══════════════════════════════════════════════════════════════════

Всего тестов:    18
Пройдено:        18
Провалено:       0
Среднее время:   15.23ms

Успешность:      100.0%
[████████████████████████████████████████]

✅ Все тесты пройдены!
```

---

## 🔧 Способ 3: Локальная сборка Go

### Требования
- Go 1.21+
- GCC (для SQLite)

### Установка Go

**Windows:**
1. Скачайте с https://go.dev/dl/
2. Установите
3. Перезапустите терминал

**Linux:**
```bash
sudo apt install golang-go
```

### Сборка и запуск

```bash
# 1. Перейдите в проект
cd E:\Personal\ZeroCode\VPf11\project\go-server

# 2. Загрузите зависимости
go mod download

# 3. Запустите сервер
go run cmd/server/main.go

# 4. В другом терминале проверьте
curl http://localhost:8080/health
```

### Сборка бинарника

```bash
# Windows
go build -o server.exe cmd/server/main.go

# Linux
go build -o server cmd/server/main.go
```

---

## ✅ Проверка работы

### Health Check

```bash
curl http://localhost:8080/health
```

**Ожидаемый ответ:**
```json
{
  "timestamp": "2026-03-05T12:00:00Z",
  "data": {
    "status": "healthy"
  }
}
```

### Создание пользователя

```bash
curl -X POST http://localhost:8080/add-user \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"Alice\"}"
```

**Ожидаемый ответ:**
```json
{
  "timestamp": "2026-03-05T12:00:00Z",
  "data": {
    "id": 1,
    "name": "Alice"
  },
  "message": "Пользователь успешно добавлен"
}
```

### Получение списка

```bash
curl http://localhost:8080/users
```

---

## 🎯 Сценарии использования

### Сценарий 1: Локальная разработка

```bash
# Терминал 1: Запуск сервера
go run cmd/server/main.go

# Терминал 2: Тесты
python test_endpoints.py

# Терминал 3: Изменение кода...
```

### Сценарий 2: Production тестирование

```bash
# Сборка и запуск
docker-compose up -d --build

# Запуск тестов
docker-compose --profile test run test

# Просмотр логов
docker-compose logs -f api

# Остановка
docker-compose down
```

### Сценарий 3: CI/CD

```bash
# В CI пайплайне
docker build -t users-api:latest .
docker-compose up -d
sleep 5
docker-compose --profile test run test
docker-compose down
```

---

## 🔧 Переменные окружения

### Для Docker

Создайте `.env` файл:

```bash
API_PORT=8080
LOG_LEVEL=info
DB_PATH=/app/data/test.db
READ_TIMEOUT=30s
WRITE_TIMEOUT=30s
```

### Для локального запуска

**Windows (PowerShell):**
```powershell
$env:SERVER_PORT="8080"
$env:LOG_LEVEL="debug"
go run cmd/server/main.go
```

**Linux/macOS:**
```bash
SERVER_PORT=8080 LOG_LEVEL=debug go run cmd/server/main.go
```

---

## 📊 Команды для управления

| Команда | Описание |
|---------|----------|
| `docker-compose up -d` | Запуск в фоне |
| `docker-compose down` | Остановка |
| `docker-compose ps` | Статус |
| `docker-compose logs -f` | Логи |
| `python test_endpoints.py` | Тесты API |
| `go run cmd/server/main.go` | Запуск сервера |
| `go build -o server` | Сборка |
| `go test ./...` | Go тесты |

---

## 🐛 Troubleshooting

### Сервер не запускается

```bash
# Проверьте что порт свободен
netstat -ano | findstr :8080

# Проверьте логи Docker
docker-compose logs api

# Пересоберите
docker-compose down
docker-compose up -d --build
```

### Тесты не проходят

```bash
# Проверите что сервер запущен
curl http://localhost:8080/health

# Проверьте Python зависимости
pip install -r requirements-test.txt

# Запустите с отладкой
python -u test_endpoints.py
```

### Ошибки Go

```bash
# Очистите кэш
go clean -cache -modcache

# Переустановите зависимости
go mod tidy
```

---

## 📚 Документация

| Файл | Описание |
|------|----------|
| [README.md](README.md) | Основная документация |
| [DOCKER.md](DOCKER.md) | Docker инструкция |
| [TEST_REPORT.md](TEST_REPORT.md) | Отчёт о тестах |
| [FINAL_CHECK.md](FINAL_CHECK.md) | Финальная проверка |
| [openapi.yaml](openapi.yaml) | OpenAPI спецификация |

---

## ✅ Чеклист успешного запуска

```
[ ] Docker установлен (если используется)
[ ] Go установлен (если используется локально)
[ ] Python установлен (для тестов)
[ ] Сервер запущен и отвечает на /health
[ ] Тесты проходят (18/18)
[ ] Логи пишутся
[ ] База данных создана
```

---

## 🎉 Готово!

Теперь у вас есть:
- ✅ Работающий API сервер
- ✅ 18 автоматических тестов
- ✅ Docker контейнер
- ✅ Полная документация

**Следующие шаги:**
1. Изучите API в [README.md](README.md)
2. Настройте переменные окружения
3. Разверните на сервере (см. [DOCKER.md](DOCKER.md))
