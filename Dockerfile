# =============================================================================
# Dockerfile для Users API на Go
# =============================================================================
# Многоэтапная сборка для минимального размера образа
# =============================================================================

# -----------------------------------------------------------------------------
# Этап 1: Сборка приложения
# -----------------------------------------------------------------------------
FROM golang:1.21-alpine AS builder

# Установка рабочих переменных
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=1
ENV GO111MODULE=on

# Установка зависимостей для компиляции с CGO (SQLite)
RUN apk add --no-cache gcc musl-dev

# Создание директории приложения
WORKDIR /build

# Копирование файлов модуля
COPY go.mod go.sum ./

# Загрузка зависимостей (кэшируется)
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка бинарного файла с оптимизациями
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -trimpath \
    -o /build/server \
    cmd/server/main.go

# -----------------------------------------------------------------------------
# Этап 2: Финальный минимальный образ
# -----------------------------------------------------------------------------
FROM alpine:3.19

# Установка часового пояса и SSL сертификатов
RUN apk add --no-cache ca-certificates tzdata

# Создание пользователя без прав root для безопасности
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -h /home/appuser -D appuser

# Создание директории для приложения
WORKDIR /app

# Копирование бинарного файла из builder
COPY --from=builder /build/server /app/server

# Копирование конфигурации (опционально)
COPY --from=builder /build/.env.example /app/.env.example

# Создание директории для базы данных
RUN mkdir -p /app/data && chown -R appuser:appgroup /app

# Переключение на пользователя без root прав
USER appuser

# Проброс порта
EXPOSE 8080

# Переменные окружения по умолчанию
ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=8080
ENV DB_PATH=/app/data/test.db
ENV LOG_LEVEL=info
ENV LOG_FILE=

# Проверка здоровья
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Запуск приложения
ENTRYPOINT ["/app/server"]
