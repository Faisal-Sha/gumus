# Trendyol Product Tracking System

A scalable microservices-based system for tracking and analyzing product information from Trendyol in real-time.

## Architecture

The system consists of three main microservices:

1. **Crawler Service**: Responsible for fetching product data from Trendyol
2. **Product Analysis Service**: Handles product data analysis, categorization, and database operations
3. **Notification Service**: Manages price drop notifications for favorited products

## Tech Stack

- **Language**: Go (Echo framework)
- **Database**: PostgreSQL
- **Message Queue**: Kafka/RabbitMQ
- **Service Communication**: gRPC
- **ORM**: GORM

## Project Structure

```
.
├── services/
│   ├── crawler/         # Product crawling service
│   ├── product-analyzer/ # Product analysis service
│   └── notification/    # Notification service
├── proto/              # Protocol buffer definitions
├── pkg/
│   ├── models/         # Shared data models
│   ├── database/       # Database configurations
│   └── config/         # Configuration management
└── scripts/           # Utility scripts
```

## Features

- Real-time product data synchronization
- New product detection
- Inactive product marking
- Price and stock tracking
- Prioritized updates for favorited products
- Price drop notifications
- Comprehensive product information tracking

## Setup Instructions

[To be added after initial development]

## Development

[To be added after initial development]

## Testing

[To be added after initial development]
