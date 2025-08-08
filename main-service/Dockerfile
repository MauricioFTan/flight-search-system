FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy file go.mod dan go.sum untuk men-download dependensi
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build aplikasi Go. Menghasilkan binary bernama 'main'
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main/main.go

# Stage 2: Run application
FROM alpine:latest

WORKDIR /app

# Copy binary yang sudah di-build dari stage 'builder'
COPY --from=builder /app/main .

# Expose port yang digunakan oleh aplikasi
EXPOSE 8080

# Perintah untuk menjalankan aplikasi saat kontainer dimulai
CMD ["./main"]