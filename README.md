# Post Comments Service

A full-stack web application for managing posts and comments with a modern React frontend and a robust Go backend service.

## 🌐 Live Deployments

- **Frontend**: [https://post-comments-service-three.vercel.app/](https://post-comments-service-three.vercel.app/)
- **Backend API**: [https://post-comments-service.onrender.com](https://post-comments-service.onrender.com)

## 📸 Screenshots
<img width="1280" alt="Screenshot 2025-06-01 at 3 16 38 AM" src="https://github.com/user-attachments/assets/6805ecc9-f8a7-40e3-a6c3-d43e9bbcf90e" />
<img width="1276" alt="Screenshot 2025-06-01 at 3 17 35 AM" src="https://github.com/user-attachments/assets/cbc63fd3-f100-41ea-866c-fb9d7e506991" />




## 🏗️ Architecture Overview

This project consists of two main components:

### Frontend (`/client`)
- **Framework**: React 18 with TypeScript and Vite
- **Styling**: Tailwind CSS v4 with shadcn/ui components
- **Features**: Modern UI, infinite scrolling, responsive design
- **Deployment**: Vercel

### Backend (`/server`)
- **Language**: Go (Golang) with Gin framework
- **Database**: PostgreSQL with migrations
- **Features**: JWT authentication, nested comments, RESTful API
- **Deployment**: Render

## 🚀 Quick Start

### Prerequisites
- Node.js 18+ (for frontend)
- Go 1.21+ (for backend)
- PostgreSQL 12+ (for backend)
- Docker & Docker Compose (optional)

### Frontend Setup
```bash
cd client
npm install
npm run dev
```

### Backend Setup
```bash
cd server
go mod tidy
# Configure .env file (see server/README.md)
make migrate-up
go run cmd/main.go
```

## 📚 Detailed Documentation

### Frontend Documentation
For detailed frontend setup, features, and development guide:
- **[Client README](./client/README.md)** - Complete frontend documentation

### Backend Documentation
For detailed backend setup, API documentation, and architecture:
- **[Server README](./server/README.md)** - Complete backend documentation
- **[API Documentation](./server/API.md)** - Comprehensive API reference
- **[Architecture Guide](./server/ARCHITECTURE.md)** - System design and database schema

## 🛠️ Tech Stack

### Frontend
- React 18 + TypeScript
- Vite (build tool)
- Tailwind CSS v4
- shadcn/ui components
- React Router
- Axios

### Backend
- Go (Golang)
- Gin web framework
- PostgreSQL
- JWT authentication
- Docker support
- Comprehensive logging

## 🌟 Key Features

- **User Management**: Authentication and user profiles
- **Post Management**: Create, read, update, and delete posts
- **Nested Comments**: Threaded comments with unlimited nesting
- **Modern UI**: Responsive design with smooth animations
- **Infinite Scrolling**: Optimized content loading
- **Real-time Updates**: Interactive comment system
- **Secure API**: JWT-based authentication
- **Production Ready**: Deployed and scalable

## 🚀 Deployment

Both applications are deployed and ready to use:

- **Frontend** is deployed on Vercel with automatic deployments from the main branch
- **Backend** is deployed on Render with PostgreSQL database

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

---

For specific setup instructions and detailed documentation, please refer to the individual README files in the `client` and `server` directories. 
