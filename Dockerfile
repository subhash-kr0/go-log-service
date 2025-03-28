# Build stage
FROM golang:1.24.1 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add logrotate aws-cli && \
    mkdir -p /var/log

# Copy logrotate config
COPY logrotate.conf /etc/logrotate.d/app

# Copy built binary
COPY --from=builder /app/app /app/

# Setup cron job for log rotation
RUN echo "0 0 * * * /usr/sbin/logrotate -f /etc/logrotate.d/app" >> /etc/crontabs/root

WORKDIR /app
CMD ["sh", "-c", "crond && /app/app"]