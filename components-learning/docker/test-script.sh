#!/usr/bin/env bash

set -e  # Exit on any error

echo "Starting the service..."
docker-compose up --build -d

echo "This is the automation script. We need to check the services to be healthy before further testing"
echo "Waiting for services to start..."
sleep 10

echo "Testing the service health check..."
if curl -f http://localhost:5000/health; then
    echo "Health check passed!"
else
    echo "Health check failed!"
    exit 1
fi

echo "Testing resource usage..."
# Run docker stats for a limited time since it runs continuously
sleep 10s docker stats --no-stream

echo "Checking logs before shutdown..."
docker-compose logs cms

echo "Testing graceful shutdown..."
docker-compose kill -s SIGTERM cms

echo "Waiting for graceful shutdown to complete..."
sleep 5

echo "Checking final logs..."
docker-compose logs cms

echo "Downing the entire suite..."
docker-compose down

echo "Script completed successfully!"