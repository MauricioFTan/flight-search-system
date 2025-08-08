# Stage 1: Build application
# PASTIKAN VERSI GO DI SINI SUDAH 1.22 atau lebih baru
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy file go.mod dan go.sum untuk men-download dependensi
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build aplikasi Go.
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main/main.go

# Stage 2: Run application
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

# Perintah untuk menjalankan aplikasi
CMD ["./main"]