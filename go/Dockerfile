# Tahap Build
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
ENV GOTOOLCHAIN=auto
RUN go build -o server main.go

# Tahap Eksekusi
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 50051
CMD ["./server"]