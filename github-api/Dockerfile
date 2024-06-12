# Dockerfile

# Build stage
FROM golang:1.19 as builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]
