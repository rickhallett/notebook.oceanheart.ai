# Multi-stage Dockerfile for notebook.oceanheart.ai Go blog engine

# Build stage
FROM golang:1.22-bookworm AS build
RUN apt-get update && apt-get install -y build-essential sqlite3 libsqlite3-dev && rm -rf /var/lib/apt/lists/*
WORKDIR /src

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .
ENV CGO_ENABLED=1
RUN go build -ldflags="-s -w" -o /out/notebook ./cmd/notebook

# Runtime stage
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y \
    ca-certificates \
    sqlite3 \
    curl \
    && rm -rf /var/lib/apt/lists/* \
    && useradd -u 10001 -m app

WORKDIR /app

# Copy binary
COPY --from=build /out/notebook ./notebook

# Copy necessary files and directories
COPY content/ ./content/
COPY migrations/ ./migrations/
COPY internal/view/assets/ ./internal/view/assets/
COPY internal/view/templates/ ./internal/view/templates/

# Create directory for database with proper permissions
RUN mkdir -p /app/data && chown -R app:app /app

USER app

# Environment variables
ENV PORT=8080
ENV ENV=production
ENV DB_PATH=/app/data/notebook.db
ENV CONTENT_DIR=/app/content
ENV SITE_BASEURL=https://notebook.oceanheart.ai

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/ || exit 1

CMD ["./notebook"]
