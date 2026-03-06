# 🐳 Docker Hub: Публикация и развёртывание Users API

> Полное пошаговое руководство по публикации образа в Docker Hub и развёртыванию на production сервере.

**Версия:** 1.0.0  
**Проект:** Users API (Go 1.21)  
**Статус:** ✅ Production Ready

---

## 📋 Оглавление

1. [Подготовка на Windows 11](#1-подготовка-на-windows-11)
2. [Создание Docker Hub аккаунта](#2-создание-docker-hub-аккаунта)
3. [Сборка и публикация образа](#3-сборка-и-публикация-образа)
4. [Развёртывание на сервере (Linux)](#4-развёртывание-на-сервере-linux)
5. [Настройка CI/CD (GitHub Actions)](#5-настройка-cicd-github-actions)
6. [Обновление и мониторинг](#6-обновление-и-мониторинг)
7. [Troubleshooting](#7-troubleshooting)

---

## 1. Подготовка на Windows 11

### 1.1 Установка Docker Desktop

**Скачивание и установка:**

```powershell
# Способ 1: Через winget (рекомендуется)
winget install Docker.DockerDesktop

# Способ 2: Вручную
# Перейдите на https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe
# Скачайте и запустите установщик
```

**Требования:**
- Windows 11 64-bit
- WSL 2 (Windows Subsystem for Linux)
- 4GB RAM (рекомендуется 8GB+)
- Виртуализация включена в BIOS

**Настройка WSL 2:**

```powershell
# Установка WSL 2
wsl --install

# Проверка
wsl --list --verbose

# Обновление
wsl --update
```

**Перезагрузите компьютер** после установки.

### 1.2 Проверка установки

```powershell
# Откройте PowerShell или Windows Terminal

docker --version
docker-compose --version
docker info
```

**Ожидаемый вывод:**
```
Docker version 25.0.0, build ...
Docker Compose version v2.24.0
```

### 1.3 Запуск Docker Desktop

1. Запустите Docker Desktop из меню Пуск
2. Дождитесь зелёного индикатора "Engine running"
3. Принимите условия использования (если появится)

---

## 2. Создание Docker Hub аккаунта

### 2.1 Регистрация

1. Перейдите на https://hub.docker.com/
2. Нажмите **Sign Up**
3. Заполните форму:
   - **Username** — уникальный идентификатор (будет в URL образа)
   - **Email** — действующий email
   - **Password** — надёжный пароль
4. Подтвердите email из письма

> **Важно:** Username должен быть уникальным. Придумайте его заранее (например, `szakharchuk`, `mycompany`, `projectname`).

### 2.2 Создание токена для CI/CD

1. Войдите в Docker Hub
2. Кликните на username (справа вверху) → **Account Settings**
3. Перейдите во вкладку **Security**
4. Нажмите **New Access Token**
5. Заполните:
   - **Token Name:** `github-ci` или `personal`
   - **Expiration:** выберите срок (рекомендуется 90 дней)
   - **Permissions:** Read & Write
6. Нажмите **Generate**
7. **Скопируйте токен** (покажется один раз!)
8. Сохраните в надёжном месте (менеджер паролей)

---

## 3. Сборка и публикация образа

### 3.1 Вход в Docker Hub из терминала

Откройте PowerShell или Windows Terminal:

```powershell
docker login
```

Введите:
- **Username:** ваш username Docker Hub
- **Password:** пароль или токен из раздела 2.2

**Ожидаемый результат:**
```
Login Succeeded
```

### 3.2 Переход в директорию проекта

```powershell
cd e:\Personal\ZeroCode\VPf11\project\go-server
```

Проверка наличия файлов:
```powershell
ls Dockerfile, docker-compose.yml, go.mod
```

### 3.3 Сборка Docker образа

```powershell
docker build -t YOUR_USERNAME/users-api:latest .
```

> **Замените `YOUR_USERNAME`** на ваш username Docker Hub (например, `szakharchuk`).

**Пример:**
```powershell
docker build -t szakharchuk/users-api:latest .
```

**Что происходит:**
- `-t szakharchuk/users-api:latest` — тег образа
- `.` — текущая директория с Dockerfile
- Multi-stage сборка создаст минимальный образ (~50MB)

**Проверка сборки:**
```powershell
docker images | findstr users-api
```

**Ожидаемый результат:**
```
szakharchuk/users-api   latest   abc123def456   2 minutes ago   50MB
```

### 3.4 Тестирование образа локально

Перед публикацией проверьте работу образа:

```powershell
# Запуск контейнера
docker run -d -p 8080:8080 --name test-users-api szakharchuk/users-api:latest

# Проверка статуса
docker ps

# Проверка логов
docker logs test-users-api

# Проверка API
curl http://localhost:8080/health
```

**Ожидаемый ответ:**
```json
{"status":"ok","timestamp":"2024-03-06T..."}
```

**Остановка тестового контейнера:**
```powershell
docker stop test-users-api
docker rm test-users-api
```

### 3.5 Публикация в Docker Hub

```powershell
docker push szakharchuk/users-api:latest
```

**Ожидаемый результат:**
```
The push refers to repository [docker.io/szakharchuk/users-api]
abc123: Pushed
def456: Pushed
ghi789: Pushed
latest: digest: sha256:xyz789 size: 1234
```

**Проверка на Docker Hub:**
- Перейдите на https://hub.docker.com/r/szakharchuk/users-api
- Вы увидите ваш образ с тегом `latest`

### 3.6 Версионирование образов (рекомендуется)

Для production используйте версии:

```powershell
# Создание тега версии
docker tag szakharchuk/users-api:latest szakharchuk/users-api:1.0.0

# Push версии
docker push szakharchuk/users-api:1.0.0

# Push всех тегов
docker push szakharchuk/users-api --all-tags
```

**Пример полного цикла:**
```powershell
# Текущая версия
docker tag szakharchuk/users-api:latest szakharchuk/users-api:v1.0.0
docker push szakharchuk/users-api:v1.0.0

# Следующая версия
docker tag szakharchuk/users-api:latest szakharchuk/users-api:v1.0.1
docker push szakharchuk/users-api:v1.0.1

# Обновление latest
docker push szakharchuk/users-api:latest
```

---

## 4. Развёртывание на сервере (Linux)

### 4.1 Подготовка сервера

**Требования:**
- Linux сервер (Ubuntu 20.04/22.04, Debian 11+, CentOS 8+)
- root доступ или sudo
- Открытый порт 8080 (или другой)
- Минимум 512MB RAM, 1 CPU

**Подключение к серверу:**
```bash
ssh user@your-server-ip
```

### 4.2 Установка Docker на сервер

#### Для Ubuntu/Debian:

```bash
# Обновление пакетов
sudo apt update

# Установка Docker и Docker Compose
sudo apt install -y docker.io docker-compose

# Проверка версии
docker --version
docker-compose --version

# Добавление пользователя в группу docker (чтобы не использовать sudo)
sudo usermod -aG docker $USER
newgrp docker

# Проверка без sudo
docker ps
```

#### Для CentOS/RHEL:

```bash
# Установка Docker
sudo yum install -y docker
sudo systemctl start docker
sudo systemctl enable docker

# Установка Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Добавление пользователя в группу docker
sudo usermod -aG docker $USER
newgrp docker
```

### 4.3 Вход в Docker Hub на сервере

```bash
docker login
```

Введите ваш username и пароль/токен.

**Ожидаемый результат:**
```
Login Succeeded
```

### 4.4 Запуск из Docker Hub (простой способ)

```bash
docker run -d \
  -p 8080:8080 \
  -v /opt/users-api/data:/app/data \
  -e SERVER_PORT=8080 \
  -e LOG_LEVEL=info \
  --name users-api \
  --restart unless-stopped \
  szakharchuk/users-api:latest
```

**Параметры:**
| Параметр | Описание |
|----------|----------|
| `-d` | Запуск в фоне (daemon mode) |
| `-p 8080:8080` | Проброс порта (хост:контейнер) |
| `-v /opt/users-api/data:/app/data` | Том для сохранения БД |
| `-e SERVER_PORT=8080` | Переменная окружения |
| `--name users-api` | Имя контейнера |
| `--restart unless-stopped` | Автоперезапуск при падении |
| `szakharchuk/users-api:latest` | Образ из Docker Hub |

**Проверка:**
```bash
# Статус контейнера
docker ps

# Логи
docker logs users-api

# Проверка API (на сервере)
curl http://localhost:8080/health

# Проверка API (с вашего компьютера)
curl http://your-server-ip:8080/health
```

### 4.5 Развёртывание через docker-compose (рекомендуется)

**Создайте директорию проекта:**

```bash
mkdir -p /opt/users-api
cd /opt/users-api
```

**Создайте файл `.env`:**

```bash
nano .env
```

**Содержимое `.env`:**
```env
API_PORT=8080
LOG_LEVEL=info
```

**Создайте `docker-compose.yml`:**

```bash
nano docker-compose.yml
```

**Содержимое `docker-compose.yml`:**
```yaml
version: '3.8'

services:
  api:
    image: szakharchuk/users-api:latest
    container_name: users-api
    restart: unless-stopped

    ports:
      - "${API_PORT:-8080}:8080"

    environment:
      - SERVER_HOST=0.0.0.0
      - SERVER_PORT=8080
      - DB_PATH=/app/data/test.db
      - LOG_LEVEL=${LOG_LEVEL:-info}

    volumes:
      - /opt/users-api/data:/app/data

    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 5s
```

**Запуск:**

```bash
docker-compose up -d
```

**Проверка:**
```bash
docker-compose ps
docker-compose logs -f
curl http://localhost:8080/health
```

### 4.6 Настройка firewall

#### Для Ubuntu (UFW):

```bash
sudo ufw allow 8080/tcp
sudo ufw reload
sudo ufw status
```

#### Для CentOS (firewalld):

```bash
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --reload
sudo firewall-cmd --list-ports
```

#### Для AWS EC2:

1. EC2 Console → Security Groups → Ваш Security Group
2. Вкладка **Inbound rules** → **Edit inbound rules**
3. **Add rule**:
   - Type: Custom TCP
   - Port: 8080
   - Source: 0.0.0.0/0
4. **Save rules**

#### Для Azure:

1. Network Security Group → Inbound Rules → Add
2. Port: 8080, Action: Allow, Priority: 100

#### Для Google Cloud:

1. VPC Network → Firewall → Create Firewall Rule
2. Name: `allow-users-api`
3. Port: 8080, Source IP: 0.0.0.0/0

---

## 5. Настройка CI/CD (GitHub Actions)

### 5.1 Создание workflow файла

В вашем репозитории GitHub создайте файл:

`.github/workflows/docker-publish.yml`

**Содержимое:**
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
            ${{ env.DOCKER_IMAGE }}:dev-${{ github.sha }}

      - name: Build and push (tag)
        if: startsWith(github.ref, 'refs/tags/v')
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: |
            ${{ env.DOCKER_IMAGE }}:${{ github.ref_name }}
            ${{ env.DOCKER_IMAGE }}:latest
```

### 5.2 Настройка секретов в GitHub

1. Перейдите в репозиторий на GitHub
2. **Settings** → **Secrets and variables** → **Actions**
3. **New repository secret**:

| Name | Value |
|------|-------|
| `DOCKERHUB_USERNAME` | Ваш username Docker Hub (например, `szakharchuk`) |
| `DOCKERHUB_TOKEN` | Токен из раздела 2.2 |

4. Нажмите **Add secret** для каждого

### 5.3 Тестирование CI/CD

```bash
# Локально (в директории проекта)
git add .github/workflows/docker-publish.yml
git commit -m "ci: add Docker Hub auto-publish workflow"
git push
```

**Проверка:**
1. Перейдите на GitHub → ваш репозиторий → вкладка **Actions**
2. Вы увидите запущенный workflow "Docker Publish"
3. После завершения (зелёная галочка) проверьте Docker Hub
4. Должен появиться новый образ с тегом `dev-<commit-hash>`

### 5.4 Автоматическая публикация по тегу

```bash
# Создание тега версии
git tag v1.0.0
git push origin v1.0.0
```

**Результат:**
- GitHub Actions автоматически соберёт образ
- Опубликует с тегами `v1.0.0` и `latest`
- Появится в Docker Hub

---

## 6. Обновление и мониторинг

### 6.1 Ручное обновление

```bash
# На сервере
cd /opt/users-api

# Pull новой версии
docker-compose pull

# Перезапуск
docker-compose up -d

# Проверка
docker-compose ps
curl http://localhost:8080/health
```

### 6.2 Скрипт автоматического обновления

**Создайте скрипт:**

```bash
nano /opt/users-api/update.sh
```

**Содержимое:**
```bash
#!/bin/bash

echo "🔄 Обновление users-api..."

cd /opt/users-api

# Pull новой версии
docker-compose pull

# Перезапуск с пересозданием контейнера
docker-compose up -d --force-recreate

# Очистка старых образов
docker image prune -f

echo "✅ Обновление завершено!"

# Проверка через 5 секунд
sleep 5
curl -s http://localhost:8080/health
```

**Запуск:**
```bash
chmod +x update.sh
./update.sh
```

### 6.3 Автоматическое обновление через Watchtower

**Добавьте в docker-compose.yml:**

```yaml
services:
  watchtower:
    image: containrrr/watchtower
    container_name: watchtower
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    command: --interval 300 users-api
    restart: unless-stopped
```

**Что делает:**
- Проверяет Docker Hub каждые 5 минут (300 секунд)
- Автоматически обновляет контейнер `users-api`
- Перезапускает при изменении образа

### 6.4 Просмотр логов

```bash
# Последние строки
docker logs users-api

# В реальном времени
docker logs -f users-api

# С временными метками
docker logs -t users-api

# Последние 100 строк
docker logs --tail 100 users-api

# Через docker-compose
docker-compose logs -f
docker-compose logs --tail 100
```

### 6.5 Проверка здоровья

```bash
# Статус контейнера
docker ps

# Детальная информация
docker inspect users-api

# Проверка healthcheck
docker inspect --format='{{.State.Health.Status}}' users-api
```

### 6.6 Метрики Docker

```bash
# Использование ресурсов (CPU, RAM)
docker stats users-api

# Информация о контейнере
docker inspect users-api
```

---

## 7. Troubleshooting

### Проблема: Docker Desktop не запускается на Windows

```powershell
# Проверка WSL
wsl --list --verbose

# Обновление WSL
wsl --update

# Сброс Docker
# Docker Desktop → Settings → Troubleshoot → Reset to factory defaults

# Переустановка WSL
wsl --unregister docker-desktop
wsl --unregister docker-desktop-data
```

### Проблема: Ошибка авторизации `docker login`

```powershell
# Очистка кэша
docker logout

# Повторный вход
docker login

# Если не работает — создайте новый токен в Docker Hub:
# Account Settings → Security → New Access Token
```

### Проблема: Образ не пушится

```powershell
# Проверка имени образа
docker images

# Должно быть: szakharchuk/users-api:latest
# Если просто users-api:latest — перетегируйте
docker tag users-api:latest szakharchuk/users-api:latest
docker push szakharchuk/users-api:latest
```

### Проблема: Контейнер не запускается на сервере

```bash
# Проверка логов
docker logs users-api

# Проверка портов
sudo netstat -tlnp | grep 8080

# Проверка прав доступа к тому
ls -la /opt/users-api/data
sudo chown -R 1000:1000 /opt/users-api/data

# Проверка firewall
sudo ufw status
sudo ufw allow 8080/tcp
```

### Проблема: API не доступен извне

```bash
# Проверка на сервере
curl http://localhost:8080/health

# Если работает локально, но не извне:
# 1. Проверьте Security Group (AWS) / Firewall (Azure)
# 2. Проверьте iptables
sudo iptables -L -n | grep 8080

# 3. Проверьте, что SERVER_HOST=0.0.0.0
docker inspect users-api | grep SERVER_HOST
```

### Проблема: GitHub Actions не запускается

1. Проверьте файл workflow на синтаксис:
   ```bash
   yamllint .github/workflows/docker-publish.yml
   ```
2. Проверьте секреты в GitHub Settings → Secrets
3. Проверьте логи workflow на вкладке Actions

---

## ✅ Чеклист полного развёртывания

```
[ ] Docker Desktop установлен на Windows 11
[ ] WSL 2 настроен и работает
[ ] Docker Hub аккаунт создан
[ ] Токен Docker Hub сохранён в менеджере паролей
[ ] Образ собран локально: docker build -t username/users-api:latest .
[ ] Образ протестирован локально: docker run -d -p 8080:8080 ...
[ ] Образ опубликован: docker push username/users-api:latest
[ ] Образ доступен на https://hub.docker.com/r/username/users-api
[ ] Сервер подготовлен (Linux, Docker установлен)
[ ] Firewall настроен (порт 8080 открыт)
[ ] Docker login выполнен на сервере
[ ] Образ запущен на сервере
[ ] Health check работает: curl http://server-ip:8080/health
[ ] CI/CD настроен (GitHub Actions workflow)
[ ] Секреты GitHub добавлены (DOCKERHUB_USERNAME, DOCKERHUB_TOKEN)
[ ] Тестовый push прошёл успешно
```

---

## 📚 Шпаргалка команд

| Задача | Команда |
|--------|---------|
| **Вход в Docker Hub** | `docker login` |
| **Сборка образа** | `docker build -t username/users-api:latest .` |
| **Публикация** | `docker push username/users-api:latest` |
| **Загрузка** | `docker pull username/users-api:latest` |
| **Запуск** | `docker run -d -p 8080:8080 --name users-api username/users-api:latest` |
| **Остановка** | `docker stop users-api && docker rm users-api` |
| **Логи** | `docker logs -f users-api` |
| **docker-compose старт** | `docker-compose up -d` |
| **docker-compose стоп** | `docker-compose down` |
| **Обновление** | `docker-compose pull && docker-compose up -d` |
| **Статус** | `docker ps` или `docker-compose ps` |
| **Метрики** | `docker stats users-api` |

---

## 🎉 Итог

**Теперь вы можете:**

1. ✅ Собирать Docker образы на Windows 11
2. ✅ Публиковать в Docker Hub
3. ✅ Запускать локально для тестов
4. ✅ Развёртывать на любом сервере с Docker
5. ✅ Автоматизировать через CI/CD (GitHub Actions)
6. ✅ Обновлять production без простоя

**Ваш образ на Docker Hub:**
```
https://hub.docker.com/r/YOUR_USERNAME/users-api
```

**Команда для развёртывания где угодно:**
```bash
docker run -d -p 8080:8080 --restart unless-stopped YOUR_USERNAME/users-api:latest
```

---

## 📞 Поддержка

При возникновении проблем:

1. Проверьте раздел [Troubleshooting](#7-troubleshooting)
2. Изучите логи: `docker logs users-api`
3. Проверьте статус: `docker ps` или `docker-compose ps`
4. Убедитесь, что firewall открыт на порту 8080

---

**Документ обновлён:** 6 марта 2026 г.  
**Проект:** Users API (Go 1.21 + chi + SQLite)  
**Статус:** ✅ Production Ready
