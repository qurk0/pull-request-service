FROM golang:1.25.4-alpine AS builder

WORKDIR /app

COPY app/go.mod app/go.sum ./
RUN go mod download

COPY app/. .

RUN CGO_ENABLED=0 GOOS=linux go build -o pr-service ./cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/pr-service /app/pr-service

COPY app/configs ./configs

EXPOSE 8080

CMD ["./pr-service"]