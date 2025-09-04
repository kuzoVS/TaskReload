# Многоэтапная сборка для Go приложения
FROM golang:1.21-alpine AS builder

# Установка необходимых пакетов
RUN apk add --no-cache git

# Установка рабочей директории
WORKDIR /app

# Копирование файлов зависимостей
COPY go.mod go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Финальный образ
FROM alpine:latest

# Установка ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копирование исполняемого файла из builder
COPY --from=builder /app/main .

# Создание директории для базы данных
RUN mkdir -p /root/data

# Открытие порта
EXPOSE 8080

# Запуск приложения
CMD ["./main"]
