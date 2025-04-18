#!/bin/bash

# Build all services
echo "Building crawler service..."
go build -o bin/crawler services/crawler/main.go

echo "Building product-analyzer service..."
go build -o bin/product-analyzer services/product-analyzer/main.go

echo "Building notification service..."
go build -o bin/notification services/notification/main.go
