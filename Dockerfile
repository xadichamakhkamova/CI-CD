  FROM golang:1.22.3-alpine AS builder

  RUN mkdir app

  COPY . /app

  WORKDIR /app

  RUN go build -o main main.go

  FROM alpine:3.18

  WORKDIR /app

  COPY --from=builder /app .

  CMD ["/app/main"]