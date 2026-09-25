#!/bin/sh
set -e

# Start the backend on port 8000 in background
echo "Starting backend on port 8000..."
./conspiracy-board > /tmp/backend.log 2>&1 &
BACKEND_PID=$!
echo "Backend PID: $BACKEND_PID"

# Wait for backend to start
sleep 3

# Verify backend is running
if ! kill -0 $BACKEND_PID 2>/dev/null; then
    echo "ERROR: Backend failed to start!"
    cat /tmp/backend.log
    exit 1
fi

echo "Backend started successfully"
echo "Starting nginx on port 8080..."

# Start nginx in the foreground (this becomes the main process)
exec nginx -c /etc/nginx/nginx.conf -g "daemon off;"
