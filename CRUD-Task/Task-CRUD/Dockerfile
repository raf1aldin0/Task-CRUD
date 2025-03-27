# Stage 1: Build Golang Binary
FROM golang:1.23 AS builder  
# ✅ Versi Golang valid

# Set working directory
WORKDIR /app

# Copy semua file proyek
COPY . .

# Download dependencies
RUN go mod tidy

# Compile Go binary (static build)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd

# Stage 2: Runtime Container
FROM debian:bullseye-slim  
# ✅ Hanya satu stage runtime

WORKDIR /root/

# Copy binary dari tahap builder
COPY --from=builder /app/main .

# Expose port yang digunakan aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]
