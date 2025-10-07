# LinkUp - Professional Team Communication Platform

A modern, feature-rich team communication platform built with Go and React, designed for professional teams and organizations.

## 🚀 Features

### Core Features
- **Real-time Messaging**: WebSocket-based instant messaging
- **Room Management**: Create public/private rooms with custom permissions
- **User Authentication**: JWT-based auth with 2FA support
- **File Sharing**: Upload and share files with drag & drop
- **Message Search**: Full-text search across all messages
- **Typing Indicators**: See when others are typing
- **Message Reactions**: React to messages with emojis
- **Online Presence**: Real-time user status tracking

### Advanced Features
- **Role-based Access Control**: Admin, moderator, and user roles
- **Two-Factor Authentication**: Enhanced security with TOTP
- **Poll System**: Create polls and vote on decisions
- **Mentions**: @mention users for notifications
- **Achievement System**: Gamification with levels and badges
- **Analytics Dashboard**: User activity and engagement metrics
- **Admin Panel**: Comprehensive user and system management
- **Content Moderation**: AI-powered content filtering
- **GitHub Integration**: Connect with GitHub for notifications
- **Calendar Integration**: Schedule meetings and events
- **Music Rooms**: Collaborative music listening
- **Chat Games**: Interactive games within chat
- **Push Notifications**: Real-time notifications
- **AI Assistant**: AI-powered chat assistance

### Technical Features
- **RESTful API**: Well-documented REST API
- **WebSocket Support**: Real-time bidirectional communication
- **Database**: PostgreSQL with GORM ORM
- **File Storage**: Local and cloud storage options
- **Swagger Documentation**: Interactive API documentation
- **Docker Support**: Containerized deployment
- **Health Monitoring**: System health checks
- **Logging**: Comprehensive logging system

## 🏗️ Architecture

### Backend (Go)
- **Framework**: Gin web framework
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT with 2FA support
- **WebSocket**: Gorilla WebSocket
- **File Upload**: Multipart form handling
- **API Documentation**: Swagger/OpenAPI

### Frontend (React)
- **Framework**: React 18 with TypeScript
- **Styling**: Tailwind CSS
- **State Management**: Zustand
- **Routing**: React Router
- **HTTP Client**: Axios
- **WebSocket**: Socket.io client
- **Forms**: React Hook Form
- **Animations**: Framer Motion

## 📁 Project Structure

```
LinkUp/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── app/
│   │   ├── app.go           # Application setup
│   │   └── deps.go          # Dependencies
│   ├── auth/
│   │   ├── jwt.go           # JWT authentication
│   │   └── auth.go          # Auth middleware
│   ├── handlers/
│   │   ├── auth.go          # Authentication handlers
│   │   ├── rooms.go         # Room management
│   │   ├── messages.go      # Message handling
│   │   ├── ws.go            # WebSocket handlers
│   │   ├── upload.go        # File upload
│   │   ├── search.go        # Message search
│   │   ├── types.go         # Request/response types
│   │   └── extended_handlers.go # Extended features
│   ├── models/
│   │   └── models.go        # Database models
│   ├── storage/
│   │   └── db.go            # Database connection
│   ├── utils/
│   │   └── hash.go          # Password hashing
│   └── err/
│       └── api_errors.go    # Error handling
├── docs/
│   ├── docs.go              # Swagger docs
│   ├── swagger.json         # Swagger JSON
│   └── swagger.yaml         # Swagger YAML
├── frontend/                # React frontend
│   ├── src/
│   │   ├── components/      # React components
│   │   ├── pages/          # Page components
│   │   ├── services/       # API services
│   │   ├── stores/        # State management
│   │   ├── types/         # TypeScript types
│   │   └── ...
│   ├── package.json
│   └── README.md
├── uploads/                 # File uploads
├── go.mod
├── go.sum
└── README.md
```

## 🚀 Quick Start

### Prerequisites
- Go 1.19+
- Node.js 16+
- PostgreSQL 12+
- Docker (optional)

### Backend Setup

1. **Clone the repository**:
```bash
git clone <repository-url>
cd LinkUp
```

2. **Install dependencies**:
```bash
go mod download
```

3. **Set up database**:
```bash
# Create PostgreSQL database
createdb linkup

# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=linkup
export JWT_SECRET=your-secret-key
```

4. **Run the server**:
```bash
go run cmd/server/main.go
```

The API will be available at `http://localhost:8080`

### Frontend Setup

1. **Navigate to frontend directory**:
```bash
cd frontend
```

2. **Install dependencies**:
```bash
npm install
```

3. **Start development server**:
```bash
npm start
```

The frontend will be available at `http://localhost:3000`

### Docker Setup

1. **Build and run with Docker Compose**:
```bash
docker-compose up --build
```

## 📚 API Documentation

Once the server is running, visit:
- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **API Docs**: `http://localhost:8080/docs/`

## 🔧 Configuration

### Environment Variables

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=5432
DB_NAME=chat

# JWT
JWT_SECRET=your-secret-key

# Server
PORT=8080
HOST=localhost

# File Upload
UPLOAD_PATH=./uploads
MAX_FILE_SIZE=10485760

# 2FA
TOTP_ISSUER=LinkUp

# AI Assistant
OPENAI_API_KEY=your-openai-key

# GitHub Integration
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
```

## 🧪 Testing

### Backend Tests
```bash
go test ./...
```

### Frontend Tests
```bash
cd frontend
npm test
```

## 📦 Deployment

### Production Build

1. **Build frontend**:
```bash
cd frontend
npm run build
```

2. **Build backend**:
```bash
go build -o linkup cmd/server/main.go
```

3. **Run production server**:
```bash
./linkup
```

### Docker Deployment

```bash
docker build -t linkup .
docker run -p 8080:8080 linkup
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin](https://github.com/gin-gonic/gin) - Web framework
- [GORM](https://gorm.io/) - ORM library
- [React](https://reactjs.org/) - Frontend framework
- [Tailwind CSS](https://tailwindcss.com/) - CSS framework
- [Zustand](https://github.com/pmndrs/zustand) - State management

## 📞 Support

For support, email support@linkup.com or join our Discord server.

---

**LinkUp** - Connecting teams, one message at a time. 🚀