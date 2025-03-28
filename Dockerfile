# Build stage
FROM golang:1.24.1 AS builder
WORKDIR /app

# Copy go.mod and go.sum for dependency installation
COPY go.mod go.sum ./
RUN go mod download

# Copy the actual source code
COPY cmd/app/ ./cmd/app/
WORKDIR /app/cmd/app

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app .

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add logrotate aws-cli && \
    mkdir -p /var/log/app

# Copy logrotate configuration
COPY configs/logrotate.conf /etc/logrotate.d/app

# Copy the built Go binary from the builder stage
COPY --from=builder /app/app /app/app

# Ensure correct permissions
RUN chmod +x /app/app

# Setup log rotation script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

WORKDIR /app
CMD ["/entrypoint.sh"]
