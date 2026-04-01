# FROM golang:1.22 AS build
# WORKDIR /app
# COPY . .
# RUN go build -o server ./cmd/server

# FROM debian:stable-slim
# WORKDIR /app
# COPY --from=build /app/server .
# CMD ["./server"]

# Stage 1 — Build
FROM golang:1.25 AS build
WORKDIR /app

# Скопировать go.mod и go.sum для кешируемой установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Скопировать весь код
COPY . .

# Собрать бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o service ./cmd/*

# Stage 2 — Runtime
FROM alpine:3.19
WORKDIR /app

# Добавить сертификаты (нужны для https запросов)
RUN apk add --no-cache ca-certificates

COPY --from=build /app/service .
COPY config.yaml .

CMD ["./service", "-config", "config.yaml"]