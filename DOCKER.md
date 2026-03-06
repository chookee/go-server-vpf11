# 🐳 Docker Инструкция по запуску

## 📋 Оглавление

1. [Быстрый старт](#быстрый-старт)
2. [Запуск локально](#запуск-локально)
3. [Запуск на сервере](#запуск-на-сервере)
4. [Тестирование API](#тестирование-api)
5. [Production развёртывание](#production-развёртывание)
6. [Troubleshooting](#troubleshooting)

---

## 🚀 Быстрый старт

### Запуск сервера

```bash
# Сборка и запуск
docker-compose up -d

# Проверка статуса
docker-compose ps

# Просмотр логов
docker-compose logs -f api
```

### Тестирование

```bash
# Запуск тестов через Python скрипт
python test_endpoints.py

# Или через Docker
docker-compose --profile test run test
```

### Остановка

```bash
docker-compose down
```

---

## 💻 Запуск локально

### Вариант 1: Docker Compose (Рекомендуется)

```bash
# 1. Перейдите в директорию проекта
cd E:\Personal\ZeroCode\VPf11\project\go-server

# 2. Запустите сервер
docker-compose up -d --build

# 3. Проверьте что сервер работает
curl http://localhost:8080/health

# 4. Запустите тесты
python test_endpoints.py
```

### Вариант 2: Чистый Docker

```bash
# 1. Сборка образа
docker build -t users-api:latest .

# 2. Запуск контейнера
docker run -d \
  --name users-api \
  -p 8080:8080 \
  -e SERVER_HOST=0.0.0.0 \
  -e SERVER_PORT=8080 \
  -e LOG_LEVEL=debug \
  -v users-api-data:/app/data \
  users-api:latest

# 3. Проверка
docker ps
docker logs -f users-api

# 4. Тест
curl http://localhost:8080/health
```

### Вариант 3: Без Docker (локальная сборка)

```bash
# 1. Сборка Go приложения
go build -o server.exe cmd/server/main.go

# 2. Запуск
./server.exe

# 3. В другом терминале - тесты
python test_endpoints.py
```

---

## 🖥️ Запуск на сервере

### Подготовка

```bash
# 1. Установите Docker и Docker Compose
# Ubuntu/Debian
sudo apt update
sudo apt install docker.io docker-compose

# 2. Добавьте пользователя в группу docker
sudo usermod -aG docker $USER

# 3. Проверьте установку
docker --version
docker-compose --version
```

### Развёртывание

```bash
# 1. Скопируйте файлы на сервер
scp -r go-server/* user@server:/opt/users-api/

# 2. Перейдите в директорию
ssh user@server
cd /opt/users-api

# 3. Настройте переменные окружения
cp .env.example .env
nano .env  # Отредактируйте значения

# 4. Запустите
docker-compose up -d

# 5. Проверьте
docker-compose ps
docker-compose logs -f
```

### Production конфигурация (.env)

```bash
# .env для production
API_PORT=8080
LOG_LEVEL=info
LOG_FILE=/app/logs/app.log

# Таймауты
READ_TIMEOUT=30s
WRITE_TIMEOUT=30s
IDLE_TIMEOUT=120s
SHUTDOWN_TIMEOUT=60s

# Лимиты
MAX_REQUEST_BODY=2097152
DB_MAX_OPEN_CONNS=50
DB_MAX_IDLE_CONNS=10
```

---

## 🧪 Тестирование API

### Через Python скрипт

**Предварительные требования:**
```bash
# Установка Python 3.8+
python --version

# Установка зависимостей
pip install -r requirements-test.txt
```

**Запуск тестов:**
```bash
# Базовый запуск
python test_endpoints.py

# С переменными окружения
API_URL=http://localhost:8080 python test_endpoints.py

# С полным путём
python e:\Personal\ZeroCode\VPf11\project\go-server\test_endpoints.py
```

**Пример вывода:**
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
       ✓ Ожидаемый статус 200

...

══════════════════════════════════════════════════════════════════
                           СВОДКА                                 
══════════════════════════════════════════════════════════════════

Всего тестов:    18
Пройдено:        18
Провалено:       0
Пропущено:       0
Среднее время:   15.23ms

Успешность:      100.0%
[████████████████████████████████████████]

✅ Все тесты пройдены!
```

### Через Docker Compose

```bash
# Запуск тестов в контейнере
docker-compose --profile test run test

# С выводом логов
docker-compose --profile test run --rm test
```

### Ручное тестирование (curl)

```bash
# Health check
curl http://localhost:8080/health

# Создать пользователя
curl -X POST http://localhost:8080/add-user \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice"}'

# Получить всех
curl http://localhost:8080/users

# Получить по ID
curl http://localhost:8080/user/1

# Обновить
curl -X PUT http://localhost:8080/user/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Updated"}'

# Удалить
curl -X DELETE http://localhost:8080/user/1

# Статистика
curl http://localhost:8080/stats/active
```

---

## 🏗️ Production развёртывание

### Docker Swarm

```bash
# Инициализация swarm
docker swarm init

# Деплой стека
docker stack deploy -c docker-compose.prod.yml users-api

# Проверка
docker stack ps users-api
```

### Kubernetes

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: users-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: users-api
  template:
    metadata:
      labels:
        app: users-api
    spec:
      containers:
      - name: api
        image: users-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: SERVER_PORT
          value: "8080"
        - name: DB_PATH
          value: "/app/data/test.db"
        volumeMounts:
        - name: data
          mountPath: /app/data
        resources:
          limits:
            memory: "512Mi"
            cpu: "1000m"
          requests:
            memory: "128Mi"
            cpu: "250m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: users-api-pvc
```

### CI/CD (GitHub Actions)

```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Build Docker image
        run: docker build -t users-api:latest .
      
      - name: Run tests
        run: |
          docker-compose up -d
          sleep 5
          docker-compose --profile test run test
      
      - name: Deploy to server
        run: |
          scp docker-compose.yml user@server:/opt/users-api/
          ssh user@server 'cd /opt/users-api && docker-compose up -d'
```

---

## 🔧 Troubleshooting

### Сервер не запускается

```bash
# Проверка логов
docker-compose logs api

# Проверка портов
docker-compose ps
netstat -tlnp | grep 8080

# Пересборка
docker-compose down
docker-compose up -d --build
```

### Тесты не проходят

```bash
# Проверка доступности API
curl http://localhost:8080/health

# Проверка Python зависимостей
pip install -r requirements-test.txt

# Запуск с отладкой
python -u test_endpoints.py
```

### Проблемы с базой данных

```bash
# Проверка тома
docker volume ls | grep users-api-data

# Очистка данных (ВНИМАНИЕ: удалит все данные!)
docker volume rm users-api-data

# Проверка внутри контейнера
docker exec -it users-api ls -la /app/data
```

### Проблемы с правами доступа

```bash
# Исправление прав (Linux)
sudo chown -R $USER:$USER /path/to/data

# Запуск от root (не рекомендуется)
docker-compose up -d --user root
```

### Мониторинг

```bash
# Статистика использования ресурсов
docker stats users-api

# Логи в реальном времени
docker-compose logs -f api

# Проверка здоровья
docker inspect --format='{{.State.Health.Status}}' users-api
```

---

## 📊 Команды для управления

| Команда | Описание |
|---------|----------|
| `docker-compose up -d` | Запуск в фоне |
| `docker-compose down` | Остановка и удаление |
| `docker-compose ps` | Статус контейнеров |
| `docker-compose logs -f` | Логи в реальном времени |
| `docker-compose restart` | Перезапуск |
| `docker-compose build` | Пересборка образов |
| `docker-compose exec api sh` | Shell в контейнере |
| `docker volume ls` | Список томов |
| `docker volume rm <name>` | Удаление тома |

---

## ✅ Чеклист развёртывания

```
[ ] Docker установлен
[ ] Docker Compose установлен
[ ] Файлы скопированы на сервер
[ ] .env настроен для production
[ ] Порты открыты в firewall
[ ] Health check работает
[ ] Тесты проходят
[ ] Логи настроены
[ ] Мониторинг настроен
[ ] Backup данных настроен
```

---

## 📚 Дополнительные ресурсы

- [Docker документация](https://docs.docker.com/)
- [Docker Compose документация](https://docs.docker.com/compose/)
- [Go Docker best practices](https://github.com/GoogleContainerTools/distroless)
