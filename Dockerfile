# Vulnerable MCP Server Dockerfile
# WARNING: This container is intentionally vulnerable for educational purposes only

FROM python:3.11-slim

# Create non-root user for security
RUN groupadd -r mcpuser && useradd -r -g mcpuser mcpuser

# Set working directory
WORKDIR /app

# Copy requirements and install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application files
COPY vulnerable_mcp_server.py .
COPY test_vulnerabilities.py .
COPY sandbox/ ./sandbox/

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
    CMD python -c "import sys; sys.exit(0)"

# Run the vulnerable server
CMD ["python", "vulnerable_mcp_server.py"]

# Security labels
LABEL security.warning="This container contains intentional vulnerabilities"
LABEL security.purpose="Educational use only"
LABEL security.production="DO NOT USE IN PRODUCTION"