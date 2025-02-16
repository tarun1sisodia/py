#!/bin/bash

set -e

# Start MySQL service if not already running
echo "Starting MySQL service..."
if systemctl is-active --quiet mysql; then
    echo "MySQL service is already running."
else
    echo "Starting MySQL..."
    sudo systemctl start mysql
fi

# Build the Go server
echo "Building the Go server..."
cd backend/cmd/server && go build -o server .

# Start the Go server in the background from backend directory so that .env is found
echo "Starting the Go server..."
cd ../../ && ./cmd/server/server &

echo "Deployment complete!" 