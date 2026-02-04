# Routes Package - Clean Architecture

## Overview

The `routes` package is organized according to **Clean Architecture** principles. Each router is responsible for its domain and encapsulates routing logic.

## Structure

```
routes/
├── public.go          # Public routes (registration, login)
├── user.go           # User routes
├── course.go         # Course routes
├── enrollment.go     # Course enrollment routes
├── assignment.go     # Assignment routes
├── submission.go     # Submission routes
├── grade.go          # Grade routes
├── attendance.go     # Attendance routes
├── chat.go           # Chat routes
├── notification.go   # Notification routes
├── schedule.go       # Schedule routes
├── analytics.go      # Analytics routes
├── session.go        # Session routes
├── student.go        # Student routes
├── teacher.go        # Teacher routes
├── manager.go        # Manager routes
├── plagiarism.go     # Plagiarism check routes
├── upload.go         # File upload routes
├── group.go          # Group routes
└── admin.go          # Reserved for future admin routes
```

## Principles

### 1. Encapsulation
Each router encapsulates:
- Constructor (`New*Router`)
- Route setup method (`SetupRoutes`)
- Required dependencies (handlers)

### 2. Separation of Concerns
- **Router** - routing and middleware only
- **Handler** - HTTP request processing
- **Service** - business logic
- **Repository** - data access

### 3. Independence
Each router can be:
- Tested independently
- Modified without affecting other routers
- Reused

## Usage Example

### Creating a New Router

```go
// routes/example.go
package routes

import (
	"ap1-final-mini-moodle/internal/domain/example"
	"ap1-final-mini-moodle/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

type ExampleRouter struct {
	exampleHandler *example.Handler
}

func NewExampleRouter(exampleHandler *example.Handler) *ExampleRouter {
	return &ExampleRouter{
		exampleHandler: exampleHandler,
	}
}

func (er *ExampleRouter) SetupRoutes(api *gin.RouterGroup) {
	examples := api.Group("/examples")
	{
		// Public routes
		examples.GET("", er.exampleHandler.List)
		examples.GET("/:id", er.exampleHandler.GetByID)

		// Protected routes (admin only)
		adminRoutes := examples.Group("")
		adminRoutes.Use(middleware.RequireRole("admin"))
		{
			adminRoutes.POST("", er.exampleHandler.Create)
			adminRoutes.PUT("/:id", er.exampleHandler.Update)
			adminRoutes.DELETE("/:id", er.exampleHandler.Delete)
		}
	}
}
```

### Integration in router.go

```go
// In the main Router's SetupRoutes method
exampleRouter := routes.NewExampleRouter(r.handlers.Example)
exampleRouter.SetupRoutes(protected)
```

## Middleware

### Authentication
```go
protected := api.Group("")
protected.Use(r.auth.Authenticate())
```

### Role-based Authorization
```go
adminRoutes := group.Group("")
adminRoutes.Use(middleware.RequireRole("admin"))
```

Available roles:
- `admin` - System administrator
- `manager` - Manager
- `teacher` - Teacher
- `student` - Student

## Special Routes

### WebSocket (chat.go)
```go
// Separate function for WebSocket setup
func SetupChatWebSocket(engine *gin.Engine, auth *middleware.AuthMiddleware, 
                        service chat.Service, hub *websocket.Hub) {
	// WebSocket requires direct access to engine
}
```

## Testing

Example unit test for a router:

```go
func TestUserRouter_SetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	
	mockHandler := &user.Handler{}
	router := routes.NewUserRouter(mockHandler)
	
	api := engine.Group("/api/v1")
	router.SetupRoutes(api)
	
	// Check registered routes
	routes := engine.Routes()
	assert.Contains(t, routes, "GET /api/v1/users/me")
}
```

## Migration

When adding new routes:

1. Create a new file in `routes/` or update an existing one
2. Follow the pattern: Router struct → Constructor → SetupRoutes
3. Add handler to `api.Handlers` if a new one is needed
4. Initialize the router in `router.go`
5. Call `SetupRoutes()` in the appropriate place

## Benefits of Current Architecture

**Modularity** - each router is independent
**Testability** - easy to write unit tests
**Scalability** - simple to add new routers
**Readability** - clear file structure
**Maintainability** - changes are localized
**Clean Architecture** - separation of concerns

## Recommendations

1. **One router = one domain** - do not mix different domain areas
2. **Minimal logic** - routers only register routes
3. **Use middleware** - for common tasks (auth, logging, etc.)
4. **Document** - comment complex routes
5. **Group** - use router groups for organization
