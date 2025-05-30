# Post-Comments Backend Service

A modular, production-grade backend service built with **Go (Golang)**, **Gin framework**, and **PostgreSQL** for managing users, posts, and nested comments.

## 🚀 Features

- **User Management**: Create, read, update, and delete users
- **Post Management**: Create, read, update, and delete posts
- **Nested Comments**: Support for threaded comments with unlimited nesting
- **Clean Architecture**: Modular design with separation of concerns
- **Professional Logging**: Structured logging with middleware
- **Error Handling**: Comprehensive error handling with custom error types
- **Database Migrations**: SQL migrations with indexes for performance
- **CORS Support**: Cross-origin resource sharing enabled
- **Input Validation**: Request validation using Gin's binding
- **Pagination**: Built-in pagination support for list endpoints
- **JWT Authentication**: Secure authentication with access and refresh tokens
- **Authorization**: Role-based access control for protected endpoints

## 📁 Project Structure

```
post-comments-service/
├── cmd/
│   └── main.go                 # Application entrypoint
├── config/
│   └── config.go              # Database and environment configuration
├── controllers/
│   ├── user_controller.go     # User HTTP handlers
│   └── post_controller.go     # Post HTTP handlers
├── models/
│   ├── user.go               # User model and DTOs
│   ├── post.go               # Post model and DTOs
│   └── comment.go            # Comment model and DTOs
├── services/
│   ├── user_service.go       # User business logic
│   └── post_service.go       # Post business logic
├── repository/
│   ├── user_repo.go          # User data access layer
│   └── post_repo.go          # Post data access layer
├── routes/
│   └── routes.go             # Route definitions
├── middleware/
│   ├── logger.go             # Logging middleware
│   └── auth.go               # Authentication middleware
├── utils/
│   ├── response.go           # Standard API response helpers
│   └── error_handler.go      # Error handling utilities
├── migrations/
│   └── 001_init.sql          # Database schema migration
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
├── .env.example              # Environment variables example
└── README.md                 # This file
```

## 🛠️ Setup Instructions

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

### 1. Clone the Repository

```bash
git clone <repository-url>
cd post-comments-service
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Database Setup

1. Create a PostgreSQL database:
```sql
CREATE DATABASE post_comments_db;
```

2. Run the migration:
```bash
psql -U postgres -d post_comments_db -f migrations/001_init.sql
```

### 4. Environment Configuration

1. Copy the example environment file:
```bash
cp .env.example .env
```

2. Update the `.env` file with your database credentials:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=post_comments_db
DB_SSLMODE=disable
PORT=8080
```

### 5. Run the Application

```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080`

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Response Format
All API responses follow this standard format:
```json
{
  "status_code": 200,
  "error_message": null,
  "data": { ... }
}
```

### User Endpoints

#### Create User
```http
POST /api/v1/users
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "securepassword",
  "display_name": "John Doe",
  "avatar_url": "https://example.com/avatar.jpg"
}
```

#### Get User by ID
```http
GET /api/v1/users/{id}
```

#### Get User by Username
```http
GET /api/v1/users/username/{username}
```

#### List Users
```http
GET /api/v1/users?limit=10&offset=0
```

#### Update User
```http
PUT /api/v1/users/{id}
Content-Type: application/json

{
  "email": "newemail@example.com",
  "display_name": "New Display Name"
}
```

#### Delete User
```http
DELETE /api/v1/users/{id}
```

### Post Endpoints

#### Create Post
```http
POST /api/v1/posts
Content-Type: application/json
X-User-ID: {user_uuid}

{
  "title": "My First Post",
  "content": "This is the content of my first post."
}
```

#### Get Post by ID
```http
GET /api/v1/posts/{id}
```

#### Get Post with Comments
```http
GET /api/v1/posts/{id}/comments
```

#### List Posts
```http
GET /api/v1/posts?limit=10&offset=0
```

#### List Posts by User
```http
GET /api/v1/users/{userId}/posts?limit=10&offset=0
```

#### Update Post
```http
PUT /api/v1/posts/{id}
Content-Type: application/json
X-User-ID: {user_uuid}

{
  "title": "Updated Title",
  "content": "Updated content"
}
```

#### Delete Post
```http
DELETE /api/v1/posts/{id}
X-User-ID: {user_uuid}
```

### Health Check
```http
GET /health
```

## 🗄️ Database Schema

### Users Table
- `id` (UUID, Primary Key)
- `username` (TEXT, Unique, Not Null)
- `email` (TEXT, Unique)
- `password_hash` (TEXT)
- `display_name` (TEXT)
- `avatar_url` (TEXT)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)
- `deleted_at` (TIMESTAMP)

### Posts Table
- `id` (UUID, Primary Key)
- `title` (TEXT, Not Null)
- `content` (TEXT, Not Null)
- `created_by` (UUID, Foreign Key to users.id)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)
- `deleted_at` (TIMESTAMP)

### Comments Table
- `id` (UUID, Primary Key)
- `content` (TEXT, Not Null)
- `post_id` (UUID, Foreign Key to posts.id)
- `parent_id` (UUID, Foreign Key to comments.id)
- `path` (UUID[], Array for nested structure)
- `thread_id` (UUID, Root comment ID)
- `created_by` (UUID, Foreign Key to users.id)
- `created_at` (TIMESTAMP)
- `deleted_at` (TIMESTAMP)

## 🔧 Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o bin/post-comments-service cmd/main.go
```

### Docker Support (Future Enhancement)
```dockerfile
# Dockerfile example
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

## 🚧 Future Enhancements

- [ ] JWT Authentication implementation
- [ ] Comment CRUD operations
- [ ] File upload for avatars
- [ ] Real-time notifications
- [ ] Rate limiting
- [ ] Caching with Redis
- [ ] Full-text search
- [ ] API versioning
- [ ] Swagger documentation
- [ ] Docker containerization
- [ ] Unit and integration tests
- [ ] CI/CD pipeline

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **Tejas Thombare** - *Initial work* - [TejasThombare20](https://github.com/TejasThombare20)

## 🙏 Acknowledgments

- Gin Web Framework
- GORM ORM
- PostgreSQL
- Go community 

## Authentication System

### JWT-Based Authentication

The service implements a robust JWT-based authentication system with the following features:

#### Token Types
- **Access Token**: Short-lived token (15 minutes) for API access
- **Refresh Token**: Long-lived token (7 days) for obtaining new access tokens

#### Authentication Endpoints

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| POST | `/api/v1/auth/register` | Register a new user | Public |
| POST | `/api/v1/auth/login` | Login user | Public |
| POST | `/api/v1/auth/refresh` | Refresh access token | Public |
| POST | `/api/v1/auth/logout` | Logout user | Public |
| GET | `/api/v1/auth/profile` | Get user profile | Required |
| POST | `/api/v1/auth/change-password` | Change password | Required |

#### Registration Request
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securepassword123",
  "display_name": "John Doe",
  "avatar_url": "https://example.com/avatar.jpg"
}
```

#### Login Request
```json
{
  "username": "johndoe",
  "password": "securepassword123"
}
```

#### Authentication Response
```json
{
  "status_code": 200,
  "error_message": null,
  "data": {
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "username": "johndoe",
      "email": "john@example.com",
      "display_name": "John Doe",
      "avatar_url": "https://example.com/avatar.jpg",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-01T00:15:00Z"
  }
}
```

#### Using Authentication

Include the access token in the Authorization header:
```
Authorization: Bearer <access_token>
```

### Protected Endpoints

The following endpoints require authentication:

#### User Management
- `POST /api/v1/users` - Create user
- `PUT /api/v1/users/user/:id` - Update user
- `DELETE /api/v1/users/user/:id` - Delete user

#### Post Management
- `POST /api/v1/posts` - Create post
- `PUT /api/v1/posts/post/:id` - Update post
- `DELETE /api/v1/posts/post/:id` - Delete post

#### Comment Management
- `POST /api/v1/comments` - Create comment
- `PUT /api/v1/comments/:id` - Update comment
- `DELETE /api/v1/comments/:id` - Delete comment

### Public Endpoints

The following endpoints are publicly accessible:

#### User Information (Read-only)
- `GET /api/v1/users` - List users
- `GET /api/v1/users/username/:username` - Get user by username
- `GET /api/v1/users/user/:id` - Get user by ID
- `GET /api/v1/users/:userId/posts` - Get user's posts

#### Post Information (Read-only)
- `GET /api/v1/posts` - List posts
- `GET /api/v1/posts/post/:id` - Get post by ID
- `GET /api/v1/posts/post/:id/comments` - Get post with comments

## API Endpoints

### Authentication

| Method | Endpoint | Description | Body | Response |
|--------|----------|-------------|------|----------|
| POST | `/api/v1/auth/register` | Register new user | RegisterRequest | AuthResponse |
| POST | `/api/v1/auth/login` | Login user | LoginRequest | AuthResponse |
| POST | `/api/v1/auth/refresh` | Refresh token | RefreshTokenRequest | AuthResponse |
| POST | `/api/v1/auth/logout` | Logout user | - | Success |
| GET | `/api/v1/auth/profile` | Get profile | - | UserResponse |
| POST | `/api/v1/auth/change-password` | Change password | ChangePasswordRequest | Success |

### Users

| Method | Endpoint | Description | Authentication | Body | Response |
|--------|----------|-------------|----------------|------|----------|
| POST | `/api/v1/users` | Create user | Required | CreateUserRequest | UserResponse |
| GET | `/api/v1/users` | List users | Optional | - | []UserResponse |
| GET | `/api/v1/users/user/:id` | Get user by ID | Optional | - | UserResponse |
| GET | `/api/v1/users/username/:username` | Get user by username | Optional | - | UserResponse |
| PUT | `/api/v1/users/user/:id` | Update user | Required | UpdateUserRequest | UserResponse |
| DELETE | `/api/v1/users/user/:id` | Delete user | Required | - | Success |
| GET | `/api/v1/users/:userId/posts` | Get user's posts | Optional | - | []PostResponse |

### Posts

| Method | Endpoint | Description | Authentication | Body | Response |
|--------|----------|-------------|----------------|------|----------|
| POST | `/api/v1/posts` | Create post | Required | CreatePostRequest | PostResponse |
| GET | `/api/v1/posts` | List posts | Optional | - | []PostResponse |
| GET | `/api/v1/posts/post/:id` | Get post by ID | Optional | - | PostResponse |
| PUT | `/api/v1/posts/post/:id` | Update post | Required | UpdatePostRequest | PostResponse |
| DELETE | `/api/v1/posts/post/:id` | Delete post | Required | - | Success |
| GET | `/api/v1/posts/post/:id/comments` | Get post with comments | Optional | - | PostWithCommentsResponse |

### Comments (Future Implementation)

| Method | Endpoint | Description | Authentication | Body | Response |
|--------|----------|-------------|----------------|------|----------|
| POST | `/api/v1/comments` | Create comment | Required | CreateCommentRequest | CommentResponse |
| GET | `/api/v1/comments/:id` | Get comment by ID | Optional | - | CommentResponse |
| PUT | `/api/v1/comments/:id` | Update comment | Required | UpdateCommentRequest | CommentResponse |
| DELETE | `/api/v1/comments/:id` | Delete comment | Required | - | Success |

## Environment Variables

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=post_comments_db
DB_SSLMODE=disable

# Server Configuration
PORT=8080

# Environment
ENV=development

# JWT Configuration
JWT_SECRET_KEY=your-super-secret-jwt-key-change-this-in-production-min-32-chars
JWT_ACCESS_TOKEN_TTL=15m
JWT_REFRESH_TOKEN_TTL=168h

# CORS Configuration
CORS_ALLOWED_ORIGINS=*
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Origin,Content-Type,Accept,Authorization,X-User-ID
```

## Getting Started

### Prerequisites

- Go 1.23 or higher
- PostgreSQL 12 or higher
- Git

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd post-comments-service
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Set up the database:
```bash
# Create database
createdb post_comments_db

# Run migrations
go run cmd/migrate.go
```

5. Build and run the application:
```bash
go build -o main cmd/main.go
./main
```

The server will start on `http://localhost:8080`

### Docker Setup

1. Start the services:
```bash
docker-compose up -d
```

2. The API will be available at `http://localhost:8080`

## Usage Examples

### Authentication Flow

1. **Register a new user:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword123",
    "display_name": "John Doe"
  }'
```

2. **Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "securepassword123"
  }'
```

3. **Use the access token for protected endpoints:**
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my first post."
  }'
```

4. **Refresh token when access token expires:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "<refresh_token>"
  }'
```

### Creating a Post

```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "title": "My Post Title",
    "content": "This is the content of my post.",
    "tags": ["technology", "programming"]
  }'
```

### Getting Posts

```bash
# Get all posts (public)
curl http://localhost:8080/api/v1/posts

# Get specific post (public)
curl http://localhost:8080/api/v1/posts/post/{post_id}

# Get posts by user (public)
curl http://localhost:8080/api/v1/users/{user_id}/posts
```

## Security Features

- **Password Hashing**: Passwords are hashed using bcrypt
- **JWT Tokens**: Secure token-based authentication
- **Token Expiration**: Access tokens expire in 15 minutes
- **Refresh Tokens**: Long-lived tokens for seamless user experience
- **CORS Protection**: Configurable CORS settings
- **Input Validation**: Comprehensive request validation
- **SQL Injection Protection**: Parameterized queries
- **Authorization**: User can only modify their own resources

## Architecture

The application follows a clean architecture pattern:

```
├── cmd/                    # Application entry points
├── config/                 # Configuration management
├── controllers/            # HTTP handlers
├── middleware/             # HTTP middleware
├── models/                 # Data models and DTOs
├── repository/             # Data access layer
├── routes/                 # Route definitions
├── services/               # Business logic layer
├── utils/                  # Utility functions
├── validator/              # Request validation
└── migrations/             # Database migrations
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 