# Gunakan image Golang terbaru
FROM golang:1.21 AS builder

# Set working directory di dalam container
WORKDIR /app

# Copy semua file proyek ke dalam container
COPY . .

# Download dependencies
RUN go mod tidy

# Build aplikasi
RUN go build -o main .

# Stage kedua: container runtime
FROM debian:bullseye-slim

WORKDIR /root/

# Copy binary dari stage builder
COPY --from=builder /app/main .

# Expose port yang digunakan aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]
