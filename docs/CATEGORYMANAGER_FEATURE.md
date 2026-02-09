# Category Manager Feature

## Overview

The Category Manager feature provides fine-grained permission management for course categories. It allows assigning users specific permission levels over individual categories, enabling a flexible hierarchical access control system.

## Purpose

- **Granular Category Access Control**: Assign different permission levels per user per category
- **Delegation of Authority**: Managers can delegate specific category management responsibilities
- **Permission Levels**: Three-tier permission system (view, edit, admin)
- **Audit Trail**: Track who has access to what and at which permission level

## Architecture

```
Domain Layer
├── Entity: CategoryManager
│   ├── ID (UUID)
│   ├── UserID (FK to users)
│   ├── CategoryID (FK to course_categories)
│   ├── PermissionLevel ("view" | "edit" | "admin")
│   ├── IsActive (boolean)
│   └── Timestamps
├── Repository Interface
│   └── Defines contract for data operations

Repository Layer
└── PostgreSQL Implementation
    ├── Full CRUD operations
    ├── Query by UserAndCategory (for uniqueness checks)
    ├── Nested data retrieval with category details
    ├── Count aggregation for dashboards

UseCase Layer
├── Service: CategoryManager
│   ├── Permission level validation
│   ├── Duplicate relationship prevention
│   ├── Ownership verification
│   └── Business logic enforcement

HTTP Delivery Layer
├── Handlers: 7 endpoints
├── Router: Route registration
└── Middleware: Authorization checks
```

## Database Schema

```sql
CREATE TABLE category_managers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES course_categories(id) ON DELETE CASCADE,
    permission_level VARCHAR(50) NOT NULL CHECK (permission_level IN ('view', 'edit', 'admin')),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, category_id)
);
```

## API Endpoints

### 1. Create Category Manager Assignment
```
POST /api/v1/category-managers
Content-Type: application/json

{
    "user_id": "uuid",
    "category_id": "uuid",
    "permission_level": "edit"
}

Response: 201 Created
{
    "id": "uuid",
    "user_id": "uuid",
    "category_id": "uuid",
    "permission_level": "edit",
    "is_active": true,
    "created_at": "2024-01-09T10:00:00Z",
    "updated_at": "2024-01-09T10:00:00Z"
}
```

### 2. Get Category Manager by ID
```
GET /api/v1/category-managers/:id
Authorization: Bearer <token>

Response: 200 OK
{
    "id": "uuid",
    "user_id": "uuid",
    "category_id": "uuid",
    "permission_level": "edit",
    "is_active": true,
    "created_at": "2024-01-09T10:00:00Z",
    "updated_at": "2024-01-09T10:00:00Z"
}
```

### 3. Get Category Manager by User and Category
```
GET /api/v1/category-managers/by-relationship?user_id=<uuid>&category_id=<uuid>
Authorization: Bearer <token>

Response: 200 OK
{
    "id": "uuid",
    "user_id": "uuid",
    "category_id": "uuid",
    "permission_level": "edit",
    ...
}
```

### 4. Get All Category Managers for User
```
GET /api/v1/category-managers/user/:user_id
Authorization: Bearer <token>

Query Parameters:
  - limit: int (default: 20, max: 100)
  - offset: int (default: 0)

Response: 200 OK
{
    "data": [
        {
            "id": "uuid",
            "user_id": "uuid",
            "category_id": "uuid",
            "permission_level": "edit",
            "is_active": true,
            ...
        }
    ],
    "limit": 20,
    "offset": 0
}
```

### 5. Get All Category Managers for Category
```
GET /api/v1/category-managers/category/:category_id
Authorization: Bearer <token>

Query Parameters:
  - limit: int (default: 20, max: 100)
  - offset: int (default: 0)

Response: 200 OK
{
    "data": [
        {
            "id": "uuid",
            "user_id": "uuid",
            "category_id": "uuid",
            "permission_level": "view",
            ...
        }
    ],
    "limit": 20,
    "offset": 0
}
```

### 6. List All Category Manager Assignments
```
GET /api/v1/category-managers
Authorization: Bearer <token>

Query Parameters:
  - limit: int (default: 20, max: 100)
  - offset: int (default: 0)

Response: 200 OK
{
    "data": [
        {...}
    ],
    "limit": 20,
    "offset": 0
}
```

### 7. Update Category Manager Assignment
```
PUT /api/v1/category-managers/:id
Authorization: Bearer <token>
Content-Type: application/json

{
    "permission_level": "admin",
    "is_active": true
}

Response: 200 OK
{
    "id": "uuid",
    "permission_level": "admin",
    ...
}
```

### 8. Delete Category Manager Assignment
```
DELETE /api/v1/category-managers/:id
Authorization: Bearer <token>

Response: 204 No Content
```

### 9. Get Category Manager with Details
```
GET /api/v1/category-managers/:id/details
Authorization: Bearer <token>

Response: 200 OK
{
    "id": "uuid",
    "user_id": "uuid",
    "category_id": "uuid",
    "permission_level": "edit",
    "is_active": true,
    "category_info": {
        "id": "uuid",
        "name": "Web Development",
        "description": "Web dev courses",
        "icon": "code",
        "order": 1,
        "is_active": true
    },
    "created_at": "2024-01-09T10:00:00Z",
    "updated_at": "2024-01-09T10:00:00Z"
}
```

## Permission Levels

### "view" - Read-Only Access
- Can view category and its courses
- Cannot modify category settings
- Can view assigned courses

### "edit" - Read/Write Access
- Can view category and its courses
- Can modify category details
- Can manage category assignments
- Cannot delete the category

### "admin" - Full Administrative Access
- Can view and modify all category settings
- Can manage category managers
- Can delete the category
- Full control over category resources

## Error Handling

### Common Error Responses

| Error | Status | Message |
|-------|--------|---------|
| Invalid permission_level | 400 | "invalid input data" |
| User/category not found | 404 | "category manager not found" |
| Duplicate relationship | 409 | "record already exists" |
| Missing required fields | 400 | "missing required field" |
| Invalid ID format | 400 | "invalid ID format" |
| Unauthorized access | 401 | "unauthorized access" |
| Insufficient permissions | 403 | "access forbidden" |

## Validation Rules

1. **Permission Level Validation**: Must be one of: "view", "edit", "admin"
2. **User Uniqueness**: One user per category max (Unique constraint on user_id, category_id)
3. **Required Fields**: user_id, category_id, permission_level must be provided
4. **UUID Validation**: user_id and category_id must be valid UUIDs

## Service Methods

### Create(ctx, input)
```go
Validates:
- input.PermissionLevel in [view, edit, admin]
- No existing relationship between user and category
- User and category exist

Returns: CategoryManager or error
```

### GetByID(ctx, id)
```go
Queries: category_managers table by ID
Returns: CategoryManager or ErrCategoryManagerNotFound
```

### GetByUserAndCategory(ctx, userID, categoryID)
```go
Used for: Duplicate checking, internal relationship lookups
Returns: CategoryManager or error
```

### GetWithDetails(ctx, id)
```go
Queries: WITH nested category information
Returns: CategoryManager with CategoryInfo populated
```

### List(ctx, filter)
```go
Supports: Pagination, filtering by user_id and category_id
Returns: []*CategoryManager, total_count, error
```

### Update(ctx, id, input)
```go
Allows: Update permission_level and is_active status
Validates: New permission_level if provided
Returns: Updated CategoryManager or error
```

### Delete(ctx, id)
```go
Soft delete via UPDATE is_active = false
Returns: error if any
```

## Integration Example

```go
// In your handler
categoryManagerService.Create(ctx, &CategoryManagerInput{
    UserID:          userID,
    CategoryID:      categoryID,
    PermissionLevel: "edit",
})

// Later, check permissions
manager, err := categoryManagerService.GetByUserAndCategory(ctx, userID, categoryID)
if err != nil || manager.PermissionLevel == "view" {
    // Deny write operations
}
```

## Testing

### Unit Tests Needed
- [ ] Validate permission level enum
- [ ] Prevent duplicate user-category relationships
- [ ] Update only specific fields
- [ ] Cascade deletions from user/category

### Integration Tests Needed
- [ ] Create and retrieve category managers
- [ ] List with pagination
- [ ] Update multiple managers
- [ ] Permission level enforcement

## Audit Considerations

- All operations logged with user context
- Track permission level changes for compliance
- Maintain historical record of access changes
- Integration with audit_logs table recommended

## Future Enhancements

1. **Batch Assignment**: Assign multiple users to category at once
2. **Permission Templates**: Pre-defined permission packages
3. **Time-Bound Access**: Temporary category manager assignments
4. **Audit Reports**: Generate category access reports
5. **Delegation Chains**: Track permission delegation history
