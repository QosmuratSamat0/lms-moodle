# aLMS Test Infrastructure

This directory contains all testing infrastructure and tests for the aLMS backend.

## Directory Structure

```
tests/
├── testutil/          # Shared test utilities
│   ├── config.go      # Test configuration (env vars)
│   ├── database.go    # Database helpers (connect, truncate)
│   ├── tx.go          # Transaction-based test isolation
│   ├── migrate.go     # Migration helpers
│   ├── clock.go       # Time mocking (Clock interface)
│   ├── mocks.go       # Common mocks and fixed UUIDs
│   ├── http.go        # HTTP testing helpers (Gin)
│   └── fixtures.go    # Test data seeding
├── integration/       # Integration tests (repository layer)
│   ├── setup_test.go  # Test setup/teardown
│   └── *_test.go      # Repository integration tests
└── e2e/               # End-to-end API tests
    ├── setup_test.go  # E2E test setup
    └── *_test.go      # API endpoint tests
```

## Quick Start

### 1. Start Test Infrastructure

```bash
# Start test postgres and redis containers
make test-docker-up

# Run migrations against test database
make test-migrate
```

### 2. Run Tests

```bash
# Run unit tests only (no external dependencies)
make test-unit

# Run integration tests (requires test db)
make test-integration

# Run e2e tests (requires running API server)
make test-e2e

# Run all tests
make test-all
```

### 3. Cleanup

```bash
# Stop test containers and remove volumes
make test-docker-down

# Full cleanup including coverage files
make test-clean
```

## Test Types

### Unit Tests (Default)
- Location: `internal/domain/<module>/*_test.go`
- Run with: `go test ./internal/...`
- No build tags required
- Mock all external dependencies

### Integration Tests
- Location: `tests/integration/*_test.go`
- Build tag: `//go:build integration`
- Run with: `go test -tags=integration ./tests/integration/...`
- Test repository layer against real PostgreSQL
- Use transaction rollback for test isolation

### E2E Tests
- Location: `tests/e2e/*_test.go`
- Build tag: `//go:build e2e`
- Run with: `go test -tags=e2e ./tests/e2e/...`
- Test full API through HTTP calls
- Require running API server

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `TEST_DATABASE_URL` | `postgres://alms_test:alms_test_secret@localhost:5433/alms_test_db?sslmode=disable` | Test PostgreSQL connection string |
| `TEST_REDIS_URL` | `redis://localhost:6380/0` | Test Redis connection string |
| `TEST_API_URL` | `http://localhost:8080` | API URL for e2e tests |
| `TEST_TIMEOUT` | `30` | Test timeout in seconds |
| `MIGRATIONS_DIR` | `./migrations` | Path to migrations directory |

## Writing Tests

### Service Tests (Default Pattern)

```go
package course

import (
    "context"
    "testing"
    "github.com/google/uuid"
)

// Mock repository
type mockCourseRepo struct {
    GetByIDFn func(ctx context.Context, id uuid.UUID) (*CourseWithTeacher, error)
    CreateFn  func(ctx context.Context, course *Course) error
}

func (m *mockCourseRepo) GetByID(ctx context.Context, id uuid.UUID) (*CourseWithTeacher, error) {
    return m.GetByIDFn(ctx, id)
}

func (m *mockCourseRepo) Create(ctx context.Context, course *Course) error {
    return m.CreateFn(ctx, course)
}

func TestService_GetByID(t *testing.T) {
    tests := []struct {
        name      string
        courseID  uuid.UUID
        mockSetup func(*mockCourseRepo)
        wantErr   bool
    }{
        {
            name:     "success",
            courseID: uuid.MustParse("..."),
            mockSetup: func(m *mockCourseRepo) {
                m.GetByIDFn = func(ctx context.Context, id uuid.UUID) (*CourseWithTeacher, error) {
                    return &CourseWithTeacher{...}, nil
                }
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := &mockCourseRepo{}
            tt.mockSetup(repo)
            svc := NewService(repo)

            _, err := svc.GetByID(context.Background(), tt.courseID)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Handler Tests Pattern

```go
func TestHandler_Create(t *testing.T) {
    router := testutil.TestRouter()

    // Create mock service
    mockSvc := &mockCourseService{
        CreateFn: func(ctx context.Context, teacherID uuid.UUID, req *CreateCourseRequest) (*CourseResponse, error) {
            return &CourseResponse{ID: "123"}, nil
        },
    }

    handler := NewHandler(mockSvc)

    // Set up route with auth context
    router.POST("/courses",
        testutil.SetAuthContext(testutil.TestTeacherID, "teacher@test.com", "teacher"),
        handler.Create,
    )

    // Test
    w := testutil.PerformRequest(router, "POST", "/courses", map[string]string{
        "title": "New Course",
    })

    testutil.AssertStatusCode(t, w, http.StatusCreated)
}
```

### Integration Test Pattern

```go
//go:build integration

func TestCourseRepository_Create(t *testing.T) {
    pool := getTestPool(t)
    fixtures := testutil.NewFixtures(pool)
    txTest := testutil.NewTxTest(pool)
    ctx := context.Background()

    // Seed dependencies
    fixtures.CreateUser(ctx, t, testutil.TestTeacherID, "teacher@test.com", "hash", "teacher")
    fixtures.CreateTeacher(ctx, t, testutil.TestTeacherID, "Test", "Teacher", "CS")

    repo := course.NewRepository(pool)

    // Use transaction rollback for isolation
    txTest.Run(t, "create course", func(t *testing.T, tx pgx.Tx) {
        c := &course.Course{Title: "Test", IsActive: true}
        err := repo.Create(ctx, c)
        if err != nil {
            t.Fatalf("Create() error = %v", err)
        }
    })
}
```

## Fixed Test UUIDs

Use these fixed UUIDs from `testutil/mocks.go` for deterministic tests:

```go
testutil.TestAdminID     // 00000000-0000-0000-0000-000000000010
testutil.TestTeacherID   // 00000000-0000-0000-0000-000000000011
testutil.TestStudentID   // 00000000-0000-0000-0000-000000000021
testutil.TestCourseID1   // 00000000-0000-0000-0001-000000000001
// ... etc
```

## Time Mocking

Use `testutil.Clock` interface for time-dependent tests:

```go
clock := testutil.NewFixedClockAt(2026, 1, 25, 10, 0, 0)
// Use clock.Now() instead of time.Now()
clock.Advance(time.Hour) // Move time forward
```
