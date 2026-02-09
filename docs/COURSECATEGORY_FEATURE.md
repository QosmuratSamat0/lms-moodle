# Course Category Feature

## Overview

The Course Category feature provides the hierarchical organizational structure for grouping courses into logical categories. This enables better course discovery, organization, and permission-based access control through category managers.

## Purpose

- **Course Organization**: Organize courses into meaningful categories (e.g., "Web Development", "Data Science")
- **Structured Discovery**: Help students and teachers discover courses by category
- **Permission Scaffolding**: Foundation for category-based access control via Category Managers
- **Categorization Metadata**: Store category details including icons and display order
- **Aggregate Statistics**: Track course counts per category

## Architecture

```
Domain Layer
├── Entity: CourseCategory
│   ├── ID (UUID)
│   ├── Name (string, required)
│   ├── Description (optional)
│   ├── Icon (optional, for UI)
│   ├── Order (int, for sorting)
│   ├── IsActive (boolean)
│   ├── TotalCourses (computed field from aggregation)
│   └── Timestamps

Repository Interface
└── Defines contract for data operations

Repository Layer
└── PostgreSQL Implementation
    ├── Full CRUD operations
    ├── LEFT JOIN courses with GROUP BY COUNT aggregation
    ├── Ordering by display order and creation date
    ├── Pagination support
    ├── Active category filtering

UseCase Layer
├── Service: CourseCategory
│   ├── Name validation (required, non-empty)
│   ├── Default IsActive status
│   ├── Business logic enforcement
│   └── Error handling

HTTP Delivery Layer
├── Handlers: 5 endpoints
├── Router: Route registration
└── Response mapping
```

## Database Schema

```sql
CREATE TABLE course_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    icon VARCHAR(100),
    "order" INT DEFAULT 0 NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Relationship Table
CREATE TABLE courses (
    ...
    category_id UUID REFERENCES course_categories(id) ON DELETE SET NULL,
    ...
);
```

## API Endpoints

### 1. Create Course Category
```
POST /api/v1/categories
Content-Type: application/json
Authorization: Bearer <token>

{
    "name": "Web Development",
    "description": "All web development related courses",
    "icon": "code",
    "order": 1
}

Response: 201 Created
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Web Development",
    "description": "All web development related courses",
    "icon": "code",
    "order": 1,
    "is_active": true,
    "total_courses": 0,
    "created_at": "2024-01-09T10:00:00Z",
    "updated_at": "2024-01-09T10:00:00Z"
}
```

### 2. List All Course Categories
```
GET /api/v1/categories
Authorization: Bearer <token>

Query Parameters:
  - limit: int (default: 20, max: 100)
  - offset: int (default: 0)

Response: 200 OK
{
    "data": [
        {
            "id": "uuid",
            "name": "Web Development",
            "description": "...",
            "icon": "code",
            "order": 1,
            "is_active": true,
            "total_courses": 5,
            "created_at": "2024-01-09T10:00:00Z",
            "updated_at": "2024-01-09T10:00:00Z"
        },
        {
            "id": "uuid",
            "name": "Data Science",
            "description": "...",
            "icon": "chart",
            "order": 2,
            "is_active": true,
            "total_courses": 3,
            "created_at": "2024-01-09T10:00:00Z",
            "updated_at": "2024-01-09T10:00:00Z"
        }
    ],
    "limit": 20,
    "offset": 0
}
```

### 3. Get Course Category by ID
```
GET /api/v1/categories/:id
Authorization: Bearer <token>

Response: 200 OK
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Web Development",
    "description": "All web development related courses",
    "icon": "code",
    "order": 1,
    "is_active": true,
    "total_courses": 5,
    "created_at": "2024-01-09T10:00:00Z",
    "updated_at": "2024-01-09T10:00:00Z"
}
```

### 4. Update Course Category
```
PUT /api/v1/categories/:id
Authorization: Bearer <token>
Content-Type: application/json

{
    "name": "Web Development (Updated)",
    "description": "Updated description",
    "icon": "globe",
    "order": 2,
    "is_active": true
}

Response: 200 OK
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Web Development (Updated)",
    "description": "Updated description",
    "icon": "globe",
    "order": 2,
    "is_active": true,
    "total_courses": 5,
    "created_at": "2024-01-09T10:00:00Z",
    "updated_at": "2024-01-09T10:00:55Z"
}
```

### 5. Delete Course Category
```
DELETE /api/v1/categories/:id
Authorization: Bearer <token>

Response: 204 No Content
```

## Entity Definition

```go
type CourseCategory struct {
    ID           string    `json:"id"`
    Name         string    `json:"name"`
    Description  string    `json:"description,omitempty"`
    Icon         string    `json:"icon,omitempty"`
    Order        int       `json:"order"`
    IsActive     bool      `json:"is_active"`
    TotalCourses int       `json:"total_courses"` // Computed via COUNT(courses.id)
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

type CreateCourseCategoryInput struct {
    Name        string `json:"name" binding:"required"`
    Description string `json:"description"`
    Icon        string `json:"icon"`
    Order       int    `json:"order"`
}

type UpdateCourseCategoryInput struct {
    Name        *string `json:"name"`
    Description *string `json:"description"`
    Icon        *string `json:"icon"`
    Order       *int    `json:"order"`
    IsActive    *bool   `json:"is_active"`
}
```

## Error Handling

### Common Error Responses

| Error | Status | Message |
|-------|--------|---------|
| Category not found | 404 | "course category not found" |
| Missing category name | 400 | "missing required field" |
| Invalid ID format | 400 | "invalid ID format" |
| Unauthorized access | 401 | "unauthorized access" |
| Invalid input data | 400 | "invalid input data" |
| Category name exists | 409 | "category with this name already exists" |

## Validation Rules

1. **Name Required**: Category name must be non-empty string
2. **UUID Format**: ID must be valid UUID format
3. **Order (Optional)**: Integer for display ordering
4. **Icon (Optional)**: String representation of icon (e.g., "code", "chart")
5. **IsActive**: Boolean flag, defaults to true on creation

## Service Methods

### Create(ctx, input)
```go
Validates:
- input.Name is not empty (returns ErrMissingRequired if empty)
- Sets IsActive = true by default

Returns: CourseCategory with generated UUID and timestamps
```

### GetByID(ctx, id)
```go
Queries: course_categories LEFT JOIN courses with GROUP BY COUNT
Aggregates: Total number of courses in category
Returns: CourseCategory with TotalCourses populated or ErrCategoryNotFound
```

### List(ctx, filter)
```go
Query Behavior:
- Orders by "order" ASC, then created_at DESC
- Supports pagination (limit, offset)
- Includes course count aggregation

Returns: []*CourseCategory, total_count, error
```

### Update(ctx, id, input)
```go
Allows: Selective field updates
Updates:
- Name (if provided)
- Description (if provided)
- Icon (if provided)
- Order (if provided)
- IsActive (if provided)

Returns: Updated CourseCategory or ErrCategoryNotFound
```

### Delete(ctx, id)
```go
Behavior: Hard delete from database
Note: Cascading behavior on courses table is SET NULL

Returns: error if any
```

## Relationships

### With Courses (1-to-Many)
```
CourseCategory ──── 0..* Courses
```
- A course belongs to at most one category
- A category can have zero to many courses
- Deleting a category sets course category_id to NULL
- Total course count available via aggregation query

### With Category Managers (1-to-Many)
```
CourseCategory ──── 0..* CategoryManagers
```
- Category managers control access to specific categories
- Enables fine-grained permission-based access control
- CategoryManager tracks permission levels per user per category

## Display Ordering

Categories are ordered by the `order` field (ascending) with creation date as tiebreaker:

```
ORDER BY "order" ASC, created_at DESC
```

This allows:
- Admin control over category display order via the `order` field
- Tie-breaking by creation date for consistent pagination
- Supporting featured/priority categories via lower order values

Example Display:
```
1. [order=0] Web Development (5 courses)
2. [order=1] Data Science (3 courses)
3. [order=2] Mobile Development (2 courses)
4. [order=3] DevOps (1 course)
```

## Integration Example

```go
// Create a category
category, err := coursecategoryService.Create(ctx, &CreateCourseCategoryInput{
    Name:        "Web Development",
    Description: "Frontend and backend development",
    Icon:        "code",
    Order:       1,
})

// Get all categories with course counts
categories, total, err := coursecategoryService.List(ctx, 20, 0)

// For each category, retrieve its category managers
managers, _, err := categorymanagerService.GetByCategoryID(ctx, category.ID, ...)

// Check user's permission level for a category
manager, err := categorymanagerService.GetByUserAndCategory(ctx, userID, categoryID)
if manager.PermissionLevel == "view" {
    // Read-only access
}
```

## Soft Delete Consideration

The current implementation uses **hard delete**. To implement soft delete (for audit trails), add:

```go
type CourseCategory struct {
    // ... existing fields
    DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// Update queries to filter: WHERE is_active = true AND deleted_at IS NULL
```

## Testing

### Unit Tests Needed
- [ ] Validate name is required
- [ ] Set IsActive=true on creation
- [ ] Update only specified fields
- [ ] Handle missing fields gracefully

### Integration Tests Needed
- [ ] Create, retrieve, update, delete categories
- [ ] Verify course count aggregation works
- [ ] Test pagination and ordering
- [ ] Cascade behavior with category managers
- [ ] Soft delete behavior (if implemented)

## Performance Considerations

1. **Aggregation Query**: LEFT JOIN with GROUP BY COUNT
   - Index on courses.category_id recommended
   - Count operation may be slow with many courses
   - Consider caching total_courses for high-traffic categories

2. **Pagination**: 
   - Limit: 20-100 recommended
   - Offset pagination may be slow for large datasets
   - Consider cursor-based pagination for 10k+ records

3. **Ordering**:
   - Current: ORDER BY "order" ASC, created_at DESC
   - Efficient with index on ("order", "created_at")

## Future Enhancements

1. **Category Hierarchies**: Support nested subcategories
2. **Category Icons Library**: Pre-defined icon set
3. **Category Analytics**: Track course enrollment by category
4. **Bulk Operations**: Create/update multiple categories
5. **Category Featured Status**: Promote categories to homepage
6. **Auto-Archival**: Archive empty categories after N days
7. **Category Templates**: Pre-configured category structures
