# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники и собираем
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api

# Run stage
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Europe/Moscow

WORKDIR /app

COPY --from=builder /app/api .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./api"]
