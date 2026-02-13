# Mini Moodle LMS API

A production-ready Learning Management System API built with Go, featuring Clean Architecture, concurrent processing, comprehensive REST endpoints, and PostgreSQL database integration.

## Overview

Mini Moodle is a backend API for educational institutions to manage courses, students, teachers, assignments, grades, attendance, and more. The system implements clean architecture principles with clear separation of concerns across domain, repository, usecase, and delivery layers.

## Core Features

- 10+ domain services covering all LMS functionality
- User authentication and role-based access control (Admin, Teacher, Manager, Student)
- Course management with teacher assignments
- Student enrollments with grade tracking
- Assignment submissions and automated grading
- Attendance tracking
- Real-time chat and notifications
- File uploads and document management
- Group management for collaborative learning
- Background worker processing for async operations
- Graceful shutdown with worker cleanup
- Comprehensive error handling and validation
- Request logging and database health monitoring

## Technology Stack

- **Language**: Go 1.21+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL 15+ with pgx driver and connection pooling
- **Concurrency**: Goroutines with channel-based communication
- **Containerization**: Docker and docker-compose
- **Documentation**: Swagger/OpenAPI specification
- **Migrations**: SQL-based database versioning

## Project Structure

```
internal/
├── app/                          # Application initialization and background workers
│   ├── app.go                    # Main app setup
│   ├── deps.go                   # Dependency injection
│   └── workers.go                # Background job processing
├── domain/                       # Business entities and interfaces
│   ├── user/                     # User/authentication domain
│   ├── student/                  # Student domain with group support
│   ├── teacher/                  # Teacher domain
│   ├── course/                   # Course management
│   ├── assignment/               # Assignment domain
│   ├── submission/               # Submission domain
│   ├── grade/                    # Grade/assessment domain
│   ├── attendance/               # Attendance tracking
│   ├── chat/                     # Messaging domain
│   ├── notification/             # Notification domain
│   ├── group/                    # Group management
│   └── ...                       # Additional domains
├── repository/                   # Data access layer (PostgreSQL)
│   ├── user/                     # User queries
│   ├── student/                  # Student queries with group joins
│   ├── course/                   # Course queries
│   └── ...                       # Repository implementations
├── usecase/                      # Business logic and services
│   ├── user/                     # User service logic
│   ├── student/                  # Student service logic
│   ├── course/                   # Course service logic
│   └── ...                       # Service implementations
├── delivery/http/                # HTTP API handlers
│   ├── student.go                # Student endpoints
│   ├── course.go                 # Course endpoints
│   ├── user.go                   # User endpoints
│   └── ...                       # HTTP handlers
└── shared/                       # Shared utilities
    ├── auth/                     # JWT authentication
    ├── config/                   # Configuration
    ├── database/                 # Database setup
    ├── errors/                   # Error definitions
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+ 
- Docker and Docker Compose (for containerized setup)
- Make tool

## Installation

### Local Development

1. Clone the repository:
```bash
cd ap1-final-mini-moodle
```

2. Install Go dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/moodle"
export PORT=8080
```

4. Run database migrations:
```bash
make migrate-up
```

5. Start the API server:
```bash
make run
```

The API will be available at `http://localhost:8080/api/v1`

### Docker Setup

For complete containerized environment with PostgreSQL and Redis:

```bash
# Start all services
make docker-up

# Apply database migrations
make migrate-up

# View logs
docker-compose logs -f api

# Stop services
make docker-down
```

## Build Commands

```bash
make build              # Compile binary to bin/
make run                # Run API locally with live reload
make test               # Execute test suite
make clean              # Remove build artifacts
make docker-up          # Start Docker containers
make docker-down        # Stop containers  
make docker-logs        # Stream application logs
make migrate-up         # Apply all pending database migrations
make migrate-down       # Rollback last database migration
make docs               # Generate Swagger documentation
```

## API Documentation

### Authentication

All endpoints require JWT bearer token in Authorization header:

```
Authorization: Bearer <JWT_TOKEN>
```

Token includes user ID, email, role, and expiration timestamp.

### Core Endpoints

**Students**
- `GET /api/v1/students` - List all students with filtering
- `GET /api/v1/students/{id}` - Get student details
- `POST /api/v1/students` - Create new student
- `PATCH /api/v1/students/{id}` - Update student info
- `DELETE /api/v1/students/{id}` - Delete student record

**Courses**
- `GET /api/v1/courses` - List courses
- `GET /api/v1/courses/{id}` - Get course details
- `POST /api/v1/courses` - Create course
- `PATCH /api/v1/courses/{id}` - Update course
- `DELETE /api/v1/courses/{id}` - Delete course

**Enrollments**
- `GET /api/v1/enrollments` - List enrollments
- `POST /api/v1/enrollments` - Enroll student in course
- `GET /api/v1/enrollments/{id}` - Get enrollment details
- `DELETE /api/v1/enrollments/{id}` - Remove enrollment

**Assignments**
- `GET /api/v1/assignments` - List assignments
- `POST /api/v1/assignments` - Create assignment
- `PATCH /api/v1/assignments/{id}` - Update assignment
- `DELETE /api/v1/assignments/{id}` - Delete assignment

**Submissions**
- `GET /api/v1/submissions` - List submissions
- `POST /api/v1/submissions` - Submit assignment
- `GET /api/v1/submissions/{id}` - Get submission details

**Grades**
- `GET /api/v1/grades` - List grades
- `POST /api/v1/grades` - Record grade
- `PATCH /api/v1/grades/{id}` - Update grade

**Attendance**
- `GET /api/v1/attendance` - List attendance records
- `POST /api/v1/attendance` - Record attendance

**Chat**
- `GET /api/v1/messages` - List messages
- `POST /api/v1/messages` - Send message
- `WS /api/v1/chat` - WebSocket chat endpoint

**Groups**
- `GET /api/v1/groups` - List groups
- `POST /api/v1/groups` - Create group
- `PATCH /api/v1/groups/{id}` - Update group
- `POST /api/v1/groups/{id}/members` - Add member to group

Additional endpoints available for Teachers, Managers, Admins, Categories, and Notifications.

## Database Schema

### Key Tables

**students**
- id (UUID)
- user_id (FK -> users)
- student_code (unique)
- major (text)
- year (int: 1-4)
- gpa (numeric)
- enrollment_status (text)
- group_id (FK -> groups, nullable)
- admitted_at (timestamp)
- created_at (timestamp)
- updated_at (timestamp)

**courses**
- id (UUID)
- title (text)
- code (unique)
- description (text)
- teacher_id (FK -> users)
- category_id (FK -> course_categories)
- credits (int)
- created_at (timestamp)
- updated_at (timestamp)

**enrollments**
- id (UUID)
- student_id (FK -> users)
- course_id (FK -> courses)
- status (text)
- grade (numeric, nullable)
- enrolled_at (timestamp)

**assignments**
- id (UUID)
- course_id (FK -> courses)
- title (text)
- description (text)
- due_date (timestamp)
- max_score (numeric)
- created_at (timestamp)

**groups**
- id (UUID)
- name (text)
- course_id (FK -> courses)
- created_at (timestamp)

Migration files in `migrations/` directory implement database schema versioning.

## Concurrency Model

The application uses Go goroutines and channels for concurrent processing:

### Background Workers

- Notification Worker: Processes user notifications asynchronously
- Email Worker: Sends email notifications
- Upload Worker: Handles file uploads and processing
- Health Check Worker: Monitors database connectivity every 30 seconds

### Implementation

```go
// Workers process jobs from channels
app.NotificationChannel <- Notification{UserID, Message, Type}

// Main gracefully shuts down workers on signal
// Waits for all pending jobs to complete
```

## Error Handling

Standardized error responses across all endpoints:

```json
{
  "error": "Error message",
  "details": "Additional context"
}
```

HTTP status codes:
- 200: Success
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 500: Internal Server Error

## Configuration

### Environment Variables

Required:
- `DATABASE_URL` - PostgreSQL connection string
- `PORT` - API server port (default: 8080)

Optional:
- `REDIS_URL` - Redis connection string (if using caching)
- `JWT_SECRET` - Secret key for JWT signing
- `LOG_LEVEL` - Log verbosity level

### Database Connection

Uses pgx with connection pooling:
- Max connections: 25
- Min connections: 5
- Connection timeout: 30 seconds

## Testing

Run the test suite:

```bash
# All tests with verbose output
go test -v ./...

# Tests with coverage report
go test -v ./... -cover

# Specific package tests
go test -v ./internal/repository/student/...
```

Test files located in `tests/` directory with unit and integration tests.

## Performance Optimization

- Connection pooling reduces database overhead
- Indexed queries on frequently searched columns
- Asynchronous background processing prevents blocking
- Goroutine-based concurrency for parallel request handling
- Clean code separation enables efficient testing and refactoring

## Deployment

### Docker Image Build

Multi-stage Docker build produces optimized image (~120MB):

```bash
docker build -t mini-moodle:latest .
```

### Production Deployment

Use docker-compose for orchestration:

```bash
docker-compose -f docker-compose.yml up -d
```

For production, update configuration:
- Use strong JWT_SECRET
- Enable HTTPS
- Set proper CORS headers
- Use managed PostgreSQL service
- Enable request logging
- Set up monitoring and alerting

## Recent Updates

### Student Entity Enhancement

Added group support to the Student domain:
- `group_id` (UUID, nullable): References group assignment
- `group_name` (string): Populated via LEFT JOIN with groups table
- Database migration: ALTER TABLE students ADD COLUMN group_id UUID REFERENCES groups(id) ON DELETE SET NULL
- Repository queries updated with proper NULL handling using COALESCE and type casting

### SQL Fixes

Fixed critical NULL scanning errors across teacher, manager, and admin repositories:
- Added COALESCE() with type casting for UUID fields: `COALESCE(field::text, '')`
- Fixed malformed SQL with missing FROM clauses
- Added proper LEFT JOIN syntax for joined fields
- Verified all 3+ endpoints return 200 OK responses

## Known Limitations

- WebSocket chat has concurrent connection limits
- File uploads limited to 50MB per file
- Redis caching optional but recommended for production
- Some endpoints require specific user roles (admin, teacher, etc.)

## Troubleshooting

**Server won't start:**
- Check DATABASE_URL is correct
- Verify PostgreSQL is running
- Check port 8080 is not in use

**Database migration errors:**
- Ensure PostgreSQL version is 15+
- Check database user has creation privileges
- Run make migrate-up after database creation

**API returns 401 Unauthorized:**
- Check JWT token has not expired
- Verify Authorization header format: `Bearer <token>`
- Confirm token contains valid user_id

**Slow query responses:**
- Check database indexes are created
- Verify connection pooling settings
- Monitor database server resources

## Support

For issues or questions, refer to:
- API specification: docs/swagger.yaml
- Architecture details: docs/architecture.md
- Individual domain documentation in internal/domain/*/

## Project Status

Production-ready. All core features implemented, tested, and deployed successfully.

## License

Private - Educational use only.
