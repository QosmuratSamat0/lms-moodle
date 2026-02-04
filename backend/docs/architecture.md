# Mini-Moodle LMS - Clean Architecture

Compact and clean LMS project using Clean Architecture.

## Project Structure

```
backend/
├── cmd/api/main.go                  # Entry point
├── internal/
│   ├── domain/                      # Domain entities + interfaces
│   │   ├── user/                    # User domain
│   │   ├── course/                  # Course domain
│   │   ├── enrollment/              # Student enrollments
│   │   ├── assignment/              # Assignments
│   │   ├── submission/              # Submissions
│   │   ├── grade/                   # Grades
│   │   ├── attendance/              # Attendance tracking
│   │   ├── chat/                    # Messaging
│   │   ├── notification/            # Notifications
│   │   └── upload/                  # File uploads
│   ├── app/                         # Application initialization
│   ├── repository/                  # Data access layer (PostgreSQL)
│   ├── usecase/                     # Business logic layer
│   ├── delivery/http/               # HTTP handlers
│   └── shared/                      # Shared utilities
├── migrations/                      # Database migrations
├── tests/
│   ├── unit/                        # Unit tests
│   └── integration/                 # Integration tests
├── Dockerfile                       # Docker image build
├── docker-compose.yml              # Local dev environment
├── Makefile                        # Build commands
└── go.mod
```

## Clean Architecture Layers

1. **Domain Layer** (internal/domain/*)
   - Pure business logic entities
   - Repository interfaces
   - No external dependencies

2. **Repository Layer** (internal/repository/*)
   - Data access implementations (PostgreSQL)
   - Query builders
   - Transaction handling

3. **Usecase Layer** (internal/usecase/*)
   - Application business rules
   - Service implementations
   - Orchestrates repositories

4. **Delivery Layer** (internal/delivery/http/)
   - HTTP handlers
   - Request validation
   - Response formatting

5. **App Layer** (internal/app/)
   - Application initialization
   - Background workers (goroutines)
   - Graceful shutdown

## Core Services

10 Core domains:
- **User**: Authentication, user management
- **Course**: Course creation and management
- **Enrollment**: Student enrollments
- **Assignment**: Assignment management
- **Submission**: Student submissions
- **Grade**: Grading system
- **Attendance**: Attendance tracking
- **Chat**: Course messaging
- **Notification**: User notifications with background processing
- **Upload**: File storage

## Concurrency & Background Processing

The app includes goroutine-based background workers:

- **Notification Worker**: Processes notifications asynchronously via channel
- **Health Check Worker**: Monitors database connection every 30 seconds
- **Graceful Shutdown**: Waits for all background workers to complete

```go
// Send notification asynchronously
app.SendNotification(userID, message, "assignment_submission")
```

## Quick Start

### Local Development

```bash
# Install dependencies
go mod download

# Run locally (requires PostgreSQL)
make run

# Run tests
make test

# Build
make build
```

### Docker

```bash
# Start all services
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

### Database Migrations

```bash
# Apply migrations
make migrate-up

# Rollback
make migrate-down
```

## API Endpoints

All endpoints under /api/v1/

### Users
- POST /users - Register
- GET /users - List
- GET /users/{id} - Get
- PATCH /users/{id} - Update
- DELETE /users/{id} - Delete

### Courses
- POST /courses - Create
- GET /courses - List
- GET /courses/{id} - Get
- PATCH /courses/{id} - Update
- DELETE /courses/{id} - Delete

### Enrollments
- POST /enrollments - Enroll
- GET /enrollments/course/{courseID} - By course
- GET /enrollments/student/{studentID} - By student
- DELETE /enrollments/{id} - Remove

Similar patterns for assignments, submissions, grades, attendance, chat, notifications, uploads

## Environment Variables

```
DATABASE_URL=postgres://user:pass@localhost:5432/moodle
REDIS_URL=redis://localhost:6379
PORT=8080
```

## Technologies

- **Framework**: Gin Web Framework
- **Database**: PostgreSQL with pgx
- **Caching**: Redis
- **Language**: Go 1.21+
- **Concurrency**: Goroutines with channels

## Testing

```bash
# Run unit tests
go test ./tests/unit/...

# Run all tests with coverage
go test -v ./... -cover
```

## Performance

- Multi-stage Docker build (120MB final image)
- Connection pooling
- Paginated queries
- Indexed database columns
- Asynchronous notification processing
- Health checks on background worker

## License

Private - Educational Use Only
