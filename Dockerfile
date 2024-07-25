# Используйте официальный образ Go как базовый
FROM golang:1.22 as builder

# Установите рабочий каталог
WORKDIR /app

# Копируйте модули Go
COPY go.mod go.sum ./
RUN go mod download

# Копируйте исходный код проекта
COPY . .

# Соберите приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /api

# Используйте образ alpine для финального образа
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируйте скомпилированный файл из предыдущего шага
COPY --from=builder /api .

# Откройте порт, который использует ваше API
EXPOSE 8080

# Запустите ваше API
CMD ["./api"]
