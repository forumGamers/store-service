# Gunakan base image dari Golang dengan versi 1.20.1
FROM golang:1.20.1-alpine

# Set working directory ke dalam folder "bin"
WORKDIR /app/bin

# Copy file main ke dalam container
COPY bin/main.go .

# Build aplikasi Go
RUN go mod tidy && \
    go build bin/main .

# Expose port 4000 untuk aplikasi
EXPOSE 4000

# Jalankan aplikasi saat container dijalankan
CMD ["./main"]
