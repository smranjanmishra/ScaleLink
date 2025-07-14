# LinkSprint System Design

## 🏗️ Architecture Overview

LinkSprint is a distributed URL shortener designed to handle millions of redirects with high availability and low latency. The system follows a microservices architecture with clear separation of concerns.

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

## 🎯 Design Goals

1. **High Availability**: 99.9% uptime
2. **Low Latency**: < 50ms average response time
3. **High Throughput**: 1M+ requests/day
4. **Fault Tolerance**: Graceful degradation
5. **Scalability**: Horizontal scaling capability

## 🏛️ System Components

### 1. API Layer (Go + Fiber)
- **Purpose**: HTTP request handling and routing
- **Technology**: Go with Fiber framework
- **Features**: 
  - RESTful API endpoints
  - Middleware for logging, CORS, rate limiting
  - JSON request/response handling
  - Error handling and validation

### 2. Caching Layer (Redis)
- **Purpose**: Fast URL lookups and session storage
- **Technology**: Redis with persistence
- **Features**:
  - URL caching with TTL
  - Click count tracking
  - Session storage
  - Rate limiting counters

### 3. Database Layer (CockroachDB)
- **Purpose**: Persistent data storage
- **Technology**: CockroachDB (distributed SQL)
- **Features**:
  - ACID transactions
  - Automatic sharding
  - Geographic distribution
  - Built-in replication

## 📊 Data Models

### URL Table
```sql
CREATE TABLE urls (
    id UUID PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    title VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    expires_at TIMESTAMP,
    INDEX idx_short_code (short_code),
    INDEX idx_created_at (created_at)
);
```

### Analytics Table
```sql
CREATE TABLE analytics (
    id UUID PRIMARY KEY,
    url_id UUID NOT NULL,
    short_code VARCHAR(10) NOT NULL,
    ip_address INET,
    user_agent TEXT,
    referer TEXT,
    country VARCHAR(100),
    city VARCHAR(100),
    clicked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_url_id (url_id),
    INDEX idx_short_code (short_code),
    INDEX idx_clicked_at (clicked_at),
    FOREIGN KEY (url_id) REFERENCES urls(id) ON DELETE CASCADE
);
```

## 🔄 Request Flow

### URL Creation Flow
1. Client sends POST request with original URL
2. API validates URL format and length
3. Generate unique short code (6 characters)
4. Store URL in CockroachDB
5. Cache URL in Redis with 24-hour TTL
6. Return short URL to client

### URL Redirection Flow
1. Client requests short URL
2. Check Redis cache first (fast path)
3. If not in cache, query CockroachDB
4. Cache result in Redis
5. Track analytics (async)
6. Redirect to original URL

### Analytics Flow
1. Track click event in database
2. Increment Redis counter
3. Aggregate data for reporting
4. Cache analytics results

## 🚀 Performance Optimizations

### 1. Caching Strategy
- **Hot URLs**: Cached in Redis with 24-hour TTL
- **Analytics**: Aggregated data cached for 5 minutes
- **User Sessions**: JWT tokens with Redis storage

### 2. Database Optimizations
- **Indexes**: On short_code, created_at, clicked_at
- **Connection Pooling**: Reuse database connections
- **Query Optimization**: Use prepared statements

### 3. Application Optimizations
- **Connection Pooling**: Redis and database connections
- **Async Processing**: Analytics tracking
- **Rate Limiting**: Per-IP request limits

## 🔧 CAP Theorem Trade-offs

### Consistency vs Availability
- **Analytics Data**: Eventual consistency (AP)
- **URL Lookups**: Strong consistency (CP)
- **Cache**: Eventual consistency with TTL

### Partition Tolerance
- **CockroachDB**: Handles network partitions automatically
- **Redis**: Master-slave replication
- **Application**: Stateless design

## 📈 Scaling Strategy

### Horizontal Scaling
1. **Load Balancer**: Distribute requests across multiple API servers
2. **Stateless API**: Any server can handle any request
3. **Database Sharding**: CockroachDB handles automatic sharding
4. **Cache Distribution**: Redis cluster for high availability

### Vertical Scaling
1. **CPU**: Multi-core processing
2. **Memory**: Increased cache capacity
3. **Storage**: SSD for database performance

## 🛡️ Fault Tolerance

### Single Points of Failure
1. **Database**: CockroachDB with automatic failover
2. **Cache**: Redis with replication
3. **Application**: Multiple instances behind load balancer

### Graceful Degradation
1. **Cache Failure**: Fallback to database
2. **Database Failure**: Serve cached data only
3. **Analytics Failure**: Continue core functionality

## 🔒 Security Considerations

### Input Validation
- URL format validation
- SQL injection prevention
- XSS protection
- Rate limiting

### Data Protection
- HTTPS for all communications
- JWT for authentication
- Secure headers
- Input sanitization

## 📊 Monitoring & Observability

### Metrics
- Request latency
- Error rates
- Cache hit ratio
- Database performance
- System resources

### Logging
- Structured logging
- Request tracing
- Error tracking
- Performance monitoring

## 🧪 Testing Strategy

### Unit Tests
- Service layer logic
- Database operations
- Cache operations
- Input validation

### Integration Tests
- API endpoints
- Database integration
- Cache integration
- End-to-end flows

### Load Tests
- Concurrent users
- High throughput
- Stress testing
- Performance benchmarks

## 🚀 Deployment Strategy

### Docker Containerization
- Multi-stage builds
- Health checks
- Resource limits
- Environment configuration

### CI/CD Pipeline
- Automated testing
- Code quality checks
- Security scanning
- Automated deployment

## 📈 Performance Targets

| Metric | Target | Current |
|--------|--------|---------|
| Response Time | < 50ms | TBD |
| Throughput | 1M req/day | TBD |
| Cache Hit Ratio | > 90% | TBD |
| Availability | 99.9% | TBD |
| Error Rate | < 0.1% | TBD |

## 🔮 Future Enhancements

### Phase 2 Features
- User authentication
- Custom domains
- Advanced analytics
- API rate limiting
- Webhook notifications

### Phase 3 Features
- Geographic distribution
- CDN integration
- Advanced caching
- Machine learning insights
- Mobile SDK

## 📚 References

- [CockroachDB Documentation](https://www.cockroachlabs.com/docs/)
- [Redis Documentation](https://redis.io/documentation)
- [Fiber Framework](https://gofiber.io/)
- [System Design Primer](https://github.com/donnemartin/system-design-primer) 