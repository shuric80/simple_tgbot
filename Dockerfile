# Используем официальный образ Go
FROM golang:1.21-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum файлы
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Используем минимальный образ для продакшена
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем собранное приложение из builder образа
COPY --from=builder /app/main .

# Запускаем приложение
CMD ["./main"]
