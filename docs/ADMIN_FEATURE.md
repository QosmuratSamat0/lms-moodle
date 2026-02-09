# Admin Feature Implementation Summary

## Overview
Complete implementation of the Admin feature with full CRUD operations, role-based access control, and permission management.

## Database Schema
The `admins` table includes:
- `id` (UUID) - Primary key
- `user_id` (UUID FK) - Reference to users table
- `employee_id` (VARCHAR 50) - Unique identifier
- `department` (VARCHAR 255)
- `access_level` (VARCHAR 50) - "super_admin" or "admin"
- `permissions` (JSONB) - Array of permission strings
- `is_active` (BOOLEAN)
- `created_at` (TIMESTAMPTZ)
- `updated_at` (TIMESTAMPTZ)

## Files Created/Modified

### Domain Layer
- **[internal/domain/admin/entity.go](internal/domain/admin/entity.go)** - Domain entities
  - `Admin` struct with all fields
  - `CreateAdminInput` for creation
  - `UpdateAdminInput` for partial updates
  - `AdminFilter` for listing with filters
  
- **[internal/domain/admin/repository.go](internal/domain/admin/repository.go)** - Interface definitions
  - `Create`, `GetByID`, `GetByUserID`, `GetByEmployeeID`
  - `List` with filtering
  - `Update`, `Delete`

### Repository Layer
- **[internal/repository/admin/postgres.go](internal/repository/admin/postgres.go)** - PostgreSQL implementation
  - Full CRUD operations
  - Filtering by department, access_level, and active status
  - JOIN with users table for email and names
  - Proper JSON marshaling/unmarshaling for permissions
  - Pagination support

### Use Case Layer
- **[internal/usecase/admin/service.go](internal/usecase/admin/service.go)** - Business logic
  - `CreateAdmin` - Validates access levels and checks duplicates
  - `GetByID`, `GetByUserID` - Retrieves admin records
  - `List` - Lists admins with pagination
  - `UpdateAdmin` - Updates admin fields with validation
  - `DeleteAdmin` - Deletes admin record
  - Input validation and error handling

### Delivery Layer (HTTP)
- **[internal/delivery/http/admin.go](internal/delivery/http/admin.go)** - HTTP handlers
  - `Create` - POST /api/v1/admins
  - `GetByID` - GET /api/v1/admins/:id
  - `GetByUserID` - GET /api/v1/admins/user/:userID
  - `GetMyProfile` - GET /api/v1/admins/me
  - `List` - GET /api/v1/admins?limit=20&offset=0
  - `Update` - PUT /api/v1/admins/:id
  - `Delete` - DELETE /api/v1/admins/:id

- **[internal/delivery/http/admin_router.go](internal/delivery/http/admin_router.go)** - Route registration
  - Registers all admin endpoints
  - Groups under /api/v1/admins

- **[internal/delivery/http/middleware.go](internal/delivery/http/middleware.go)** - Authorization & authentication
  - `RequireRole(allowedRoles ...)` - Checks user role
  - `RequirePermission(permission)` - Checks specific permission
  - `RequireOwnership(resourceType)` - Checks resource ownership
  - `AuthMiddleware()` - JWT validation (placeholder)
  - `ContextMiddleware()` - Sets user context

### Dependency Injection
- **[internal/app/deps.go](internal/app/deps.go)** - Updated
  - Added `adminRepo` import
  - Added `adminUC` import
  - Added `AdminSvc` to `Deps` struct
  - Added admin repository initialization in `BuildDeps`
  - Added admin service initialization in `BuildDeps`
  - Added admin handler in `BuildHTTPModules`
  - Added admin module in `BuildHTTPModules`

## Features

### Admin Roles & Access Levels
- **super_admin** - Full system access, manage all resources
- **admin** - Limited admin access, manage assigned resources

### Permissions
- Admins can have custom permission arrays (JSON)
- Wildcard "*" grants all permissions
- Examples: "manage_users", "manage_courses", "manage_enrollments"

### Admin Capabilities
- Manage all users
- Manage all courses and enrollments
- Create and manage teachers and managers
- View system statistics
- Modify system configuration

### API Endpoints
```
POST   /api/v1/admins              - Create new admin
GET    /api/v1/admins              - List all admins (with pagination)
GET    /api/v1/admins/:id          - Get admin by ID
GET    /api/v1/admins/me           - Get current admin profile
PUT    /api/v1/admins/:id          - Update admin
DELETE /api/v1/admins/:id          - Delete admin
GET    /api/v1/admins/user/:userID - Get admin by user ID
```

## Middleware Usage

### Usage Examples
```go
// Require specific role
r.Use(RequireRole("admin", "super_admin"))

// Require specific permission
r.Use(RequirePermission("manage_users"))

// Require resource ownership
r.Use(RequireOwnership("course"))
```

## Request/Response Examples

### Create Admin
```json
POST /api/v1/admins
{
  "user_id": "uuid",
  "employee_id": "ADM001",
  "department": "IT",
  "access_level": "admin",
  "permissions": ["manage_users", "manage_courses"]
}
```

### Update Admin
```json
PUT /api/v1/admins/:id
{
  "department": "System Administration",
  "access_level": "super_admin",
  "permissions": ["*"],
  "is_active": true
}
```

### Response
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "employee_id": "ADM001",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "department": "IT",
  "access_level": "admin",
  "permissions": ["manage_users", "manage_courses"],
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## Validation Rules
- UserID: Required
- EmployeeID: Required, must be unique
- AccessLevel: Required, must be "super_admin" or "admin"
- Department: Optional
- Permissions: Optional array of strings

## Error Handling
- `ErrMissingRequired` - Missing required fields
- `ErrAlreadyExists` - Employee ID or User already has admin profile
- `ErrInvalidID` - Invalid ID format
- `ErrAdminNotFound` - Admin not found
- `ErrInvalidInput` - Invalid access level or other input

## Integration
All components are integrated into the dependency injection system:
1. Repository is initialized with database pool
2. Service wraps repository with business logic
3. Handler wraps service for HTTP endpoints
4. Module registers routes in Gin engine
5. Middleware provides authorization layer


