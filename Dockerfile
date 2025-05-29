# FROM golang:1.24-alpine AS builder

# WORKDIR /cmd/app
# COPY . .
# RUN go mod download
# RUN CGO_ENABLED=0 GOOS=linux go build -o /wallet ./cmd/app

# FROM alpine:latest
# WORKDIR /app
# COPY --from=builder /wallet .
# COPY .env .

# EXPOSE 4000
# CMD ["./wallet"]
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /wallet ./cmd/app

# Стадия для тестов (используем тот же образ builder)
FROM golang:1.24-alpine AS tester
WORKDIR /app
COPY --from=builder /app .
RUN go test -v ./internal/tests

# Финальный образ
FROM alpine:latest
WORKDIR /app
COPY --from=builder /wallet .
COPY .env .

EXPOSE 4000
CMD ["./wallet"]