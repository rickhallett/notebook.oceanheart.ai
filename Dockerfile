# Multi-stage Dockerfile for notebook.oceanheart.ai Go blog engine

# Build stage
FROM golang:1.22-bookworm AS build
RUN apt-get update && apt-get install -y build-essential sqlite3 libsqlite3-dev && rm -rf /var/lib/apt/lists/*
WORKDIR /src
COPY . .
ENV CGO_ENABLED=1
RUN go build -o /out/notebook ./cmd/notebook

# Runtime stage
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/* && useradd -u 10001 -m app
WORKDIR /app
COPY --from=build /out/notebook ./notebook
COPY content/ ./content/
COPY internal/view/assets/ ./internal/view/assets/
COPY internal/view/templates/ ./internal/view/templates/
USER app
ENV PORT=8080
EXPOSE 8080
CMD ["./notebook"]
