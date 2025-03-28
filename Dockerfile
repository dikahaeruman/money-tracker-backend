# Use the official Golang image as the builder
FROM golang:1.21 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project and build the binary
COPY . .
RUN go build -o main .

# Use a lightweight image for the final container
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the compiled binary from the builder
COPY --from=builder /app/main .

# Expose the application's port
EXPOSE 8080

# Run the application
CMD ["./main"]
