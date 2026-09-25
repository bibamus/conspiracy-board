# Stage 1: Build frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend

COPY frontend .
RUN npm ci
RUN npm run build

# Stage 2: Build backend
FROM golang:1.27-alpine AS backend-builder
WORKDIR /app/backend

RUN apk add --no-cache gcc musl-dev

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend .

# Build the binary
RUN CGO_ENABLED=1 GOOS=linux go build -o conspiracy-board .

# Stage 3: Runtime
FROM alpine:latest
RUN apk add --no-cache ca-certificates nginx

WORKDIR /app

# Create data directory
RUN mkdir -p /data

# Set environment variables
ENV DB_PATH=/data/data.db

# Copy nginx config
COPY nginx.conf /etc/nginx/nginx.conf

# Create entrypoint script
RUN printf '#!/bin/sh\nset -e\n\necho "Starting backend on port 8000..."\n./conspiracy-board > /tmp/backend.log 2>&1 &\nBACKEND_PID=$!\necho "Backend PID: $BACKEND_PID"\n\nsleep 3\n\nif ! kill -0 $BACKEND_PID 2>/dev/null; then\n    echo "ERROR: Backend failed to start!"\n    cat /tmp/backend.log\n    exit 1\nfi\n\necho "Backend started successfully"\necho "Starting nginx on port 8080..."\n\nexec nginx -c /etc/nginx/nginx.conf -g "daemon off;"\n' > /app/entrypoint.sh && chmod +x /app/entrypoint.sh

# Copy backend binary
COPY --from=backend-builder /app/backend/conspiracy-board .

# Copy frontend static files
COPY --from=frontend-builder /app/frontend/dist ./public

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
