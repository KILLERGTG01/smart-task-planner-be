# 🎯 Smart Task Planner Backend

<div align="center">

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Fiber](https://img.shields.io/badge/fiber-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Auth0](https://img.shields.io/badge/Auth0-EB5424?style=for-the-badge&logo=auth0&logoColor=white)
![Google AI](https://img.shields.io/badge/Google%20AI-4285F4?style=for-the-badge&logo=google&logoColor=white)

**A high-performance Go backend API for intelligent task planning powered by Google's Gemini AI**

[🚀 Live API](https://api.anurag-goel.com) • [📖 Documentation](#api-documentation) • [🔧 Setup](#installation) • [🧪 Try on Scalar](#try-the-api)

</div>

---

## 📋 Table of Contents

- [✨ Features](#-features)
- [🏗️ Architecture](#️-architecture)
- [🚀 Tech Stack](#-tech-stack)
- [📁 Project Structure](#-project-structure)
- [⚡ Installation](#-installation)
- [🔧 Configuration](#-configuration)
- [🧪 Try the API](#-try-the-api)
- [📖 API Documentation](#-api-documentation)
- [🔐 Authentication](#-authentication)
- [💡 Examples](#-examples)
- [🚀 Deployment](#-deployment)
- [🤝 Contributing](#-contributing)
- [📄 License](#-license)

---

## ✨ Features

### 🎯 **Core Functionality**
- **AI-Powered Plan Generation** - Create detailed task plans using Google Gemini 2.5 Flash Lite
- **Anonymous & Authenticated Access** - Generate plans without signup, save with authentication
- **Real-time Streaming** - Server-sent events for live plan generation updates
- **Plan History** - Save and retrieve previous plans for authenticated users

### 🔐 **Authentication & Security**
- **OAuth Integration** - Google & GitHub login via Auth0
- **JWT Token Management** - Secure token-based authentication
- **Rate Limiting** - Nginx-based rate limiting for API protection
- **CORS Support** - Configurable cross-origin resource sharing

### ⚡ **Performance & Reliability**
- **High Performance** - Built with Go Fiber for maximum throughput
- **Connection Pooling** - Optimized PostgreSQL connections
- **Graceful Shutdown** - Clean server termination handling
- **Structured Logging** - Comprehensive logging with Zap
- **Fallback Plans** - Guaranteed response even if AI fails

---

## 🏗️ Architecture

```mermaid
graph TB
    subgraph "Client Layer"
        FE[Frontend App]
        API_CLIENT[API Clients]
    end
    
    subgraph "Load Balancer"
        NGINX[Nginx<br/>Rate Limiting & SSL]
    end
    
    subgraph "Application Layer"
        GO_APP[Go Fiber Server<br/>Port 8080]
        
        subgraph "Handlers"
            AUTH_H[Auth Handlers]
            PLAN_H[Plan Handlers]
            HEALTH_H[Health Handler]
        end
        
        subgraph "Middleware"
            AUTH_M[JWT Auth]
            CORS_M[CORS]
            LOG_M[Logging]
        end
        
        subgraph "Services"
            GEMINI_S[Gemini AI Service]
            VALIDATION_S[Validation Service]
        end
    end
    
    subgraph "External Services"
        AUTH0[Auth0<br/>OAuth Provider]
        GEMINI[Google Gemini AI<br/>2.5 Flash Lite]
    end
    
    subgraph "Database Layer"
        POSTGRES[(PostgreSQL<br/>Connection Pool)]
    end
    
    FE --> NGINX
    API_CLIENT --> NGINX
    NGINX --> GO_APP
    
    GO_APP --> AUTH_H
    GO_APP --> PLAN_H
    GO_APP --> HEALTH_H
    
    AUTH_H --> AUTH0
    PLAN_H --> GEMINI_S
    GEMINI_S --> GEMINI
    
    GO_APP --> POSTGRES
    
    AUTH_M --> AUTH0
    
    style GO_APP fill:#00ADD8
    style POSTGRES fill:#316192
    style GEMINI fill:#4285F4
    style AUTH0 fill:#EB5424
```

### 🔄 **Request Flow**

1. **Client Request** → Nginx (Rate Limiting + SSL)
2. **Nginx** → Go Fiber Server
3. **Middleware Chain** → CORS → Auth (if required) → Logging
4. **Route Handler** → Business Logic
5. **External APIs** → Gemini AI / Auth0 (if needed)
6. **Database** → PostgreSQL (for authenticated users)
7. **Response** → JSON/SSE back to client

---

## 🚀 Tech Stack

### **Backend Core**
- **[Go 1.25+](https://golang.org/)** - Primary language for high performance
- **[Fiber v2](https://gofiber.io/)** - Express-inspired web framework
- **[PostgreSQL](https://postgresql.org/)** - Primary database with ACID compliance
- **[pgx/v5](https://github.com/jackc/pgx)** - High-performance PostgreSQL driver

### **AI & Authentication**
- **[Google Gemini 2.5 Flash Lite](https://ai.google.dev/)** - AI-powered task generation
- **[Auth0](https://auth0.com/)** - OAuth authentication provider
- **[JWT](https://jwt.io/)** - Stateless authentication tokens

### **Infrastructure & Monitoring**
- **[Zap](https://github.com/uber-go/zap)** - Structured, high-performance logging
- **[Nginx](https://nginx.org/)** - Reverse proxy, load balancing, rate limiting
- **[Docker](https://docker.com/)** - Containerization for deployment
- **[golang-migrate](https://github.com/golang-migrate/migrate)** - Database migrations

### **Development Tools**
- **[Air](https://github.com/cosmtrek/air)** - Live reload for development
- **[go-playground/validator](https://github.com/go-playground/validator)** - Request validation
- **[godotenv](https://github.com/joho/godotenv)** - Environment variable management

---

## 📁 Project Structure

```
smart-task-planner-be/
├── 📁 cmd/
│   ├── 📁 server/           # Main application entry point
│   │   └── main.go
│   └── 📁 migrate/          # Database migration tool
│       └── main.go
├── 📁 internal/
│   ├── 📁 config/           # Configuration management
│   │   └── config.go
│   ├── 📁 db/               # Database connection & models
│   │   ├── connection.go
│   │   └── models.go
│   ├── 📁 handlers/         # HTTP request handlers
│   │   ├── auth_handler.go
│   │   ├── plan_handler.go
│   │   └── test_handler.go
│   ├── 📁 logger/           # Logging configuration
│   │   └── logger.go
│   ├── 📁 middleware/       # HTTP middleware
│   │   ├── auth.go
│   │   └── config.go
│   ├── 📁 routes/           # Route definitions
│   │   ├── routes.go
│   │   ├── auth.go
│   │   ├── health.go
│   │   └── plan.go
│   ├── 📁 services/         # Business logic services
│   │   └── gemini_service.go
│   └── 📁 validation/       # Request validation
│       └── validator.go
├── 📁 migrations/           # Database migration files
├── 📄 .env.example          # Environment variables template
├── 📄 .gitignore           # Git ignore rules
├── 📄 Makefile             # Build and development commands
├── 📄 go.mod               # Go module dependencies
├── 📄 go.sum               # Go module checksums
├── 📄 nginx.conf.example   # Nginx configuration template
├── 📄 README.md            # This file
└── 📄 LICENSE              # MIT License
```

---

## ⚡ Installation

### **Prerequisites**
- **Go 1.25+** - [Download](https://golang.org/dl/)
- **PostgreSQL 13+** - [Download](https://postgresql.org/download/)
- **Auth0 Account** - [Sign up](https://auth0.com/)
- **Google AI API Key** - [Get API Key](https://ai.google.dev/)

### **1. Clone Repository**
```bash
git clone https://github.com/KILLERGTG01/smart-task-planner-be.git
cd smart-task-planner-be
```

### **2. Install Dependencies**
```bash
make install-deps
# or
go mod download && go mod tidy
```

### **3. Environment Setup**
```bash
cp .env.example .env
# Edit .env with your configuration
```

### **4. Database Setup**
```bash
# Run migrations
make migrate
# or
./bin/migrate
```

### **5. Start Development Server**
```bash
# Full development setup (migrations + server)
make dev

# Or just run the server
make run
```

The server will start on `http://localhost:8080`

---

## 🔧 Configuration

### **Environment Variables**

Create a `.env` file with the following configuration:

```bash
# Application
APP_ENV=development
APP_PORT=8080
LOG_LEVEL=debug

# Database
DATABASE_URL=postgresql://username:password@localhost:5432/smart_planner?sslmode=disable

# Auth0 Configuration
AUTH0_DOMAIN=your-domain.auth0.com
AUTH0_AUDIENCE=https://your-api-identifier
AUTH0_ISSUER=https://your-domain.auth0.com/
AUTH0_CLIENT_ID=your_client_id
AUTH0_CLIENT_SECRET=your_client_secret
AUTH0_MANAGEMENT_TOKEN=your_management_token
AUTH0_REDIRECT_URI=http://localhost:8080/auth/callback

# Google Gemini AI
GEMINI_API_KEY=your_gemini_api_key
GEMINI_BASE_URL=https://generativelanguage.googleapis.com

# CORS & Frontend
FRONTEND_URL=http://localhost:3000
ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com

# Optional
RATE_LIMIT_REQUESTS=10
RATE_LIMIT_PERIOD_SECONDS=1
TRUST_PROXY=false
```

### **Auth0 Setup**

1. **Create Application** in Auth0 Dashboard
2. **Configure Settings:**
   ```
   Application Type: Regular Web Application
   
   Allowed Callback URLs:
   http://localhost:8080/auth/callback,
   https://yourdomain.com/auth/callback
   
   Allowed Logout URLs:
   http://localhost:3000,
   https://yourdomain.com
   
   Allowed Web Origins:
   http://localhost:3000,
   https://yourdomain.com
   ```

3. **Enable Social Connections:**
   - Google OAuth
   - GitHub OAuth

4. **Create API** for JWT tokens

---

## 🧪 Try the API

### **🔗 Interactive API Documentation**

**Try the live API using Scalar:**

[![Try with Scalar](https://img.shields.io/badge/Try%20API-Scalar-blue?style=for-the-badge)](https://api.anurag-goel.com/scalar)

**Or use the Scalar CLI:**

```bash
# Install Scalar CLI
npm install -g @scalar/cli

# Generate interactive docs
scalar serve https://api.anurag-goel.com/openapi.json
```

### **🚀 Quick Test**

```bash
# Test health endpoint
curl https://api.anurag-goel.com/health

# Generate a plan (no auth required)
curl -X POST https://api.anurag-goel.com/api/generate \
  -H "Content-Type: application/json" \
  -d '{"goal": "Learn React in 30 days", "title": "React Learning Plan"}'

# Test streaming endpoint
curl -X POST https://api.anurag-goel.com/api/generate/stream \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -d '{"goal": "Build a mobile app", "title": "Mobile App Development"}'
```

---

## 📖 API Documentation

### **🔓 Public Endpoints**

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/health` | Health check | ❌ |
| `POST` | `/api/generate` | Generate task plan | ❌ |
| `POST` | `/api/generate/stream` | Generate plan (streaming) | ❌ |
| `GET` | `/auth/login` | Get OAuth login URL | ❌ |
| `GET` | `/auth/callback` | OAuth callback handler | ❌ |
| `POST` | `/auth/refresh` | Refresh JWT token | ❌ |
| `GET` | `/auth/logout` | Get logout URL | ❌ |

### **🔐 Protected Endpoints**

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/api/history` | Get user's plan history | ✅ |
| `GET` | `/auth/profile` | Get user profile | ✅ |

### **📊 Request/Response Examples**

#### **Generate Plan**
```http
POST /api/generate
Content-Type: application/json

{
  "goal": "Learn Python programming",
  "title": "Python Mastery Journey"
}
```

**Response:**
```json
{
  "plan": [
    {
      "task": "Set up Python development environment",
      "duration_days": 1,
      "depends_on": []
    },
    {
      "task": "Learn Python syntax and basic concepts",
      "duration_days": 7,
      "depends_on": ["Set up Python development environment"]
    },
    {
      "task": "Practice with data structures and algorithms",
      "duration_days": 10,
      "depends_on": ["Learn Python syntax and basic concepts"]
    },
    {
      "task": "Build a complete Python project",
      "duration_days": 14,
      "depends_on": ["Practice with data structures and algorithms"]
    }
  ],
  "saved": false,
  "message": "Plan generated but not saved. Login to save your plans."
}
```

#### **Streaming Plan Generation**
```http
POST /api/generate/stream
Content-Type: application/json
Accept: text/event-stream

{
  "goal": "Start a fitness routine",
  "title": "30-Day Fitness Challenge"
}
```

**Server-Sent Events Response:**
```
event: status
data: {"message": "Starting plan generation..."}

event: progress
data: {"message": "Plan generated successfully!"}

event: plan
data: {"plan": [{"task": "Assess current fitness level", "duration_days": 1, "depends_on": []}]}

event: complete
data: {"saved": false}
```

---

## 🔐 Authentication

### **OAuth Flow**

```mermaid
sequenceDiagram
    participant Client
    participant Backend
    participant Auth0
    participant Google/GitHub

    Client->>Backend: GET /auth/login?provider=google
    Backend->>Client: {auth_url, state}
    Client->>Auth0: Redirect to auth_url
    Auth0->>Google/GitHub: OAuth flow
    Google/GitHub->>Auth0: User consent
    Auth0->>Backend: GET /auth/callback?code=xxx
    Backend->>Auth0: Exchange code for tokens
    Auth0->>Backend: {access_token, user_info}
    Backend->>Client: Redirect with token
```

### **Using JWT Tokens**

```javascript
// Store token after successful auth
localStorage.setItem('access_token', token);

// Use token for API calls
fetch('/api/history', {
  headers: {
    'Authorization': `Bearer ${localStorage.getItem('access_token')}`
  }
});
```

---

## 💡 Examples

### **🔓 Anonymous Plan Generation**

**Request:**
```bash
curl -X POST https://api.anurag-goel.com/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "goal": "Learn machine learning in 60 days",
    "title": "ML Mastery Path"
  }'
```

**Response:**
```json
{
  "plan": [
    {
      "task": "Understand mathematics foundations (linear algebra, statistics)",
      "duration_days": 14,
      "depends_on": []
    },
    {
      "task": "Learn Python for data science (NumPy, Pandas, Matplotlib)",
      "duration_days": 10,
      "depends_on": ["Understand mathematics foundations (linear algebra, statistics)"]
    },
    {
      "task": "Study machine learning algorithms and concepts",
      "duration_days": 15,
      "depends_on": ["Learn Python for data science (NumPy, Pandas, Matplotlib)"]
    },
    {
      "task": "Practice with scikit-learn and real datasets",
      "duration_days": 12,
      "depends_on": ["Study machine learning algorithms and concepts"]
    },
    {
      "task": "Build and deploy a complete ML project",
      "duration_days": 9,
      "depends_on": ["Practice with scikit-learn and real datasets"]
    }
  ],
  "saved": false,
  "message": "Plan generated but not saved. Login to save your plans."
}
```

### **🔐 Authenticated Plan Generation**

**Request:**
```bash
curl -X POST https://api.anurag-goel.com/api/generate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "goal": "Build a SaaS application",
    "title": "SaaS Development Journey"
  }'
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "plan": [
    {
      "task": "Market research and competitor analysis",
      "duration_days": 7,
      "depends_on": []
    },
    {
      "task": "Define MVP features and user stories",
      "duration_days": 5,
      "depends_on": ["Market research and competitor analysis"]
    },
    {
      "task": "Design system architecture and database schema",
      "duration_days": 8,
      "depends_on": ["Define MVP features and user stories"]
    },
    {
      "task": "Develop backend API and authentication",
      "duration_days": 21,
      "depends_on": ["Design system architecture and database schema"]
    },
    {
      "task": "Build frontend user interface",
      "duration_days": 18,
      "depends_on": ["Develop backend API and authentication"]
    },
    {
      "task": "Implement payment processing and billing",
      "duration_days": 10,
      "depends_on": ["Build frontend user interface"]
    },
    {
      "task": "Testing, deployment, and launch preparation",
      "duration_days": 12,
      "depends_on": ["Implement payment processing and billing"]
    }
  ],
  "saved": true
}
```

### **📚 Plan History (Authenticated)**

**Request:**
```bash
curl -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..." \
  https://api.anurag-goel.com/api/history
```

**Response:**
```json
{
  "plans": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "SaaS Development Journey",
      "goal": "Build a SaaS application",
      "plan": [
        {
          "task": "Market research and competitor analysis",
          "duration_days": 7,
          "depends_on": []
        }
      ],
      "createdAt": "2024-01-15T10:30:00Z"
    },
    {
      "id": "660f9511-f3ac-52e5-b827-557766551111",
      "title": "ML Mastery Path",
      "goal": "Learn machine learning in 60 days",
      "plan": [
        {
          "task": "Understand mathematics foundations",
          "duration_days": 14,
          "depends_on": []
        }
      ],
      "createdAt": "2024-01-14T15:20:00Z"
    }
  ]
}
```

---

## 🚀 Deployment

### **🐳 Docker Deployment**

```dockerfile
# Dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/server .
CMD ["./server"]
```

```bash
# Build and run
docker build -t smart-task-planner-be .
docker run -p 8080:8080 --env-file .env smart-task-planner-be
```

### **🔧 Production Setup**

1. **Build the application:**
```bash
make build
```

2. **Setup Nginx** (use `nginx.conf.example`):
```nginx
server {
    listen 443 ssl;
    server_name api.yourdomain.com;
    
    location /api/ {
        limit_req zone=api burst=20 nodelay;
        proxy_pass http://localhost:8080;
    }
}
```

3. **Run migrations:**
```bash
make migrate
```

4. **Start the server:**
```bash
./bin/server
```

### **📊 Monitoring & Scaling**

- **Health Check**: `GET /health`
- **Metrics**: Structured logging with Zap
- **Scaling**: Stateless design allows horizontal scaling
- **Database**: Connection pooling for high concurrency

---

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md).

### **Development Workflow**

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/amazing-feature`
3. **Commit** your changes: `git commit -m 'Add amazing feature'`
4. **Push** to the branch: `git push origin feature/amazing-feature`
5. **Open** a Pull Request

### **Code Standards**

- Follow Go best practices and idioms
- Write tests for new functionality
- Use structured logging with Zap
- Validate all inputs
- Handle errors gracefully

---

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

**Built with ❤️ by [Anurag Goel](https://github.com/KILLERGTG01)**

[⭐ Star this repo](https://github.com/KILLERGTG01/smart-task-planner-be) • [🐛 Report Bug](https://github.com/KILLERGTG01/smart-task-planner-be/issues) • [💡 Request Feature](https://github.com/KILLERGTG01/smart-task-planner-be/issues)

</div>