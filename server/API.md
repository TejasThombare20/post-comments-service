# Post-Comments Service API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## Response Format

All API responses follow this standard format:
```json
{
  "status_code": 200,
  "error_message": null,
  "data": { ... }
}
```

### Error Response Format
```json
{
  "status_code": 400,
  "error_message": "Validation failed",
  "data": null
}
```

## Status Codes

| Code | Description |
|------|-------------|
| 200  | Success |
| 201  | Created |
| 400  | Bad Request |
| 401  | Unauthorized |
| 403  | Forbidden |
| 404  | Not Found |
| 409  | Conflict |
| 422  | Unprocessable Entity |
| 500  | Internal Server Error |

---

## Authentication Endpoints

### Register User
Create a new user account.

**Endpoint:** `POST /api/v1/auth/register`

**Request Body:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "securepassword",
  "display_name": "John Doe"
}
```

**Response:**
```json
{
  "status_code": 201,
  "error_message": null,
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "email": "john@example.com",
      "display_name": "John Doe",
      "avatar_url": null,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-15T10:45:00Z"
  }
}
```

### Login
Authenticate a user and receive tokens.

**Endpoint:** `POST /api/v1/auth/login`

**Request Body:**
```json
{
  "username": "john_doe",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "status_code": 200,
  "error_message": null,
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "email": "john@example.com",
      "display_name": "John Doe",
      "avatar_url": null
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-15T10:45:00Z"
  }
}
```

### Refresh Token
Get a new access token using a refresh token.

**Endpoint:** `POST /api/v1/auth/refresh`

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**
```json
{
  "status_code": 200,
  "error_message": null,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-15T11:00:00Z"
  }
}
```

### Get User Profile
Get the authenticated user's profile.

**Endpoint:** `GET /api/v1/auth/profile`

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "status_code": 200,
  "error_message": null,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "john_doe",
    "email": "john@example.com",
    "display_name": "John Doe",
    "avatar_url": null,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

---

## User Management Endpoints

### List Users
Get a paginated list of users.

**Endpoint:** `GET /api/v1/users`

**Query Parameters:**
- `limit` (optional): Number of users per page (default: 10, max: 100)
- `offset` (optional): Number of users to skip (default: 0)

**Example:** `GET /api/v1/users?limit=20&offset=0`

**Response:**
```json
{
  "status_code": 200,
  "error_message": null,
  "data": {
    "users": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "username": "john_doe",
        "email": "john@example.com",
        "display_name": "John Doe",
        "avatar_url": null,
        "created_at": "2024-01-15T10:30:00Z"
      }
    ],
    "pagination": {
      "limit": 20,
      "offset": 0,
      "total": 1,
      "has_more": false
    }
  }
}
```

### Get User by ID
Get a specific user by their ID.

**Endpoint:** `GET /api/v1/users/user/{id}`

**Path Parameters:**
- `id`: User UUID

**Response:**
```json
{
  "status_code": 200,
  "error_message": null,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "john_doe",
    "email": "john@example.com",
    "display_name": "John Doe",
    "avatar_url": null,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### Update User
Update user information (authenticated users can only update their own profile).

**Endpoint:** `PUT /api/v1/users/user/{id}`

**Headers:** `Authorization: Bearer <token>`

**Path Parameters:**
- `id`: User UUID

**Request Body:**
```json
{
  "email": "newemail@example.com",
  "display_name": "New Display Name",
  "avatar_url": "https://example.com/avatar.jpg"
}
```

**Response:**
```json
{
  "status_code": 200,
  "error_message": null,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "john_doe",
    "email": "newemail@example.com",
    "display_name": "New Display Name",
    "avatar_url": "https://example.com/avatar.jpg",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

### Delete User
Delete a user account (authenticated users can only delete their own account).

**Endpoint:** `DELETE /api/v1/users/user/{id}`

**Headers:** `Authorization: Bearer <token>`

**Path Parameters:**
- `id`: User UUID

**Response:**
```json
{
  "status_code": 200,
  "error_message": null,
  "data": {
    "message": "User deleted successfully"
  }
}
```

---

## Post Management Endpoints

### Create Post
Create a new post.

**Endpoint:** `POST /api/v1/posts`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "title": "My First Post",
  "content": "This is the content of my first post. It can contain HTML."
}
```

**Response:**
```json
{
  "status_code": 201,
  "error_message": null,
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "title": "My First Post",
    "content": "This is the content of my first post. It can contain HTML.",
    "created_by": "550e8400-e29b-41d4-a716-446655440000",
    "author": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "display_name": "John Doe"
    },
    "created_at": "2024-01-15T11:00:00Z",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

### Get Post by ID
Get a specific post by its ID.

**Endpoint:** `GET /api/v1/posts/post/{id}`

**Path Parameters:**
- `id`: Post UUID

**Response:**
```json
{
  "status_code": 200,
  "error_message": null,
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "title": "My First Post",
    "content": "This is the content of my first post.",
    "created_by": "550e8400-e29b-41d4-a716-446655440000",
    "author": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "display_name": "John Doe",
      "avatar_url": null
    },
    "created_at": "2024-01-15T11:00:00Z",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

### List Posts
Get a paginated list of posts.

**Endpoint:** `GET /api/v1/posts`

**Query Parameters:**
- `limit` (optional): Number of posts per page (default: 10, max: 100)
- `offset` (optional): Number of posts to skip (default: 0)
- `author_id` (optional): Filter posts by author ID

**Example:** `GET /api/v1/posts?limit=20&offset=0&author_id=550e8400-e29b-41d4-a716-446655440000`

**Response:**
```json
{
  "status_code": 200,
  "error_message": null,
  "data": {
    "posts": [
      {
        "id": "660e8400-e29b-41d4-a716-446655440000",
        "title": "My First Post",
        "content": "This is the content of my first post.",
        "created_by": "550e8400-e29b-41d4-a716-446655440000",
        "author": {
          "id": "550e8400-e29b-41d4-a716-446655440000",
          "username": "john_doe",
          "display_name": "John Doe"
        },
        "created_at": "2024-01-15T11:00:00Z",
        "updated_at": "2024-01-15T11:00:00Z"
      }
    ],
    "pagination": {
      "limit": 20,
      "offset": 0,
      "total": 1,
      "has_more": false
    }
  }
}
```

### Update Post
Update a post (only the author can update their post).

**Endpoint:** `PUT /api/v1/posts/post/{id}`

**Headers:** `Authorization: Bearer <token>`

**Path Parameters:**
- `id`: Post UUID

**Request Body:**
```json
{
  "title": "Updated Post Title",
  "content": "Updated post content with new information."
}
```

**Response:**
```json
{
  "status_code": 200,
  "error_message": null,
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "title": "Updated Post Title",
    "content": "Updated post content with new information.",
    "created_by": "550e8400-e29b-41d4-a716-446655440000",
    "author": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "display_name": "John Doe"
    },
    "created_at": "2024-01-15T11:00:00Z",
    "updated_at": "2024-01-15T11:30:00Z"
  }
}
```

### Delete Post
Delete a post (only the author can delete their post).

**Endpoint:** `DELETE /api/v1/posts/post/{id}`

**Headers:** `Authorization: Bearer <token>`

**Path Parameters:**
- `id`: Post UUID

**Response:**
```json
{
  "status_code": 200,
  "error_message": null,
  "data": {
    "message": "Post deleted successfully"
  }
}
```

---

## Comment Management Endpoints

### Create Comment
Create a new comment on a post (supports nested comments).

**Endpoint:** `POST /api/v1/posts/post-comments/{postId}`

**Headers:** `Authorization: Bearer <token>`

**Path Parameters:**
- `postId`: Post UUID

**Request Body:**
```json
{
  "content": "This is a comment on the post.",
  "parent_id": null
}
```

**For nested comment:**
```json
{
  "content": "This is a reply to another comment.",
  "parent_id": "770e8400-e29b-41d4-a716-446655440000"
}
```

**Response:**
```json
{
  "status_code": 201,
  "error_message": null,
  "data": {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "content": "This is a comment on the post.",
    "post_id": "660e8400-e29b-41d4-a716-446655440000",
    "parent_id": null,
    "thread_id": "770e8400-e29b-41d4-a716-446655440000",
    "path": ["770e8400-e29b-41d4-a716-446655440000"],
    "created_by": "550e8400-e29b-41d4-a716-446655440000",
    "author": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "display_name": "John Doe"
    },
    "created_at": "2024-01-15T11:15:00Z"
  }
}
```

### Get Comments for Post
Get all comments for a specific post with nested structure.

**Endpoint:** `GET /api/v1/posts/post-comments/{postId}`

**Path Parameters:**
- `postId`: Post UUID

**Query Parameters:**
- `limit` (optional): Number of comments per page (default: 10, max: 100)
- `offset` (optional): Number of comments to skip (default: 0)

**Example:** `GET /api/v1/posts/post-comments/660e8400-e29b-41d4-a716-446655440000?limit=20&offset=0`

**Response:**
```json
{
  "status_code": 200,
  "error_message": null,
  "data": {
    "comments": [
      {
        "id": "770e8400-e29b-41d4-a716-446655440000",
        "content": "This is a top-level comment.",
        "post_id": "660e8400-e29b-41d4-a716-446655440000",
        "parent_id": null,
        "thread_id": "770e8400-e29b-41d4-a716-446655440000",
        "path": ["770e8400-e29b-41d4-a716-446655440000"],
        "created_by": "550e8400-e29b-41d4-a716-446655440000",
        "author": {
          "id": "550e8400-e29b-41d4-a716-446655440000",
          "username": "john_doe",
          "display_name": "John Doe"
        },
        "created_at": "2024-01-15T11:15:00Z",
        "replies": [
          {
            "id": "880e8400-e29b-41d4-a716-446655440000",
            "content": "This is a reply to the comment.",
            "post_id": "660e8400-e29b-41d4-a716-446655440000",
            "parent_id": "770e8400-e29b-41d4-a716-446655440000",
            "thread_id": "770e8400-e29b-41d4-a716-446655440000",
            "path": ["770e8400-e29b-41d4-a716-446655440000", "880e8400-e29b-41d4-a716-446655440000"],
            "created_by": "990e8400-e29b-41d4-a716-446655440000",
            "author": {
              "id": "990e8400-e29b-41d4-a716-446655440000",
              "username": "jane_doe",
              "display_name": "Jane Doe"
            },
            "created_at": "2024-01-15T11:20:00Z",
            "replies": []
          }
        ]
      }
    ],
    "pagination": {
      "limit": 20,
      "offset": 0,
      "total": 1,
      "has_more": false
    }
  }
}
```

### Get Comment by ID
Get a specific comment by its ID.

**Endpoint:** `GET /api/v1/comments/{id}`

**Path Parameters:**
- `id`: Comment UUID

**Response:**
```json
{
  "status_code": 200,
  "error_message": null,
  "data": {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "content": "This is a comment.",
    "post_id": "660e8400-e29b-41d4-a716-446655440000",
    "parent_id": null,
    "thread_id": "770e8400-e29b-41d4-a716-446655440000",
    "path": ["770e8400-e29b-41d4-a716-446655440000"],
    "created_by": "550e8400-e29b-41d4-a716-446655440000",
    "author": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "display_name": "John Doe"
    },
    "created_at": "2024-01-15T11:15:00Z"
  }
}
```

### Update Comment
Update a comment (only the author can update their comment).

**Endpoint:** `PUT /api/v1/comments/{id}`

**Headers:** `Authorization: Bearer <token>`

**Path Parameters:**
- `id`: Comment UUID

**Request Body:**
```json
{
  "content": "Updated comment content."
}
```

**Response:**
```json
{
  "status_code": 200,
  "error_message": null,
  "data": {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "content": "Updated comment content.",
    "post_id": "660e8400-e29b-41d4-a716-446655440000",
    "parent_id": null,
    "thread_id": "770e8400-e29b-41d4-a716-446655440000",
    "path": ["770e8400-e29b-41d4-a716-446655440000"],
    "created_by": "550e8400-e29b-41d4-a716-446655440000",
    "author": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "john_doe",
      "display_name": "John Doe"
    },
    "created_at": "2024-01-15T11:15:00Z"
  }
}
```

### Delete Comment
Delete a comment (only the author can delete their comment).

**Endpoint:** `DELETE /api/v1/comments/{id}`

**Headers:** `Authorization: Bearer <token>`

**Path Parameters:**
- `id`: Comment UUID

**Response:**
```json
{
  "status_code": 200,
  "error_message": null,
  "data": {
    "message": "Comment deleted successfully"
  }
}
```

---

## Health Check Endpoint

### Health Check
Check if the service is running and healthy.

**Endpoint:** `GET /health`

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T11:30:00Z",
  "version": "1.0.0",
  "database": "connected"
}
```

---

## Error Handling

### Common Error Responses

#### Validation Error (400)
```json
{
  "status_code": 400,
  "error_message": "Validation failed: username is required",
  "data": null
}
```

#### Unauthorized (401)
```json
{
  "status_code": 401,
  "error_message": "Invalid or expired token",
  "data": null
}
```

#### Forbidden (403)
```json
{
  "status_code": 403,
  "error_message": "You don't have permission to access this resource",
  "data": null
}
```

#### Not Found (404)
```json
{
  "status_code": 404,
  "error_message": "User not found",
  "data": null
}
```

#### Conflict (409)
```json
{
  "status_code": 409,
  "error_message": "Username already exists",
  "data": null
}
```

#### Internal Server Error (500)
```json
{
  "status_code": 500,
  "error_message": "Internal server error",
  "data": null
}
```

---

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **Default limit**: 100 requests per minute per IP
- **Headers included in response**:
  - `X-RateLimit-Limit`: Maximum requests allowed
  - `X-RateLimit-Remaining`: Remaining requests in current window
  - `X-RateLimit-Reset`: Time when the rate limit resets

When rate limit is exceeded:
```json
{
  "status_code": 429,
  "error_message": "Rate limit exceeded. Try again later.",
  "data": null
}
```

---

## Pagination

List endpoints support pagination with the following parameters:

- `limit`: Number of items per page (default: 10, max: 100)
- `offset`: Number of items to skip (default: 0)

Pagination response format:
```json
{
  "pagination": {
    "limit": 10,
    "offset": 0,
    "total": 25,
    "has_more": true
  }
}
```

---

## Content Security

### HTML Sanitization
All user-generated content (posts and comments) is automatically sanitized to prevent XSS attacks. The following HTML tags and attributes are allowed:

**Allowed tags**: `p`, `br`, `strong`, `em`, `u`, `a`, `ul`, `ol`, `li`, `blockquote`, `code`, `pre`

**Allowed attributes**: `href` (for `a` tags only)

### Input Validation
All input is validated according to the following rules:

- **Username**: 3-50 characters, alphanumeric and underscores only
- **Email**: Valid email format
- **Password**: Minimum 8 characters
- **Display Name**: 1-100 characters
- **Post Title**: 1-200 characters
- **Post/Comment Content**: 1-10,000 characters

---

## Authentication Flow Example

### Complete Authentication Flow

1. **Register a new user**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "securepassword",
    "display_name": "John Doe"
  }'
```

2. **Login to get tokens**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "securepassword"
  }'
```

3. **Use access token for authenticated requests**:
```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

4. **Refresh token when access token expires**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

---

## Testing the API

### Using cURL

Create a post:
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my first post."
  }'
```

Get all posts:
```bash
curl -X GET "http://localhost:8080/api/v1/posts?limit=10&offset=0"
```

Create a comment:
```bash
curl -X POST http://localhost:8080/api/v1/posts/post-comments/POST_ID \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "This is a comment on the post."
  }'
```

### Using Postman

1. Import the API endpoints into Postman
2. Set up environment variables for `base_url` and `access_token`
3. Use the authentication endpoints to get tokens
4. Test all endpoints with proper authentication headers

---

## API Versioning

The API uses URL versioning with the `/api/v1/` prefix. Future versions will use `/api/v2/`, etc.

### Version History

- **v1.0.0**: Initial API release with user, post, and comment management

---

## Support

For API support and questions:
- Check the main [README.md](README.md) for setup instructions
- Open an issue in the GitHub repository
- Contact the development team

---

*Last updated: January 2024* 