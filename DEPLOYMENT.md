# 🐳 Docker Hub — Полное руководство по развёртыванию

## 📋 Оглавление

1. [Что установить на Windows 11](#что-установить-на-windows-11)
2. [Создание Docker Hub аккаунта](#создание-docker-hub-аккаунта)
3. [Сборка и публикация образа](#сборка-и-публикация-образа)
4. [Запуск из Docker Hub](#запуск-из-docker-hub)
5. [Развёртывание на сервере](#развёртывание-на-сервере)
6. [CI/CD с GitHub Actions](#cicd-с-github-actions)
7. [Troubleshooting](#troubleshooting)

---

## 💻 Что установить на Windows 11

### 1. Docker Desktop (Обязательно)

**Ссылка:** https://www.docker.com/products/docker-desktop/

**Требования:**
- Windows 11 64-bit
- WSL 2 (Windows Subsystem for Linux)
- 4GB RAM (рекомендуется 8GB+)
- Виртуализация включена в BIOS

**Установка:**

```powershell
# 1. Скачайте установщик с официального сайта
# https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe

# 2. Запустите установщик
# Docker Desktop Installer.exe

# 3. Или через winget (менеджер пакетов Windows)
winget install Docker.DockerDesktop

# 4. Проверьте установку
docker --version
docker-compose --version
```

**Настройка после установки:**

1. Запустите Docker Desktop
2. Дождитесь запуска (зелёный индикатор внизу)
3. Настройки → Settings → General:
   - ✅ Use WSL 2 instead of Hyper-V
   - ✅ Start Docker Desktop when you log in
4. Settings → Resources:
   - CPUs: 2-4
   - Memory: 2-4 GB
   - Swap: 1 GB

---

### 2. Git (Рекомендуется)

**Ссылка:** https://git-scm.com/download/win

**Установка:**
```powershell
# Через winget
winget install Git.Git

# Или скачайте с https://git-scm.com/download/win
```

---

### 3. Visual Studio Code (Рекомендуется)

**Ссылка:** https://code.visualstudio.com/

**Расширения:**
- Docker (Microsoft)
- Docker Compose
- Remote - WSL
- GitLens

---

### 4. Windows Terminal (Рекомендуется)

**Ссылка:** Microsoft Store или https://github.com/microsoft/terminal

**Установка:**
```powershell
winget install Microsoft.WindowsTerminal
```

---

### 5. Python 3.8+ (Для тестов)

**Ссылка:** https://www.python.org/downloads/

**Установка:**
```powershell
winget install Python.Python.3.11

# Проверка
python --version
pip --version
```

---

## ✅ Проверка установки

```powershell
# Откройте PowerShell или Windows Terminal

# 1. Проверка Docker
docker --version
docker info

# 2. Проверка Docker Compose
docker-compose --version

# 3. Проверка WSL
wsl --list --verbose

# 4. Проверка Git
git --version

# 5. Проверка Python
python --version
```

**Ожидаемый вывод:**
```
Docker version 25.0.0, build ...
Docker Compose version v2.24.0
WSL version: 2.0.14.0
git version 2.43.0.windows.1
Python 3.11.8
```

---

## 🐳 Создание Docker Hub аккаунта

### 1. Регистрация

1. Перейдите на https://hub.docker.com/
2. Нажмите **Sign Up**
3. Введите:
   - Username (уникальный, будет в URL образа)
   - Email
   - Password
4. Подтвердите email

### 2. Создание токена (для CI/CD)

1. Войдите в Docker Hub
2. Кликните на username → **Account Settings**
3. **Security** → **New Access Token**
4. Введите название (например, `github-ci`)
5. Скопируйте токен (покажется один раз!)
6. Сохраните в надёжном месте

---

## 📦 Сборка и публикация образа

### 1. Вход в Docker Hub из терминала

```powershell
# Windows PowerShell
docker login

# Введите:
# Username: ваш username Docker Hub
# Password: ваш пароль или токен
```

**Пример:**
```powershell
PS C:\Users\YourName> docker login
Login with your Docker ID to push and pull images from Docker Hub.
Username: yourusername
Password: ********
Login Succeeded
```

---

### 2. Сборка образа

```powershell
# Перейдите в директорию проекта
cd E:\Personal\ZeroCode\VPf11\project\go-server

# Сборка образа с тегом
docker build -t yourusername/users-api:latest .

# Пример с конкретным username
docker build -t szakharchuk/users-api:latest .
```

**Проверка сборки:**
```powershell
# Список локальных образов
docker images

# Детальная информация
docker inspect szakharchuk/users-api:latest
```

---

### 3. Публикация в Docker Hub

```powershell
# Push образа в Docker Hub
docker push yourusername/users-api:latest

# Пример
docker push szakharchuk/users-api:latest
```

**Пример вывода:**
```
The push refers to repository [docker.io/szakharchuk/users-api]
abc123: Pushed
def456: Pushed
ghi789: Pushed
latest: digest: sha256:... size: 1234
```

---

### 4. Версионирование образов

```powershell
# Создание тега для версии
docker tag szakharchuk/users-api:latest szakharchuk/users-api:1.0.0
docker tag szakharchuk/users-api:latest szakharchuk/users-api:v1.0.0

# Push всех версий
docker push szakharchuk/users-api:1.0.0
docker push szakharchuk/users-api:v1.0.0
docker push szakharchuk/users-api:latest
```

**Проверка на Docker Hub:**
- Перейдите на https://hub.docker.com/r/yourusername/users-api
- Проверьте вкладку **Tags**

---

## 🚀 Запуск из Docker Hub

### Локальный запуск

```powershell
# Простой запуск
docker run -d -p 8080:8080 --name users-api yourusername/users-api:latest

# С переменными окружения
docker run -d \
  -p 8080:8080 \
  -e SERVER_PORT=8080 \
  -e LOG_LEVEL=info \
  -e DB_PATH=/app/data/test.db \
  --name users-api \
  yourusername/users-api:latest

# С томом для сохранения данных
docker run -d \
  -p 8080:8080 \
  -v users-api-data:/app/data \
  --name users-api \
  yourusername/users-api:latest
```

### Проверка работы

```powershell
# Статус контейнера
docker ps

# Логи
docker logs users-api

# Проверка API
curl http://localhost:8080/health

# Остановка
docker stop users-api
docker rm users-api
```

---

## 🖥️ Развёртывание на сервере

### Вариант 1: Ручное развёртывание

**На сервере (Linux):**

```bash
# 1. Установка Docker
sudo apt update
sudo apt install -y docker.io docker-compose

# 2. Добавление пользователя в группу docker
sudo usermod -aG docker $USER
newgrp docker

# 3. Вход в Docker Hub
docker login
# Username: yourusername
# Password: ********

# 4. Запуск образа из Docker Hub
docker run -d \
  -p 8080:8080 \
  -v /opt/users-api/data:/app/data \
  -e SERVER_PORT=8080 \
  -e LOG_LEVEL=info \
  --name users-api \
  --restart unless-stopped \
  yourusername/users-api:latest

# 5. Проверка
docker ps
curl http://localhost:8080/health
```

---

### Вариант 2: Через docker-compose на сервере

**Создайте `docker-compose.yml` на сервере:**

```yaml
version: '3.8'

services:
  api:
    image: yourusername/users-api:latest
    container_name: users-api
    restart: unless-stopped
    
    ports:
      - "8080:8080"
    
    environment:
      - SERVER_HOST=0.0.0.0
      - SERVER_PORT=8080
      - DB_PATH=/app/data/test.db
      - LOG_LEVEL=info
    
    volumes:
      - /opt/users-api/data:/app/data
      - /opt/users-api/logs:/app/logs
    
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3
```

**Запуск:**

```bash
# На сервере
cd /opt/users-api
docker-compose up -d

# Проверка
docker-compose ps
docker-compose logs -f
```

---

### Вариант 3: Автоматическое обновление

**Скрипт для обновления (`update.sh`):**

```bash
#!/bin/bash

echo "🔄 Обновление users-api..."

# Остановка старого контейнера
docker stop users-api
docker rm users-api

# Pull новой версии
docker pull yourusername/users-api:latest

# Запуск нового контейнера
docker run -d \
  -p 8080:8080 \
  -v /opt/users-api/data:/app/data \
  -e SERVER_PORT=8080 \
  -e LOG_LEVEL=info \
  --name users-api \
  --restart unless-stopped \
  yourusername/users-api:latest

echo "✅ Обновление завершено!"

# Проверка
sleep 5
curl http://localhost:8080/health
```

**Использование:**
```bash
chmod +x update.sh
./update.sh
```

---

## 🔄 CI/CD с GitHub Actions

### Автоматическая публикация при push

**Создайте `.github/workflows/docker-publish.yml`:**

```yaml
name: Docker Publish

on:
  push:
    branches: [main]
    tags: ['v*']
  pull_request:
    branches: [main]

env:
  DOCKER_IMAGE: szakharchuk/users-api
  DOCKERHUB_USERNAME: ${{ secrets.DOCKERHUB_USERNAME }}
  DOCKERHUB_TOKEN: ${{ secrets.DOCKERHUB_TOKEN }}

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      
      - name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ env.DOCKERHUB_USERNAME }}
          password: ${{ env.DOCKERHUB_TOKEN }}
      
      - name: Build and push (main branch)
        if: github.ref == 'refs/heads/main'
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: |
            ${{ env.DOCKER_IMAGE }}:latest
            ${{ env.DOCKER_IMAGE }}:dev
      
      - name: Build and push (tag)
        if: startsWith(github.ref, 'refs/tags/v')
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: |
            ${{ env.DOCKER_IMAGE }}:${{ github.ref_name }}
            ${{ env.DOCKER_IMAGE }}:latest
      
      - name: Run tests
        run: |
          docker-compose up -d
          sleep 5
          pip install -r requirements-test.txt
          python test_endpoints.py
          docker-compose down
```

### Настройка секретов в GitHub

1. Перейдите в репозиторий на GitHub
2. **Settings** → **Secrets and variables** → **Actions**
3. **New repository secret**:
   - `DOCKERHUB_USERNAME`: ваш username Docker Hub
   - `DOCKERHUB_TOKEN`: токен Docker Hub

---

## 📊 Полный рабочий процесс

### Локальная разработка

```powershell
# 1. Изменение кода
# Редактируете файлы в VS Code

# 2. Тестирование
docker-compose up -d
python test_endpoints.py
docker-compose down

# 3. Сборка образа
docker build -t szakharchuk/users-api:dev .

# 4. Локальная проверка
docker run -d -p 8080:8080 szakharchuk/users-api:dev
curl http://localhost:8080/health
```

### Публикация

```powershell
# 1. Коммит изменений
git add .
git commit -m "feat: new feature"
git push

# 2. Сборка и push (если нет CI/CD)
docker build -t szakharchuk/users-api:latest .
docker push szakharchuk/users-api:latest
```

### Развёртывание на сервере

```bash
# На сервере
ssh user@server

# Обновление
cd /opt/users-api
docker-compose pull
docker-compose up -d

# Проверка
docker-compose ps
curl http://localhost:8080/health
```

---

## 🔧 Troubleshooting

### Docker Desktop не запускается

```powershell
# Проверка WSL
wsl --list --verbose

# Обновление WSL
wsl --update

# Переустановка WSL
wsl --unregister docker-desktop
wsl --unregister docker-desktop-data

# Перезапуск Docker Desktop
```

### Ошибка авторизации Docker Hub

```powershell
# Выход и повторный вход
docker logout
docker login

# Проверка токена
# Создайте новый токен в Docker Hub
```

### Образ не пушится

```powershell
# Проверка имени образа
docker images

# Должно быть: yourusername/users-api:latest
# Если просто users-api:latest — перетегируйте
docker tag users-api:latest yourusername/users-api:latest
docker push yourusername/users-api:latest
```

### Контейнер не запускается на сервере

```bash
# Проверка логов
docker logs users-api

# Проверка портов
netstat -tlnp | grep 8080

# Проверка прав доступа
ls -la /opt/users-api/data

# Исправление прав
sudo chown -R 1000:1000 /opt/users-api/data
```

---

## 📚 Шпаргалка команд

| Команда | Описание |
|---------|----------|
| `docker login` | Вход в Docker Hub |
| `docker build -t user/app:tag .` | Сборка образа |
| `docker push user/app:tag` | Публикация в Docker Hub |
| `docker pull user/app:tag` | Загрузка из Docker Hub |
| `docker run -d -p 8080:8080 user/app` | Запуск контейнера |
| `docker ps` | Список контейнеров |
| `docker logs container` | Просмотр логов |
| `docker stop/rm container` | Остановка/удаление |
| `docker images` | Список образов |
| `docker tag old new` | Переименование тега |

---

## ✅ Чеклист развёртывания

```
[ ] Docker Desktop установлен на Windows 11
[ ] WSL 2 настроен
[ ] Docker Hub аккаунт создан
[ ] Токен Docker Hub сохранён
[ ] Образ собран локально
[ ] Образ опубликован в Docker Hub
[ ] Тесты проходят
[ ] Сервер настроен
[ ] Docker установлен на сервере
[ ] Образ запущен на сервере
[ ] Health check работает
[ ] CI/CD настроен (опционально)
```

---

## 🎉 Итог

**Теперь вы можете:**

1. ✅ Собирать Docker образы на Windows 11
2. ✅ Публиковать в Docker Hub
3. ✅ Запускать локально для тестов
4. ✅ Развёртывать на любом сервере
5. ✅ Автоматизировать через CI/CD

**Ваш образ на Docker Hub:**
```
https://hub.docker.com/r/yourusername/users-api
```

**Команда для развёртывания где угодно:**
```bash
docker run -d -p 8080:8080 yourusername/users-api:latest
```
