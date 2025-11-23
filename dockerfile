FROM golang:1.23.4-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git postgresql-client

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/PR-service
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:latest

RUN apk add --no-cache ca-certificates postgresql-client

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/.env .env

EXPOSE 8080

CMD ["/root/main"]