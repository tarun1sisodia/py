# Use an official Golang runtime as a parent image
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main cmd/server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
ENV PORT 8080
CMD ["/app/main"]
