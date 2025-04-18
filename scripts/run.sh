#!/bin/bash

# Start all services in the background
echo "Starting product-analyzer service..."
./bin/product-analyzer &
ANALYZER_PID=$!

echo "Starting notification service..."
./bin/notification &
NOTIFICATION_PID=$!

echo "Starting crawler service..."
./bin/crawler &
CRAWLER_PID=$!

# Wait for all services
wait $ANALYZER_PID $NOTIFICATION_PID $CRAWLER_PID
