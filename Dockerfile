# Stage 1: Build frontend
FROM node:24-alpine3.24 AS frontend-builder
WORKDIR /app/frontend

COPY frontend .
RUN npm ci
RUN npm run build

# Stage 2: Build backend
FROM golang:1.27-alpine3.24 AS backend-builder
WORKDIR /app/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend .

# Build the binary (modernc.org/sqlite is pure Go, so no cgo toolchain is needed)
RUN CGO_ENABLED=0 GOOS=linux go build -o conspiracy-board .

# Stage 3: Runtime
FROM alpine:3.24
RUN apk add --no-cache ca-certificates nginx

WORKDIR /app

# Create data directory
RUN mkdir -p /data
VOLUME /data

# Set environment variables
ENV DB_PATH=/data/data.db

# Copy nginx config
COPY nginx.conf /etc/nginx/nginx.conf

# Copy entrypoint script (strip CR in case it was checked out with CRLF line endings)
COPY entrypoint.sh /app/entrypoint.sh
RUN sed -i 's/\r$//' /app/entrypoint.sh && chmod +x /app/entrypoint.sh

# Copy backend binary
COPY --from=backend-builder /app/backend/conspiracy-board .

# Copy frontend static files
COPY --from=frontend-builder /app/frontend/dist ./public

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
