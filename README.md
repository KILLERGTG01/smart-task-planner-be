# Smart Task Planner Backend

A high-performance Go backend API for intelligent task planning and management using AI-powered plan generation.

## Tech Stack

### Core Technologies
- **Go 1.25+** - Primary backend language
- **Fiber v2** - High-performance HTTP web framework
- **PostgreSQL** - Primary database with connection pooling
- **Zap** - Structured logging library
- **Auth0** - Authentication and authorization
- **Google Gemini AI** - AI-powered task plan generation

### Key Libraries
- **pgx/v5** - PostgreSQL driver and connection pooling
- **golang-migrate** - Database migration management
- **go-playground/validator** - Request validation
- **golang-jwt** - JWT token handling
- **keyfunc** - JWKS key management

## Architecture

### Project Structure
```
smart-task-planner-be/
├── cmd/
│   ├── server/          # Main application entry point
│   └── migrate/         # Database migration tool
├── internal/
│   ├── config/          # Configuration management
│   ├── db/              # Database connection and migrations
│   ├── handlers/        # HTTP request handlers
│   ├── logger/          # Logging configuration
│   ├── middleware/      # Authentication middleware
│   ├── routes/          # Route definitions
│   ├── services/        # Business logic services
│   └── validation/      # Request validation
├── migrations/          # Database migration files
├── .env.example         # Environment configuration template
├── Makefile            # Build and development commands
└── nginx.conf.example  # Nginx configuration for production
```

### System Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Nginx         │    │   Go Backend    │
│   (React/Vue)   │◄──►│   (Rate Limit   │◄──►│   (Fiber)       │
│                 │    │    & Proxy)     │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                       ┌─────────────────┐             │
                       │   Auth0         │◄────────────┤
                       │   (JWT Auth)    │             │
                       └─────────────────┘             │
                                                        │
                       ┌─────────────────┐             │
                       │   PostgreSQL    │◄────────────┤
                       │   (Database)    │             │
                       └─────────────────┘             │
                                                        │
                       ┌─────────────────┐             │
                       │   Gemini AI     │◄────────────┘
                       │   (Plan Gen)    │
                       └─────────────────┘
```

### Key Features

#### Authentication & Security
- **JWT-based authentication** via Auth0
- **JWKS key rotation** support
- **Request validation** with structured error responses
- **Security headers** via Helmet middleware
- **CORS configuration** with environment-based origins

#### Performance & Scalability
- **Connection pooling** for database efficiency
- **Rate limiting** handled by Nginx
- **Structured logging** with Zap for observability
- **Graceful shutdown** handling
- **Fast startup** with separated migration process

#### AI Integration
- **Gemini AI integration** for intelligent task planning
- **Context-aware timeouts** for AI requests
- **JSON response parsing** with error handling

## Getting Started

### Prerequisites
- Go 1.25+
- PostgreSQL database
- Auth0 account
- Google Gemini API key

### Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd smart-task-planner-be
```

2. **Install dependencies**
```bash
make install-deps
```

3. **Setup environment**
```bash
cp .env.example .env
# Edit .env with your actual credentials
```

4. **Run database migrations**
```bash
make migrate
```

5. **Start the server**
```bash
make run
```

### Development Commands

```bash
# Build the application
make build

# Run with migrations (full setup)
make dev

# Run migrations only
make migrate

# Run tests
make test

# Clean build artifacts
make clean
```

## API Endpoints

### Authentication Required
All API endpoints require valid JWT token in Authorization header:
```
Authorization: Bearer <jwt_token>
```

### Endpoints

#### Health Check
```http
GET /health
```
Returns server status and timestamp.

#### Generate Plan
```http
POST /api/generate
Content-Type: application/json

{
  "goal": "Learn React in 30 days",
  "title": "React Learning Plan"
}
```

#### Get Plan History
```http
GET /api/history
```
Returns user's plan generation history.

## Configuration

### Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `APP_ENV` | Application environment | `development` |
| `APP_PORT` | Server port | `8080` |
| `DATABASE_URL` | PostgreSQL connection string | `postgresql://...` |
| `AUTH0_DOMAIN` | Auth0 domain | `https://your-domain.auth0.com` |
| `AUTH0_AUDIENCE` | Auth0 API audience | `https://your-api` |
| `AUTH0_ISSUER` | Auth0 issuer URL | `https://your-domain.auth0.com/` |
| `GEMINI_API_KEY` | Google Gemini API key | `your_api_key` |
| `GEMINI_BASE_URL` | Gemini API base URL | `https://generativelanguage.googleapis.com` |
| `ALLOWED_ORIGINS` | CORS allowed origins | `*` or `https://yourdomain.com` |

## Deployment

### Production Setup

1. **Build the application**
```bash
make build
```

2. **Setup Nginx** (recommended)
```bash
# Use provided nginx.conf.example
# Configure rate limiting and SSL
```

3. **Environment Configuration**
```bash
# Set production environment variables
export APP_ENV=production
export ALLOWED_ORIGINS=https://yourdomain.com
```

4. **Database Migration**
```bash
make migrate
```

5. **Start the server**
```bash
./bin/server
```

### Docker Support
```dockerfile
# Example Dockerfile structure
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

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
