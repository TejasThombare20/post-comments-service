# Post-Comments Service Architecture

## 📋 Table of Contents

1. [System Overview](#system-overview)
2. [Architecture Patterns](#architecture-patterns)
3. [Database Design](#database-design)
4. [Nested Comments System](#nested-comments-system)
5. [Performance Optimizations](#performance-optimizations)
6. [Security Architecture](#security-architecture)
7. [API Design](#api-design)
8. [Deployment Architecture](#deployment-architecture)
9. [Monitoring & Observability](#monitoring--observability)
10. [Scalability Considerations](#scalability-considerations)

---

## 🏗️ System Overview

The Post-Comments Service is a **RESTful backend service** built using **Clean Architecture** principles with **Go (Golang)**, designed to handle user management, post creation, and **efficient nested comment systems** at scale.

### Core Components

```
┌─────────────────────────────────────────────────────────────┐
│                    Client Applications                      │
│              (Web, Mobile, Third-party APIs)               │
└─────────────────────┬───────────────────────────────────────┘
                      │ HTTP/HTTPS
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                   Load Balancer                            │
│                 (nginx/HAProxy)                             │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                 Go Backend Service                          │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │ Controllers │ │ Middleware  │ │   Routes    │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │  Services   │ │ Validators  │ │  Utilities  │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
│  ┌─────────────┐ ┌─────────────┐                           │
│  │ Repository  │ │   Models    │                           │
│  └─────────────┘ └─────────────┘                           │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                  Data Layer                                 │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │ PostgreSQL  │ │    Redis    │ │ File Storage│           │
│  │  (Primary)  │ │   (Cache)   │ │ (Optional)  │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────────────────────────────────────────────┘
```

---

## 🏛️ Architecture Patterns

### Clean Architecture Implementation

The service follows **Clean Architecture** principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                    External Layer                           │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │ Controllers │ │ Middleware  │ │   Routes    │           │
│  │ (HTTP)      │ │ (Auth/CORS) │ │ (Gin)       │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────┬───────────────────────────────────────┘
                      │ Dependency Injection
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                  Application Layer                          │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │  Services   │ │ Validators  │ │    DTOs     │           │
│  │ (Business)  │ │ (Input)     │ │ (Transfer)  │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────┬───────────────────────────────────────┘
                      │ Interface Contracts
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                   Domain Layer                              │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │   Models    │ │ Interfaces  │ │   Errors    │           │
│  │ (Entities)  │ │ (Contracts) │ │ (Custom)    │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────┬───────────────────────────────────────┘
                      │ Implementation
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                Infrastructure Layer                         │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │ Repository  │ │  Database   │ │   Config    │           │
│  │ (Data)      │ │ (PostgreSQL)│ │ (Environment)│          │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────────────────────────────────────────────┘
```

### Key Architectural Benefits

1. **Testability**: Each layer can be tested independently
2. **Maintainability**: Clear separation of concerns
3. **Scalability**: Easy to scale individual components
4. **Flexibility**: Easy to swap implementations
5. **Dependency Inversion**: High-level modules don't depend on low-level modules

---

## 🗄️ Database Design

### Entity Relationship Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                        Users                                │
├─────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                               │
│ username (TEXT, UNIQUE, NOT NULL)                           │
│ email (TEXT, UNIQUE)                                        │
│ password_hash (TEXT)                                        │
│ display_name (TEXT)                                         │
│ avatar_url (TEXT)                                           │
│ created_at (TIMESTAMP)                                      │
│ updated_at (TIMESTAMP)                                      │
│ deleted_at (TIMESTAMP)                                      │
└─────────────────────┬───────────────────────────────────────┘
                      │ 1:N
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                        Posts                                │
├─────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                               │
│ title (TEXT, NOT NULL)                                      │
│ content (TEXT, NOT NULL)                                    │
│ created_by (UUID, FK → users.id)                            │
│ created_at (TIMESTAMP)                                      │
│ updated_at (TIMESTAMP)                                      │
│ deleted_at (TIMESTAMP)                                      │
└─────────────────────┬───────────────────────────────────────┘
                      │ 1:N
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                      Comments                               │
├─────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                               │
│ content (TEXT, NOT NULL)                                    │
│ post_id (UUID, FK → posts.id)                               │
│ parent_id (UUID, FK → comments.id, NULLABLE)                │
│ path (UUID[], Array for materialized path)                  │
│ thread_id (UUID, Root comment ID)                           │
│ created_by (UUID, FK → users.id)                            │
│ created_at (TIMESTAMP)                                      │
│ deleted_at (TIMESTAMP)                                      │
└─────────────────────────────────────────────────────────────┘
```

### Database Indexes for Performance

```sql
-- Users table indexes
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Posts table indexes
CREATE INDEX idx_posts_created_by ON posts(created_by);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
CREATE INDEX idx_posts_title ON posts USING gin(to_tsvector('english', title));

-- Comments table indexes (Critical for nested comments performance)
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id);
CREATE INDEX idx_comments_thread_id ON comments(thread_id);
CREATE INDEX idx_comments_path ON comments USING gin(path);
CREATE INDEX idx_comments_created_at ON comments(created_at);
CREATE INDEX idx_comments_post_thread ON comments(post_id, thread_id);
```

---

## 🌳 Nested Comments System

### The Challenge

Traditional approaches to nested comments have significant limitations:

1. **Adjacency List**: Simple but requires recursive queries
2. **Nested Sets**: Fast reads but complex writes
3. **Closure Table**: Flexible but storage overhead

### Our Solution: Materialized Path + Thread ID

We use a **hybrid approach** combining **Materialized Path** with **Thread ID** for optimal performance:

```sql
CREATE TABLE comments (
    id UUID PRIMARY KEY,
    content TEXT NOT NULL,
    post_id UUID NOT NULL REFERENCES posts(id),
    parent_id UUID REFERENCES comments(id),
    path UUID[] NOT NULL,           -- Materialized path
    thread_id UUID NOT NULL,        -- Root comment ID
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);
```

### How It Works

#### 1. Root Comment Creation
```sql
-- Creating a root comment
INSERT INTO comments (id, content, post_id, parent_id, path, thread_id, created_by)
VALUES (
    'comment-1-uuid',
    'This is a root comment',
    'post-uuid',
    NULL,
    ARRAY['comment-1-uuid'],  -- Path contains only itself
    'comment-1-uuid',         -- Thread ID is itself
    'user-uuid'
);
```

#### 2. Nested Comment Creation
```sql
-- Creating a reply to comment-1
INSERT INTO comments (id, content, post_id, parent_id, path, thread_id, created_by)
VALUES (
    'comment-2-uuid',
    'This is a reply',
    'post-uuid',
    'comment-1-uuid',
    ARRAY['comment-1-uuid', 'comment-2-uuid'],  -- Path includes parent path
    'comment-1-uuid',                           -- Same thread ID as root
    'user-uuid'
);

-- Creating a nested reply (comment-2 → comment-3)
INSERT INTO comments (id, content, post_id, parent_id, path, thread_id, created_by)
VALUES (
    'comment-3-uuid',
    'This is a nested reply',
    'post-uuid',
    'comment-2-uuid',
    ARRAY['comment-1-uuid', 'comment-2-uuid', 'comment-3-uuid'],
    'comment-1-uuid',
    'user-uuid'
);
```

### Query Patterns

#### 1. Get All Comments for a Post (Hierarchical)
```sql
SELECT 
    c.*,
    u.username,
    u.display_name,
    u.avatar_url,
    array_length(c.path, 1) as depth
FROM comments c
JOIN users u ON c.created_by = u.id
WHERE c.post_id = $1 
  AND c.deleted_at IS NULL
ORDER BY c.thread_id, c.path;
```

#### 2. Get Comment Thread (All replies to a specific comment)
```sql
SELECT 
    c.*,
    u.username,
    u.display_name
FROM comments c
JOIN users u ON c.created_by = u.id
WHERE c.thread_id = $1 
  AND c.deleted_at IS NULL
ORDER BY c.path;
```

#### 3. Get Direct Replies to a Comment
```sql
SELECT 
    c.*,
    u.username,
    u.display_name
FROM comments c
JOIN users u ON c.created_by = u.id
WHERE c.parent_id = $1 
  AND c.deleted_at IS NULL
ORDER BY c.created_at;
```

#### 4. Get Comment with All Ancestors
```sql
SELECT 
    c.*,
    u.username,
    u.display_name
FROM comments c
JOIN users u ON c.created_by = u.id
WHERE c.id = ANY(
    SELECT unnest(path) 
    FROM comments 
    WHERE id = $1
)
ORDER BY array_position(
    (SELECT path FROM comments WHERE id = $1), 
    c.id
);
```

### Performance Benefits

#### 1. **Fast Reads** ⚡
- **Single Query**: Get entire comment tree with one query
- **No Recursion**: Avoid expensive recursive CTEs
- **Index Optimized**: GIN index on path array for fast lookups

#### 2. **Efficient Writes** 📝
- **Simple Inserts**: Just append to parent's path
- **No Complex Updates**: Path calculation is straightforward
- **Atomic Operations**: Each comment insert is independent

#### 3. **Scalable Structure** 📈
- **Unlimited Nesting**: No depth limitations
- **Thread Isolation**: Comments grouped by thread_id
- **Pagination Friendly**: Easy to paginate by thread or depth

### Example Comment Tree Structure

```
Post: "How to optimize database queries?"
│
├── Comment 1 (thread_id: comment-1, path: [comment-1])
│   "Great question! Here are some tips..."
│   │
│   ├── Comment 2 (thread_id: comment-1, path: [comment-1, comment-2])
│   │   "Thanks for the tips! What about indexing?"
│   │   │
│   │   └── Comment 3 (thread_id: comment-1, path: [comment-1, comment-2, comment-3])
│   │       "For indexing, consider..."
│   │
│   └── Comment 4 (thread_id: comment-1, path: [comment-1, comment-4])
│       "Also, don't forget about query planning!"
│
└── Comment 5 (thread_id: comment-5, path: [comment-5])
    "I disagree with the first approach..."
    │
    └── Comment 6 (thread_id: comment-5, path: [comment-5, comment-6])
        "Can you explain why?"
```

### Database Storage Example

| id | content | post_id | parent_id | path | thread_id | depth |
|----|---------|---------|-----------|------|-----------|-------|
| comment-1 | "Great question!..." | post-1 | NULL | [comment-1] | comment-1 | 1 |
| comment-2 | "Thanks for tips!..." | post-1 | comment-1 | [comment-1, comment-2] | comment-1 | 2 |
| comment-3 | "For indexing..." | post-1 | comment-2 | [comment-1, comment-2, comment-3] | comment-1 | 3 |
| comment-4 | "Don't forget..." | post-1 | comment-1 | [comment-1, comment-4] | comment-1 | 2 |
| comment-5 | "I disagree..." | post-1 | NULL | [comment-5] | comment-5 | 1 |
| comment-6 | "Can you explain..." | post-1 | comment-5 | [comment-5, comment-6] | comment-5 | 2 |

---

## ⚡ Performance Optimizations

### Database Level

#### 1. **Connection Pooling**
```go
// config/database.go
db.SetMaxOpenConns(25)    // Maximum open connections
db.SetMaxIdleConns(5)     // Maximum idle connections
db.SetConnMaxLifetime(5 * time.Minute)
```

#### 2. **Query Optimization**
- **Prepared Statements**: Reduce parsing overhead
- **Batch Operations**: Group multiple inserts
- **Selective Fields**: Only fetch required columns
- **Pagination**: Limit result sets

#### 3. **Index Strategy**
```sql
-- Composite indexes for common query patterns
CREATE INDEX idx_comments_post_thread_created 
ON comments(post_id, thread_id, created_at DESC);

-- Partial indexes for active comments
CREATE INDEX idx_active_comments 
ON comments(post_id, created_at) 
WHERE deleted_at IS NULL;
```

### Application Level

#### 1. **Caching Strategy**
```go
// Redis caching for frequently accessed data
type CacheService struct {
    redis *redis.Client
    ttl   time.Duration
}

// Cache comment trees for popular posts
func (c *CacheService) GetCommentTree(postID string) (*CommentTree, error) {
    key := fmt.Sprintf("comments:post:%s", postID)
    // Check cache first, fallback to database
}
```

#### 2. **Lazy Loading**
```go
// Load comments on demand
type CommentResponse struct {
    ID       string `json:"id"`
    Content  string `json:"content"`
    Author   User   `json:"author"`
    Replies  []CommentResponse `json:"replies,omitempty"`
    HasMore  bool   `json:"has_more"`
}
```

#### 3. **Response Optimization**
- **JSON Streaming**: For large comment trees
- **Compression**: GZIP compression for responses
- **Field Selection**: Allow clients to specify required fields

### Memory Management

#### 1. **Efficient Data Structures**
```go
// Use sync.Pool for frequently allocated objects
var commentPool = sync.Pool{
    New: func() interface{} {
        return &Comment{}
    },
}

func GetComment() *Comment {
    return commentPool.Get().(*Comment)
}

func PutComment(c *Comment) {
    c.Reset() // Clear fields
    commentPool.Put(c)
}
```

#### 2. **Streaming Responses**
```go
// Stream large comment trees
func (h *CommentHandler) StreamComments(c *gin.Context) {
    c.Header("Content-Type", "application/json")
    c.Header("Transfer-Encoding", "chunked")
    
    encoder := json.NewEncoder(c.Writer)
    // Stream comments as they're processed
}
```

---

## 🔒 Security Architecture

### Authentication & Authorization

#### 1. **JWT Token Strategy**
```go
type TokenClaims struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.StandardClaims
}

// Dual token system
type AuthResponse struct {
    AccessToken  string `json:"access_token"`  // 15 minutes
    RefreshToken string `json:"refresh_token"` // 7 days
    ExpiresAt    int64  `json:"expires_at"`
}
```

#### 2. **Permission System**
```go
// Resource-based permissions
type Permission struct {
    Resource string // "post", "comment", "user"
    Action   string // "create", "read", "update", "delete"
    Owner    bool   // Can only access own resources
}

func (m *AuthMiddleware) CheckPermission(resource, action string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("user_id")
        resourceID := c.Param("id")
        
        if !m.hasPermission(userID, resource, action, resourceID) {
            c.JSON(403, gin.H{"error": "Forbidden"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### Input Validation & Sanitization

#### 1. **Request Validation**
```go
type CreatePostRequest struct {
    Title   string `json:"title" validate:"required,min=1,max=200"`
    Content string `json:"content" validate:"required,min=1,max=10000"`
}

func (v *Validator) ValidateCreatePost(req *CreatePostRequest) error {
    return v.validate.Struct(req)
}
```

#### 2. **HTML Sanitization**
```go
// Prevent XSS attacks
func SanitizeHTML(content string) string {
    p := bluemonday.UGCPolicy()
    p.AllowElements("p", "br", "strong", "em", "u", "a", "ul", "ol", "li")
    p.AllowAttrs("href").OnElements("a")
    return p.Sanitize(content)
}
```

### Rate Limiting

```go
// IP-based rate limiting
type RateLimiter struct {
    requests map[string][]time.Time
    mutex    sync.RWMutex
    limit    int
    window   time.Duration
}

func (rl *RateLimiter) Allow(ip string) bool {
    rl.mutex.Lock()
    defer rl.mutex.Unlock()
    
    now := time.Now()
    requests := rl.requests[ip]
    
    // Remove old requests outside window
    var validRequests []time.Time
    for _, req := range requests {
        if now.Sub(req) < rl.window {
            validRequests = append(validRequests, req)
        }
    }
    
    if len(validRequests) >= rl.limit {
        return false
    }
    
    validRequests = append(validRequests, now)
    rl.requests[ip] = validRequests
    return true
}
```

---

## 🌐 API Design

### RESTful Principles

#### 1. **Resource-Based URLs**
```
GET    /api/v1/posts                    # List posts
POST   /api/v1/posts                    # Create post
GET    /api/v1/posts/{id}               # Get specific post
PUT    /api/v1/posts/{id}               # Update post
DELETE /api/v1/posts/{id}               # Delete post

GET    /api/v1/posts/{id}/comments      # Get post comments
POST   /api/v1/posts/{id}/comments      # Create comment on post
```

#### 2. **Consistent Response Format**
```go
type APIResponse struct {
    StatusCode   int         `json:"status_code"`
    ErrorMessage *string     `json:"error_message"`
    Data         interface{} `json:"data"`
}

type PaginatedResponse struct {
    Items      interface{} `json:"items"`
    Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
    Limit   int  `json:"limit"`
    Offset  int  `json:"offset"`
    Total   int  `json:"total"`
    HasMore bool `json:"has_more"`
}
```

#### 3. **Error Handling**
```go
type APIError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

// Standard error responses
var (
    ErrValidation    = APIError{Code: "VALIDATION_ERROR", Message: "Invalid input"}
    ErrUnauthorized  = APIError{Code: "UNAUTHORIZED", Message: "Authentication required"}
    ErrForbidden     = APIError{Code: "FORBIDDEN", Message: "Access denied"}
    ErrNotFound      = APIError{Code: "NOT_FOUND", Message: "Resource not found"}
    ErrConflict      = APIError{Code: "CONFLICT", Message: "Resource already exists"}
)
```

### Content Negotiation

```go
func (h *Handler) HandleRequest(c *gin.Context) {
    accept := c.GetHeader("Accept")
    
    switch accept {
    case "application/json":
        h.respondJSON(c, data)
    case "application/xml":
        h.respondXML(c, data)
    default:
        h.respondJSON(c, data) // Default to JSON
    }
}
```

---

## 🚀 Deployment Architecture

### Container Strategy

#### 1. **Multi-Stage Dockerfile**
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

# Production stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

#### 2. **Docker Compose for Development**
```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis
    
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: post_comments_db
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    
  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
```

### Production Deployment

#### 1. **Kubernetes Deployment**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: post-comments-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: post-comments-service
  template:
    metadata:
      labels:
        app: post-comments-service
    spec:
      containers:
      - name: app
        image: post-comments-service:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: "postgres-service"
        - name: JWT_SECRET_KEY
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: jwt-secret
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

#### 2. **Load Balancer Configuration**
```nginx
upstream backend {
    least_conn;
    server app1:8080 max_fails=3 fail_timeout=30s;
    server app2:8080 max_fails=3 fail_timeout=30s;
    server app3:8080 max_fails=3 fail_timeout=30s;
}

server {
    listen 80;
    server_name api.example.com;
    
    location / {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Rate limiting
        limit_req zone=api burst=20 nodelay;
        
        # Timeouts
        proxy_connect_timeout 5s;
        proxy_send_timeout 10s;
        proxy_read_timeout 10s;
    }
}
```

---

## 📊 Monitoring & Observability

### Health Checks

```go
type HealthChecker struct {
    db    *sql.DB
    redis *redis.Client
}

func (h *HealthChecker) Check() HealthStatus {
    status := HealthStatus{
        Status:    "healthy",
        Timestamp: time.Now(),
        Version:   version,
        Checks:    make(map[string]string),
    }
    
    // Database health
    if err := h.db.Ping(); err != nil {
        status.Status = "unhealthy"
        status.Checks["database"] = "disconnected"
    } else {
        status.Checks["database"] = "connected"
    }
    
    // Redis health
    if err := h.redis.Ping().Err(); err != nil {
        status.Checks["redis"] = "disconnected"
    } else {
        status.Checks["redis"] = "connected"
    }
    
    return status
}
```

### Metrics Collection

```go
// Prometheus metrics
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "endpoint"},
    )
    
    dbConnectionsActive = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "db_connections_active",
            Help: "Number of active database connections",
        },
    )
)
```

### Logging Strategy

```go
type Logger struct {
    *logrus.Logger
}

func (l *Logger) LogRequest(c *gin.Context, duration time.Duration) {
    l.WithFields(logrus.Fields{
        "method":     c.Request.Method,
        "path":       c.Request.URL.Path,
        "status":     c.Writer.Status(),
        "duration":   duration,
        "ip":         c.ClientIP(),
        "user_agent": c.Request.UserAgent(),
        "user_id":    c.GetString("user_id"),
    }).Info("HTTP request")
}
```

---

## 📈 Scalability Considerations

### Horizontal Scaling

#### 1. **Stateless Design**
- No server-side sessions
- JWT tokens for authentication
- Database for persistent state
- Redis for shared cache

#### 2. **Database Scaling**
```sql
-- Read replicas for query distribution
-- Master-slave configuration
-- Connection pooling
-- Query optimization

-- Partitioning strategy for large tables
CREATE TABLE comments_2024_01 PARTITION OF comments
FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');
```

#### 3. **Caching Strategy**
```go
// Multi-level caching
type CacheManager struct {
    l1 *sync.Map        // In-memory cache
    l2 *redis.Client    // Redis cache
    l3 Database         // Database
}

func (cm *CacheManager) Get(key string) (interface{}, error) {
    // Try L1 cache first
    if value, ok := cm.l1.Load(key); ok {
        return value, nil
    }
    
    // Try L2 cache
    if value, err := cm.l2.Get(key).Result(); err == nil {
        cm.l1.Store(key, value) // Populate L1
        return value, nil
    }
    
    // Fallback to database
    value, err := cm.l3.Get(key)
    if err == nil {
        cm.l2.Set(key, value, time.Hour) // Populate L2
        cm.l1.Store(key, value)          // Populate L1
    }
    
    return value, err
}
```

### Performance Monitoring

#### 1. **Key Metrics**
- **Response Time**: P50, P95, P99 latencies
- **Throughput**: Requests per second
- **Error Rate**: 4xx and 5xx responses
- **Database Performance**: Query execution time
- **Memory Usage**: Heap size, GC frequency
- **Connection Pools**: Active/idle connections

#### 2. **Alerting Rules**
```yaml
# Prometheus alerting rules
groups:
- name: post-comments-service
  rules:
  - alert: HighErrorRate
    expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
    for: 2m
    annotations:
      summary: "High error rate detected"
      
  - alert: HighLatency
    expr: histogram_quantile(0.95, http_request_duration_seconds) > 1
    for: 5m
    annotations:
      summary: "High latency detected"
      
  - alert: DatabaseConnectionsHigh
    expr: db_connections_active > 20
    for: 1m
    annotations:
      summary: "Database connection pool nearly exhausted"
```

---

## 🔄 Future Enhancements

### Planned Features

1. **Real-time Updates**
   - WebSocket support for live comments
   - Server-Sent Events for notifications
   - Redis pub/sub for real-time messaging

2. **Advanced Comment Features**
   - Comment reactions (likes, dislikes)
   - Comment moderation system
   - Spam detection and filtering
   - Rich text formatting support

3. **Performance Optimizations**
   - GraphQL API for flexible queries
   - Database sharding for massive scale
   - CDN integration for static content
   - Advanced caching strategies

4. **Analytics & Insights**
   - Comment engagement metrics
   - User behavior analytics
   - Content popularity tracking
   - Performance dashboards

### Migration Strategies

#### 1. **Database Migrations**
```sql
-- Example: Adding comment reactions
CREATE TABLE comment_reactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    comment_id UUID NOT NULL REFERENCES comments(id),
    user_id UUID NOT NULL REFERENCES users(id),
    reaction_type VARCHAR(20) NOT NULL, -- 'like', 'dislike', 'love', etc.
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(comment_id, user_id, reaction_type)
);

CREATE INDEX idx_comment_reactions_comment_id ON comment_reactions(comment_id);
CREATE INDEX idx_comment_reactions_user_id ON comment_reactions(user_id);
```

#### 2. **API Versioning**
```go
// Support multiple API versions
func SetupRoutes(r *gin.Engine) {
    v1 := r.Group("/api/v1")
    {
        v1.GET("/posts", handlers.V1.GetPosts)
        v1.POST("/posts", handlers.V1.CreatePost)
    }
    
    v2 := r.Group("/api/v2")
    {
        v2.GET("/posts", handlers.V2.GetPosts) // Enhanced with reactions
        v2.POST("/posts", handlers.V2.CreatePost)
    }
}
```

---

## 📚 References & Best Practices

### Database Design References
- [Storing Hierarchical Data in a Database](https://www.sitepoint.com/hierarchical-data-database/)
- [PostgreSQL Array Performance](https://www.postgresql.org/docs/current/arrays.html)
- [Materialized Path Pattern](https://docs.mongodb.com/manual/tutorial/model-tree-structures-with-materialized-paths/)

### Go Best Practices
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Clean Architecture in Go](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

### Security Guidelines
- [OWASP API Security Top 10](https://owasp.org/www-project-api-security/)
- [JWT Best Practices](https://auth0.com/blog/a-look-at-the-latest-draft-for-jwt-bcp/)
- [Go Security Checklist](https://github.com/Checkmarx/Go-SCP)

---

*This architecture documentation is a living document and should be updated as the system evolves.* 