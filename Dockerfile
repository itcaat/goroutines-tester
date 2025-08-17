# Многоэтапная сборка для минимизации размера финального образа
FROM golang:1.25.0-alpine AS builder

# Build arguments для версионирования
ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown

# Устанавливаем необходимые пакеты для сборки
RUN apk add --no-cache git ca-certificates

# Создаем рабочую директорию
WORKDIR /app

# Копируем go.mod для кэширования зависимостей
COPY go.mod ./
COPY go.sum ./

# Загружаем зависимости (кэшируется если go.mod не изменился)
RUN go mod download && go mod verify

# Копируем только необходимые файлы (лучшее кэширование)
COPY *.go ./
COPY internal/ ./internal/

# Собираем приложение с версионированием и оптимизацией
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE} -X main.builtBy=docker" \
    -trimpath \
    -o goroutines-tester .

# Финальный образ
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS и curl для healthcheck
RUN apk --no-cache add ca-certificates curl

# Создаем пользователя для безопасности
RUN adduser -D -s /bin/sh appuser

WORKDIR /app

# Копируем бинарный файл из builder стадии
COPY --from=builder /app/goroutines-tester .

# Меняем владельца на appuser
RUN chown appuser:appuser /app/goroutines-tester

# Переключаемся на непривилегированного пользователя
USER appuser

EXPOSE 8888

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8888/metrics || exit 1


CMD ["./goroutines-tester"]
