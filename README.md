# LinkSprint: Distributed URL Shortener & Analytics

A high-performance, distributed URL shortener built with Go, featuring analytics, caching, and fault tolerance.

## 🚀 Features

- **URL Shortening**: Generate short, unique URLs
- **Analytics Dashboard**: Track clicks, geographic data, and performance metrics
- **High Performance**: Redis caching for sub-50ms response times
- **Fault Tolerant**: Distributed architecture with CockroachDB
- **Scalable**: Designed to handle millions of redirects
- **Docker Ready**: Easy deployment with containerization

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Load Balancer │    │   Go + Fiber    │    │   Redis Cache   │
│   (Nginx)       │───▶│   API Server    │───▶│   (Hot Data)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │  CockroachDB    │
                       │  (Persistent)   │
                       └─────────────────┘
```

## 🛠️ Tech Stack

- **Backend**: Go + Fiber (high-performance web framework)
- **Cache**: Redis (for fast URL lookups)
- **Database**: CockroachDB (distributed SQL database)
- **Containerization**: Docker + Docker Compose
- **Monitoring**: Prometheus + Grafana (optional)

## 📊 Performance Targets

- **Latency**: < 50ms average response time
- **Throughput**: 1M+ requests/day
- **Cache Hit Ratio**: > 90%
- **Availability**: 99.9% uptime

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)

### Running with Docker

```bash
# Clone the repository
git clone <repository-url>
cd linksprint

# Start all services
docker-compose up -d

# Access the application
open http://localhost:8080
```

### Local Development

```bash
# Install dependencies
go mod download

# Start Redis and CockroachDB
docker-compose up -d redis cockroachdb

# Run the application
go run cmd/server/main.go
```

## 📈 API Endpoints

### URL Shortening
- `POST /api/v1/shorten` - Create a short URL
- `GET /:shortCode` - Redirect to original URL
- `GET /api/v1/urls` - List all URLs (with pagination)

### Analytics
- `GET /api/v1/analytics/:shortCode` - Get analytics for a URL
- `GET /api/v1/analytics/global` - Global analytics dashboard

### Health & Monitoring
- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics

## 🔧 Configuration

Environment variables:

```env
# Server
PORT=8080
ENV=development

# Database
COCKROACHDB_URL=postgresql://root@localhost:26257/linksprint?sslmode=disable

# Redis
REDIS_URL=localhost:6379

# Security
JWT_SECRET=your-secret-key
```

## 📊 Load Testing

Run load tests to verify performance:

```bash
# Install wrk
brew install wrk  # macOS
# or download from https://github.com/wg/wrk

# Run load test
wrk -t12 -c400 -d30s http://localhost:8080/health
```

## 🏗️ System Design Decisions

### CAP Theorem Trade-offs
- **Consistency**: Eventual consistency for analytics data
- **Availability**: High availability through Redis caching
- **Partition Tolerance**: CockroachDB handles network partitions

### Caching Strategy
- **Hot URLs**: Cached in Redis with TTL
- **Analytics**: Aggregated data cached for 5 minutes
- **User Sessions**: JWT tokens with Redis storage

### Scaling Considerations
- **Horizontal Scaling**: Stateless API servers behind load balancer
- **Database Sharding**: CockroachDB handles automatic sharding
- **Cache Distribution**: Redis cluster for high availability
