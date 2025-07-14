# LinkSprint Quick Start Guide

Get LinkSprint running in 5 minutes! 🚀

## Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)
- Git

## 🚀 Quick Start with Docker

### 1. Clone the Repository
```bash
git clone <repository-url>
cd linksprint
```

### 2. Start All Services
```bash
docker-compose up -d
```

### 3. Access the Application
- **Web Interface**: http://localhost:8080
- **API Documentation**: http://localhost:8080/api/v1
- **Health Check**: http://localhost:8080/health

## 🛠️ Local Development

### 1. Install Dependencies
```bash
go mod download
```

### 2. Start Dependencies
```bash
docker-compose up -d redis cockroachdb
```

### 3. Run the Application
```bash
go run cmd/server/main.go
```

## 📊 Testing the API

### Create a Short URL
```bash
curl -X POST http://localhost:8080/api/v1/urls/shorten \
  -H "Content-Type: application/json" \
  -d '{
    "original_url": "https://google.com",
    "title": "Google Search"
  }'
```

### Get URL Statistics
```bash
curl http://localhost:8080/api/v1/urls/YOUR_SHORT_CODE/stats
```

### List All URLs
```bash
curl http://localhost:8080/api/v1/urls
```

### Get Global Analytics
```bash
curl http://localhost:8080/api/v1/analytics/global
```

## 🧪 Load Testing

Run the included load test script:

```bash
chmod +x scripts/load_test.sh
./scripts/load_test.sh
```

## 📈 Monitoring

### Prometheus
- **URL**: http://localhost:9090
- **Purpose**: Metrics collection and visualization

### Grafana
- **URL**: http://localhost:3000
- **Username**: admin
- **Password**: admin
- **Purpose**: Dashboard and alerting

## 🔧 Configuration

### Environment Variables
Create a `.env` file in the root directory:

```env
PORT=8080
ENV=development
COCKROACHDB_URL=postgresql://root@localhost:26257/linksprint?sslmode=disable
REDIS_URL=localhost:6379
JWT_SECRET=your-secret-key
```

### Docker Configuration
The `docker-compose.yml` file includes:
- **App**: LinkSprint application
- **Redis**: Caching layer
- **CockroachDB**: Database
- **Prometheus**: Metrics collection
- **Grafana**: Monitoring dashboard

## 🐛 Troubleshooting

### Common Issues

#### 1. Port Already in Use
```bash
# Check what's using the port
lsof -i :8080

# Kill the process
kill -9 <PID>
```

#### 2. Database Connection Issues
```bash
# Check if CockroachDB is running
docker-compose ps

# Restart the database
docker-compose restart cockroachdb
```

#### 3. Redis Connection Issues
```bash
# Check Redis logs
docker-compose logs redis

# Restart Redis
docker-compose restart redis
```

### Logs
```bash
# View all logs
docker-compose logs

# View specific service logs
docker-compose logs app
docker-compose logs redis
docker-compose logs cockroachdb
```

## 📚 Next Steps

### 1. Explore the Codebase
- **API Handlers**: `internal/handlers/`
- **Business Logic**: `internal/services/`
- **Data Models**: `internal/models/`
- **Database**: `internal/database/`

### 2. Run Tests
```bash
go test ./...
```

### 3. Build for Production
```bash
docker build -t linksprint .
```

### 4. Deploy
```bash
# Using Docker Compose
docker-compose -f docker-compose.prod.yml up -d

# Using Kubernetes (see k8s/ directory)
kubectl apply -f k8s/
```

## 🎯 Performance Testing

### Using wrk (High-Performance Load Testing)
```bash
# Install wrk
# macOS: brew install wrk
# Ubuntu: sudo apt-get install wrk

# Run load test
wrk -t12 -c400 -d30s http://localhost:8080/health
```

### Using Apache Bench
```bash
# Test URL creation
ab -n 1000 -c 10 -p test_data.json -T application/json http://localhost:8080/api/v1/urls/shorten
```

## 📖 API Documentation

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/urls/shorten` | Create short URL |
| GET | `/api/v1/urls` | List all URLs |
| GET | `/api/v1/urls/:shortCode/stats` | Get URL statistics |
| GET | `/api/v1/analytics/:shortCode` | Get URL analytics |
| GET | `/api/v1/analytics/global` | Get global analytics |
| GET | `/:shortCode` | Redirect to original URL |

### Request Examples

#### Create URL
```json
{
  "original_url": "https://example.com",
  "title": "Example Website",
  "description": "A sample website",
  "custom_code": "my-link"
}
```

#### Response
```json
{
  "short_code": "abc123",
  "original_url": "https://example.com",
  "short_url": "http://localhost:8080/abc123",
  "created_at": "2024-01-01T12:00:00Z"
}
```

## 🚀 Production Deployment

### 1. Environment Setup
```bash
# Set production environment variables
export ENV=production
export JWT_SECRET=your-super-secret-key
export BASE_URL=https://your-domain.com
```

### 2. Database Setup
```bash
# Initialize CockroachDB cluster
cockroach start-single-node --insecure
cockroach sql --insecure -e "CREATE DATABASE linksprint;"
```

### 3. Deploy with Docker
```bash
docker-compose -f docker-compose.prod.yml up -d
```

### 4. Monitor
```bash
# Check application health
curl https://your-domain.com/health

# View logs
docker-compose logs -f app
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📞 Support

- **Issues**: GitHub Issues
- **Documentation**: `/docs` directory
- **API**: http://localhost:8080/api/v1

---

**Happy Coding! 🎉** 