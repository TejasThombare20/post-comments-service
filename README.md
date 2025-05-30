# Post-Comments Backend Service

A modular, production-grade backend service built with **Go (Golang)**, **Gin framework**, and **PostgreSQL** for managing users, posts, and nested comments with comprehensive development tools and containerization support.

## 📚 Documentation

- **[API Documentation](API.md)** - Comprehensive API reference with examples
- **[Architecture Guide](ARCHITECTURE.md)** - System architecture, database design, and nested comments implementation
- **[Setup Instructions](#setup-instructions)** - Getting started with Docker and local development

## 🚀 Features

- **User Management**: Create, read, update, and delete users with authentication
- **Post Management**: Create, read, update, and delete posts with author information
- **Nested Comments**: Support for threaded comments with unlimited nesting
- **JWT Authentication**: Secure authentication with access and refresh tokens
- **Clean Architecture**: Modular design with separation of concerns
- **Professional Logging**: Structured logging with middleware
- **Error Handling**: Comprehensive error handling with custom error types
- **Database Migrations**: SQL migrations with indexes for performance
- **CORS Support**: Cross-origin resource sharing enabled
- **Input Validation**: Request validation using custom validators
- **Pagination**: Built-in pagination support for list endpoints
- **HTML Sanitization**: Safe HTML content processing for comments
- **Configuration Validation**: Startup validation for all environment variables
- **Development Tools**: Pre-commit hooks, linting, formatting, and testing tools
- **Containerization**: Docker support for easy deployment and development

## 📁 Project Structure

```
post-comments-service/
├── cmd/
│   └── main.go                    # Application entrypoint
├── config/
│   └── config.go                  # Configuration management with validation
├── controllers/
│   ├── auth_controller.go         # Authentication HTTP handlers
│   ├── user_controller.go         # User HTTP handlers
│   ├── post_controller.go         # Post HTTP handlers
│   └── comment_controller.go      # Comment HTTP handlers
├── models/
│   ├── user.go                    # User model and DTOs
│   ├── post.go                    # Post model and DTOs
│   ├── comment.go                 # Comment model and DTOs
│   └── auth.go                    # Authentication models
├── services/
│   ├── auth_service.go            # Authentication business logic
│   ├── user_service.go            # User business logic
│   ├── post_service.go            # Post business logic
│   ├── comment_service.go         # Comment business logic
│   └── jwt_service.go             # JWT token management
├── repository/
│   ├── user_repo.go               # User data access layer
│   ├── post_repo.go               # Post data access layer
│   └── comment_repo.go            # Comment data access layer
├── routes/
│   └── routes.go                  # Route definitions
├── middleware/
│   ├── logger.go                  # Logging middleware
│   ├── auth.go                    # Authentication middleware
│   └── cors.go                    # CORS middleware
├── utils/
│   ├── response.go                # Standard API response helpers
│   ├── error_handler.go           # Error handling utilities
│   ├── auth_helpers.go            # Authentication utilities
│   └── html_sanitizer.go          # HTML content sanitization
├── validator/
│   └── validator.go               # Custom validation logic
├── migrations/
│   └── 001_init.sql               # Database schema migration
├── .pre-commit-config.yaml        # Pre-commit hooks configuration
├── .golangci.yml                  # Go linting configuration
├── Dockerfile                     # Docker container configuration
├── docker-compose.yml             # Docker Compose for development
├── Makefile                       # Development automation tasks
├── go.mod                         # Go module definition
├── go.sum                         # Go module checksums
├── .env.example                   # Environment variables example
├── API.md                         # Comprehensive API documentation
└── README.md                      # This file
```

## 🛠️ Setup Instructions

### Prerequisites

- **Go 1.21+** (for non-Docker setup)
- **Docker & Docker Compose** (for Docker setup)
- **PostgreSQL 12+** (for non-Docker setup)
- **Git**

### Option 1: Docker Setup (Recommended)

This is the easiest way to get started with the project.

#### 1. Clone the Repository

```bash
git clone <repository-url>
cd post-comments-service
```

#### 2. Environment Configuration

```bash
cp .env.example .env
# Edit .env file with your preferred settings
```

#### 3. Start with Docker Compose

```bash
# Start all services (database, app, redis, adminer)
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

#### 4. Access the Application

- **API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **Database Admin (Adminer)**: http://localhost:8081
- **Redis**: localhost:6379

### Option 2: Local Development Setup

#### 1. Clone and Install Dependencies

```bash
git clone <repository-url>
cd post-comments-service

# Install Go dependencies
go mod tidy

# Install development tools
make install-tools
```

#### 2. Database Setup

```bash
# Create PostgreSQL database
createdb post_comments_db

# Run migrations
make migrate-up
```

#### 3. Environment Configuration

```bash
cp .env.example .env
```

Edit `.env` with your database credentials:
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=post_comments_db
DB_SSLMODE=disable

# Server Configuration
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info
DEBUG=true

# JWT Configuration (REQUIRED)
JWT_SECRET_KEY=your-super-secret-jwt-key-that-is-at-least-32-characters-long
JWT_ACCESS_TOKEN_DURATION=15m
JWT_REFRESH_TOKEN_DURATION=168h

# Database Pool Configuration
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m

# Server Timeouts
SERVER_READ_TIMEOUT=15s
SERVER_WRITE_TIMEOUT=15s
SERVER_IDLE_TIMEOUT=60s
```

#### 4. Setup Development Tools

```bash
# Setup pre-commit hooks
make setup-hooks

# Run code formatting and linting
make quality
```

#### 5. Run the Application

```bash
# Development mode with hot reload
make dev

# Or regular run
make run
```

## 🔧 Development Tools

### Available Make Commands

```bash
# Code Quality
make fmt           # Format code with gofmt
make imports       # Organize imports with goimports
make lint          # Run golangci-lint
make vet           # Run go vet
make quality       # Run all quality checks

# Testing
make test          # Run tests
make coverage      # Run tests with coverage report

# Building
make build         # Build binary
make build-linux   # Build for Linux
make clean         # Clean build artifacts

# Development
make run           # Run application
make dev           # Run with hot reload (requires air)
make deps          # Install dependencies
make tidy          # Tidy go modules

# Docker
make docker-build  # Build Docker image
make docker-run    # Run Docker container

# Database
make migrate-up    # Run database migrations

# Security
make security      # Run security scan (requires gosec)

# Setup
make install-tools # Install development tools
make setup-hooks   # Setup pre-commit hooks

# CI/CD
make ci            # Simulate CI pipeline
make help          # Show all available commands
```

### Pre-commit Hooks

The project includes pre-commit hooks that run automatically before each commit:

- **Go formatting** (gofmt)
- **Import organization** (goimports)
- **Go vet** checks
- **Go mod tidy**
- **Linting** (golangci-lint)
- **Security scanning** (gosec)
- **YAML/JSON validation**
- **Trailing whitespace removal**

Install hooks:
```bash
make setup-hooks
```

### Code Quality Tools

- **golangci-lint**: Comprehensive Go linting
- **gofmt**: Code formatting
- **goimports**: Import organization
- **gosec**: Security vulnerability scanning
- **go vet**: Static analysis

## 📚 API Documentation

For comprehensive API documentation including all endpoints, request/response formats, authentication flows, and examples, see **[API.md](API.md)**.

### Quick API Overview

- **Base URL**: `http://localhost:8080/api/v1`
- **Authentication**: JWT Bearer tokens
- **Response Format**: Standardized JSON responses
- **Health Check**: `GET /health`

### Key Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/auth/register` | Register new user | No |
| POST | `/auth/login` | User login | No |
| GET | `/auth/profile` | Get user profile | Yes |
| GET | `/users` | List users | No |
| POST | `/posts` | Create post | Yes |
| GET | `/posts` | List posts | No |
| POST | `/posts/post-comments/{id}` | Create comment | Yes |
| GET | `/posts/post-comments/{id}` | Get comments | No |

For detailed documentation with request/response examples, authentication flows, error handling, and testing instructions, see **[API.md](API.md)**.

## 🔒 Environment Variables

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `JWT_SECRET_KEY` | JWT signing key (min 32 chars) | `your-super-secret-jwt-key-that-is-at-least-32-characters-long` |
| `DB_HOST` | Database host | `localhost` |
| `DB_USER` | Database user | `postgres` |
| `DB_NAME` | Database name | `post_comments_db` |

### Optional Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `DB_PORT` | `5432` | Database port |
| `DB_PASSWORD` | `` | Database password |
| `DB_SSLMODE` | `disable` | Database SSL mode |
| `ENVIRONMENT` | `development` | Environment (development/staging/production) |
| `LOG_LEVEL` | `info` | Log level (debug/info/warn/error) |
| `DEBUG` | `false` | Debug mode |
| `JWT_ACCESS_TOKEN_DURATION` | `15m` | Access token duration |
| `JWT_REFRESH_TOKEN_DURATION` | `168h` | Refresh token duration |

For a complete list of environment variables, see `.env.example`.

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run specific package tests
go test ./services/...
```

## 🚀 Deployment

### Docker Deployment

```bash
# Build production image
docker build -t post-comments-service:latest .

# Run with environment file
docker run -p 8080:8080 --env-file .env post-comments-service:latest
```

### Binary Deployment

```bash
# Build for Linux
make build-linux

# Copy binary and run on server
./post-comments-service_unix
```

## 🔧 Configuration Validation

The application validates all configuration on startup and will fail to start if:

- Required environment variables are missing
- JWT secret key is less than 32 characters
- Database connection parameters are invalid
- Port numbers are out of range
- Duration values are invalid

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

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run quality checks (`make quality`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Development Workflow

1. **Setup**: `make install-tools && make setup-hooks`
2. **Code**: Make your changes
3. **Quality**: `make quality` (runs formatting, linting, vetting)
4. **Test**: `make test`
5. **Build**: `make build`
6. **Commit**: Git will automatically run pre-commit hooks

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Troubleshooting

### Common Issues

1. **Config validation errors**: Check your `.env` file against `.env.example`
2. **Database connection issues**: Ensure PostgreSQL is running and credentials are correct
3. **Port already in use**: Change the `PORT` environment variable
4. **JWT errors**: Ensure `JWT_SECRET_KEY` is at least 32 characters long

### Docker Issues

```bash
# Rebuild containers
docker-compose down && docker-compose up --build

# View logs
docker-compose logs -f app

# Reset database
docker-compose down -v && docker-compose up
```

### Development Issues

```bash
# Clean and rebuild
make clean && make build

# Update dependencies
go mod tidy

# Reset pre-commit hooks
make setup-hooks
```

## 📞 Support

For support, please open an issue in the GitHub repository or contact the development team.

## 👥 Authors

- **Tejas Thombare** - 

## 🙏 Acknowledgments

- Gin Web Framework
- GORM ORM
- PostgreSQL
- Go community 
