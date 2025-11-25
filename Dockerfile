# ЭТАП 1: сборка в образе golang (работает как компилятор)
FROM golang:1.25.3 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы проекта (включая go.mod, go.sum)
COPY . .

# Загружаем зависимости
RUN go mod download

# Компилируем приложение ДЛЯ LINUX (даже если сборка идёт на Windows)
# Ключевые переменные:
# GOOS=linux   - целевая ОС
# GOARCH=amd64 - целевая архитектура (для x86_64)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o go-task-manager

# ЭТАП 2: финальный образ (без инструментов сборки)
FROM alpine:latest

# Устанавливаем необходимые зависимости (если нужны)
RUN apk update && apk add --no-cache ca-certificates

# Определяем переменные окружения
ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db

# Создаём рабочую директорию
WORKDIR /app

# Копируем бинарник из этапа сборки
COPY --from=builder /app/go-task-manager /app/go-task-manager

# Копируем веб‑директорию
COPY web /app/web

# Открываем порт
EXPOSE 7540

# Команда запуска с явным указанием пути и аргументов (если нужны)
CMD ["/app/go-task-manager"]