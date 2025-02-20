# Use the latest Go version that matches your go.mod requirement
FROM golang:1.24 AS builder

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go application
RUN go build -o main .

# Use a minimal base image for the final container
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy compiled binary from builder
COPY --from=builder /app/main .

# Expose application port
EXPOSE 3000

# Set environment variables
ENV MONGO_URI=mongodb://mongo:27017 \
    JWT_SECRET=your_jwt_secret

# Copy the .env file to the container
COPY .env .env

# Run the application
CMD ["./main"]
