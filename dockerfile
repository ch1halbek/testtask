# Use a base image for building Go applications
FROM golang:1.23.0-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the application
RUN go build -o test_task ./cmd/main.go

# Create a minimal image to run the application
FROM alpine:latest

# Set the working directory for the final container
WORKDIR /root/

# Copy the binary file from the builder container
COPY --from=builder /app/test_task .

# Копируем конфигурационные файлы (например, config.yaml)
COPY --from=builder /app/internal/config/config.yaml ./internal/config/config.yaml

# Port to expose
EXPOSE 8080

# Specify the command to run the application
CMD ["./test_task"]