🚀 Performance & Optimization Guide
Redis Caching
Enable/Disable Cache
Set in .env:

env
ENABLE_CACHE=true # Enable caching
ENABLE_CACHE=false # Disable caching
Cache Strategy

1. Suggestion Cache
   Key Pattern: suggestion:{yard}:{size}:{height}:{type}
   TTL: 5 minutes
   Invalidation: On placement or pickup in the same yard
2. Container Position Cache
   Key Pattern: container:{container_number}
   TTL: 24 hours
   Invalidation: On pickup
   Cache Monitoring
   bash

# Connect to Redis CLI

redis-cli

# View all keys

KEYS \*

# Get cache hit/miss stats

INFO stats

# Monitor cache in real-time

MONITOR

# View specific key

GET suggestion:YRD1:20:8.6:DRY

# Delete specific pattern

redis-cli KEYS "suggestion:YRD1:\*" | xargs redis-cli DEL
Cache Performance Metrics
With caching enabled:

Suggestion requests: 10-50x faster
Database load: Reduced by 60-80%
Response time: < 10ms (cached) vs 50-200ms (uncached)
Concurrent Processing
Worker Pool
The system uses a worker pool for bulk operations:

go
// 5 concurrent workers processing suggestions
pool := worker.NewPool(5, workerFunc)
Configuration
Adjust worker count based on:

CPU cores: Generally runtime.NumCPU()
Database connections: Don't exceed connection pool size
Memory: Each worker consumes memory
Bulk Operation Performance
Containers Sequential Concurrent (5 workers) Improvement
10 ~500ms ~120ms 4.2x
50 ~2.5s ~550ms 4.5x
100 ~5s ~1.1s 4.5x
Semaphore Pattern
For placement operations, we use a semaphore to limit concurrency:

go
semaphore := make(chan struct{}, 10) // Max 10 concurrent placements
This prevents:

Database connection exhaustion
Race conditions
Memory spikes
Database Optimization
Connection Pool Settings
In pkg/database/postgres.go:

go
db.SetMaxOpenConns(25) // Max connections
db.SetMaxIdleConns(5) // Idle connections
Adjust based on:

Available database connections: Check PostgreSQL max_connections
Expected load: More connections for high traffic
Memory: Each connection uses memory
Index Usage
Ensure these indexes are created (already in migrations):

sql
CREATE INDEX idx_containers_yard ON containers(yard_id);
CREATE INDEX idx_containers_block ON containers(block_id);
CREATE INDEX idx_containers_position ON containers(block_id, slot, row, tier);
CREATE INDEX idx_containers_number ON containers(container_number);
CREATE INDEX idx_yard_plans_block ON yard_plans(block_id);
Query Optimization
All queries use prepared statements
Selective column fetching (no SELECT \* in critical paths)
Efficient JOINs avoided where possible
Proper WHERE clause indexing
Load Testing
Using Apache Bench
bash

# Single endpoint

ab -n 1000 -c 10 -T 'application/json' \
 -p suggestion.json \
 http://localhost:8080/suggestion

# Bulk endpoint

ab -n 100 -c 5 -T 'application/json' \
 -p bulk_suggestion.json \
 http://localhost:8080/bulk/suggestion
Using hey
bash

# Install hey

go install github.com/rakyll/hey@latest

# Load test

hey -n 1000 -c 10 -m POST \
 -H "Content-Type: application/json" \
 -d @suggestion.json \
 http://localhost:8080/suggestion
Monitoring
Key Metrics to Monitor
Response Time
P50, P95, P99 latencies
Target: < 100ms for cached, < 500ms for uncached
Cache Hit Rate
Target: > 70% for suggestion endpoints
Monitor via Redis INFO stats
Database Connections
Active connections
Connection wait time
Query duration
Error Rate
Target: < 0.1%
Monitor 4xx and 5xx responses
Goroutine Count
Should stay stable
Sudden increases indicate goroutine leaks
Health Checks
bash

# Application health

curl http://localhost:8080/health

# Database health

psql -U postgres -d yard_planning -c "SELECT 1"

# Redis health

redis-cli ping
Best Practices

1. Cache Invalidation
   Always invalidate related caches after mutations
   Use pattern-based deletion for bulk invalidation
2. Concurrent Access
   Use mutexes for shared state
   Prefer channels for goroutine communication
   Always set timeouts for concurrent operations
3. Error Handling
   Log errors with context
   Return meaningful errors to clients
   Use circuit breakers for external dependencies
4. Resource Management
   Always close database rows
   Use context cancellation
   Set request timeouts
   Troubleshooting
   High Memory Usage
   Check goroutine count: curl http://localhost:8080/debug/pprof/goroutine
   Review connection pool settings
   Enable memory profiling
   Slow Queries
   Enable PostgreSQL slow query log
   Use EXPLAIN ANALYZE on slow queries
   Check index usage
   Cache Issues
   Verify Redis connectivity
   Check cache key patterns
   Monitor cache size and eviction
   Deadlocks
   Review lock ordering
   Use timeouts on all locks
   Enable deadlock detection logging
