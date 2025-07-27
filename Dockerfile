# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o htb-mcp-server main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -s /bin/sh htb

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/htb-mcp-server .

# Change ownership to non-root user
RUN chown htb:htb /app/htb-mcp-server

# Switch to non-root user
USER htb

# Expose port (if needed for future HTTP transport)
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD echo '{"jsonrpc":"2.0","id":1,"method":"ping"}' | ./htb-mcp-server || exit 1

# Run the binary
ENTRYPOINT ["./htb-mcp-server"]