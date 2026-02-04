# Mini-Moodle LMS API

Production-ready Learning Management System built with Go, featuring Clean Architecture, concurrency support, and containerization.

## Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+

### Local Setup

```bash
# Install dependencies
go mod download

# Set database URL
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/moodle"

# Apply migrations
make migrate-up

# Run API
make run
```

### Docker Deployment

```bash
# Start all services (PostgreSQL, Redis, API)
make docker-up

# Apply migrations
make migrate-up

# View logs
make docker-logs

# Stop services
make docker-down
```

API runs on `http://localhost:8080/api/v1`

## Architecture

**Clean Architecture** with 4 layers + Concurrency:

- **Domain**: Pure business entities and interfaces
- **Repository**: PostgreSQL data access layer
- **Usecase**: Business logic and services
- **Delivery**: HTTP API handlers
- **App**: Background workers and graceful shutdown

## Key Features

- 10 domain services (User, Course, Enrollment, Assignment, Submission, Grade, Attendance, Chat, Notification, Upload)
- Goroutine-based background workers with channel-based processing
- Database health monitoring
- Asynchronous notification system
- Graceful shutdown with worker cleanup
- PostgreSQL with connection pooling
- Redis support for caching
- Comprehensive error handling
- Request validation

## Commands

```bash
make build          # Compile binary
make run            # Run locally
make test           # Run tests
make clean          # Clean build artifacts
make docker-up      # Start Docker containers
make docker-down    # Stop containers
make docker-logs    # View API logs
make migrate-up     # Apply migrations
make migrate-down   # Rollback migrations
```

## API Endpoints

All endpoints prefixed with `/api/v1/`

### Users
- `POST /users` - Register
- `GET /users` - List all
- `GET /users/{id}` - Get user
- `PATCH /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user

### Courses
- `POST /courses` - Create course
- `GET /courses` - List courses
- `GET /courses/{id}` - Get course details
- `PATCH /courses/{id}` - Update course
- `DELETE /courses/{id}` - Delete course

### Enrollments
- `POST /enrollments` - Enroll student
- `GET /enrollments/course/{courseID}` - List by course
- `GET /enrollments/student/{studentID}` - List by student
- `DELETE /enrollments/{id}` - Remove enrollment

Similar patterns for Assignments, Submissions, Grades, Attendance, Chat, Notifications, Uploads.

## Concurrency

Background workers run asynchronously:

```go
// Send notification (goes to background worker)
app.SendNotification(userID, message, "assignment_submitted")

// Health check runs every 30 seconds
// Graceful shutdown waits for all workers to finish
```

## Documentation

- [docs/architecture.md](docs/architecture.md) - Detailed architecture documentation
- [docs/swagger.yaml](docs/swagger.yaml) - API specifications

## Environment Variables

```
DATABASE_URL=postgres://user:password@localhost:5432/moodle
REDIS_URL=redis://localhost:6379
PORT=8080
```

## Technologies

- **Framework**: Gin Web Framework
- **Database**: PostgreSQL with pgx connection pooling
- **Concurrency**: Go goroutines with channels
- **Containerization**: Docker + docker-compose
- **Language**: Go 1.21+

## Testing

```bash
# Unit tests
go test ./tests/unit/...

# All tests with coverage
go test -v ./... -cover
```

## Performance

- Multi-stage Docker build (optimized ~120MB image)
- Connection pooling (default 25 max connections)
- Indexed database queries
- Asynchronous background processing
- Health monitoring on 30-second intervals

## Project Structure

```
backend/
├── cmd/api/main.go              # API entry point
├── internal/
│   ├── app/                     # Application init + workers
│   ├── domain/                  # Business entities (10 services)
│   ├── repository/              # PostgreSQL implementations
│   ├── usecase/                 # Business logic services
│   ├── delivery/http/           # HTTP handlers
│   └── shared/                  # Utilities
├── migrations/                  # Database migrations
├── tests/                       # Test suites
├── docs/                        # Documentation
├── Dockerfile                   # Container build
├── docker-compose.yml          # Development environment
├── Makefile                    # Build commands
└── go.mod                      # Dependencies
```

## License

Private - Educational Use Only

## Status

Production Ready - All components verified and compiled successfully.
