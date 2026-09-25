# Docker Setup Guide

Single-container deployment with both frontend and backend.

## Files Included

- **Dockerfile** - Multi-stage build (frontend + backend in one container)
- **.dockerignore** - Files to exclude from Docker builds

## Quick Start

### Build the image
```bash
docker build -t conspiracy-board:latest .
```

### Run the container
```bash
# Basic run
docker run -p 8080:8080 conspiracy-board:latest

# With database persistence (named volume)
docker run -p 8080:8080 -v conspiracy-data:/data conspiracy-board:latest

# Run in background
docker run -d -p 8080:8080 -v conspiracy-data:/data conspiracy-board:latest
```

The application will be available at:
- Frontend: http://localhost:8080
- Backend API: http://localhost:8080/api
- Health check: http://localhost:8080/health

## How It Works

1. **Stage 1** - Builds React frontend with Vite
2. **Stage 2** - Builds Go backend binary
3. **Stage 3** - Creates minimal Alpine runtime with:
   - Go backend binary
   - Frontend static files in `./public` directory
   - SQLite database

The backend Gin server serves both the API and static frontend files.

## Container Management

```bash
# View running containers
docker ps

# View logs
docker logs <container-id>

# Stop container
docker stop <container-id>

# Remove container
docker rm <container-id>

# Remove image
docker rmi conspiracy-board:latest
```

## Environment Variables

- `GIN_MODE` - Set to `release` for production (optional)
- `DB_PATH` - SQLite database path (default in container: `/data/data.db`)

Example:
```bash
docker run -e GIN_MODE=release -p 8080:8080 conspiracy-board:latest
```

## Database Persistence

The database is stored at `/data/data.db` (configurable via `DB_PATH`). Mount a volume or
host **directory** at `/data` to persist it (SQLite WAL mode also creates `-wal`/`-shm`
files next to the database, so mount the directory, not the single file):
```bash
docker run -v conspiracy-data:/data -p 8080:8080 conspiracy-board:latest
# or a host directory
docker run -v /path/to/data:/data -p 8080:8080 conspiracy-board:latest
```

## Development Workflow

For local development without Docker:
```bash
# Terminal 1: Run backend
cd backend
go run main.go database.go handlers.go models.go

# Terminal 2: Run frontend
cd frontend
npm install
npm run dev
```

Access at:
- Frontend: http://localhost:5173 (Vite dev server)
- Backend: http://localhost:8080
