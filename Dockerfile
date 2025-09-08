# Vulnerable MCP Server Dockerfile (Go)
# WARNING: This container is intentionally vulnerable for educational purposes only

# Build stage
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY vulnerable_mcp_server.go .
COPY test_vulnerabilities.go .
COPY sandbox/ ./sandbox/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o vulnerable_mcp_server vulnerable_mcp_server.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o test_vulnerabilities test_vulnerabilities.go

# Final stage
FROM alpine:latest

# Install ca-certificates and shell for command execution
RUN apk --no-cache add ca-certificates bash

# Create non-root user for security
RUN addgroup -g 1000 mcpuser && adduser -D -s /bin/bash -u 1000 -G mcpuser mcpuser

# Set working directory
WORKDIR /app

# Copy built binaries from builder stage
COPY --from=builder /app/vulnerable_mcp_server .
COPY --from=builder /app/test_vulnerabilities .
COPY --from=builder /app/sandbox/ ./sandbox/

# Create sandbox directory with proper permissions
RUN mkdir -p /app/sandbox && \
    chown -R mcpuser:mcpuser /app && \
    chmod -R 755 /app/sandbox

# Switch to non-root user
USER mcpuser

# Expose port (if needed for future web interface)
EXPOSE 8000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ./vulnerable_mcp_server --version || exit 1

# Run the vulnerable server
CMD ["./vulnerable_mcp_server"]

# Security labels
LABEL security.warning="This container contains intentional vulnerabilities"
LABEL security.purpose="Educational use only"
LABEL security.production="DO NOT USE IN PRODUCTION"
