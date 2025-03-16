# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY webapp/go.mod webapp/go.sum ./
RUN go mod download

COPY webapp/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o webapp .


# Runtime 
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/webapp .

EXPOSE 8080

CMD ["./webapp"]

