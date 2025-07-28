FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o ./url-shortener ./cmd/api

FROM alpine:latest
COPY --from=builder /app/migration /app/migration
COPY --from=builder /app/url-shortener /app/
WORKDIR /app

EXPOSE 8080

CMD ["./url-shortener"]