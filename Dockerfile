FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /quotation ./cmd/app


FROM golang:1.24-alpine AS tester
WORKDIR /app
COPY --from=builder /app .
RUN go test -v ./internal/tests


FROM alpine:latest
WORKDIR /app
COPY --from=builder /quotation .
COPY .env .

EXPOSE 8080
CMD ["./quotation"]