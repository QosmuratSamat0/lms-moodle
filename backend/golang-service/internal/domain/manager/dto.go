package manager

import "github.com/google/uuid"

// UpdateManagerRequest represents a request to update manager profile
type UpdateManagerRequest struct {
	FirstName *string `json:"first_name" validate:"omitempty,max=100"`
	LastName  *string `json:"last_name" validate:"omitempty,max=100"`
}

// ManagerResponse represents a manager in API responses
type ManagerResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	FirstName *string   `json:"first_name"`
	LastName  *string   `json:"last_name"`
	IsActive  bool      `json:"is_active"`
}

// ManagerListResponse represents a paginated list of managers
type ManagerListResponse struct {
	Managers   []ManagerResponse `json:"managers"`
	TotalCount int               `json:"total_count"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
}

// ManagerFilter represents filters for listing managers
type ManagerFilter struct {
	IsActive *bool   `form:"is_active"`
	Search   *string `form:"search"`
}

// UserManagementAction represents an action on a user account
type UserManagementAction struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Action string    `json:"action" validate:"required,oneof=activate deactivate"`
}

// BulkUserActionRequest represents a bulk action on users
type BulkUserActionRequest struct {
	UserIDs []uuid.UUID `json:"user_ids" validate:"required,min=1"`
	Action  string      `json:"action" validate:"required,oneof=activate deactivate"`
}
