# Build Stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to leverage Docker cache during builds
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Copy the migrations folder as well
COPY migrations /app/migrations

# Build the Go application
RUN go build -o myapp main.go

# Production Stage
FROM alpine:latest

WORKDIR /app

# Install any necessary dependencies (e.g., certificates for HTTPS)
RUN apk add --no-cache ca-certificates

# Copy the compiled binary from the builder stage
COPY --from=builder /app/myapp /app/myapp

# Copy the migrations folder from the builder stage (if not copied in previous step)
COPY --from=builder /app/migrations /app/migrations

# Expose the application port
EXPOSE 5555

# Add health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:5555/api/health || exit 1

# Set environment variables here, but can also be overridden at runtime
ENV APP_ENV=production
ENV PORT=5555

# Run the Go application with environment variables
CMD ["/app/myapp"]
