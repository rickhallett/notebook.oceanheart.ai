## Project Overview

This project is a minimalist, single-binary blog engine built with Go. It uses SQLite for the database and supports Turso (libSQL) for remote databases. The content is written in Markdown with YAML front matter. The engine is designed for speed, simplicity, and easy deployment. It includes features like syntax highlighting, SEO optimization, and Atom feeds.

## Building and Running

### Prerequisites
- Go 1.22.7 or later
- CGO enabled for SQLite support

### Local Development

To run the development server with hot-reloading:
```bash
./scripts/dev.sh
```
This script requires one of `watchexec`, `reflex`, or `entr` to be installed.

To run the server manually:
```bash
# Build the binary
go build -o notebook ./cmd/notebook

# Run in development mode
ENV=dev ./notebook
```

### Production Build

To build an optimized production binary:
```bash
CGO_ENABLED=1 go build -ldflags "-s -w" -o notebook ./cmd/notebook
```

The project can also be deployed using Docker or to Fly.io.

### Testing

To run all tests:
```bash
go test ./...
```

To run tests with coverage:
```bash
go test -v -cover ./...
```

## Development Conventions

- **Code Style**: The code is formatted with `go fmt`.
- **Linting**: The code is vetted with `go vet`.
- **Dependencies**: Go modules are used for dependency management.
- **Testing**: The project has extensive test coverage across all packages.
- **Content**: Blog posts are created as Markdown files in the `/content/` directory with the filename pattern `YYYY-MM-DD-slug.md`.
- **Configuration**: Configuration is managed through environment variables.
- **Deployment**: The project is configured for deployment on Fly.io via `fly.toml`.
