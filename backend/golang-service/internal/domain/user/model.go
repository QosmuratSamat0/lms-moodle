// internal/domain/user/model.go
package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         string    `json:"role" db:"role"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type UserProfile struct {
	User    User
	Student *StudentProfile
	Teacher *TeacherProfile
	Manager *ManagerProfile
}

type StudentProfile struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	FirstName *string   `json:"first_name" db:"first_name"`
	LastName  *string   `json:"last_name" db:"last_name"`
	GroupName *string   `json:"group_name" db:"group_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type TeacherProfile struct {
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	FirstName  *string   `json:"first_name" db:"first_name"`
	LastName   *string   `json:"last_name" db:"last_name"`
	Department *string   `json:"department" db:"department"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type ManagerProfile struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	FirstName *string   `json:"first_name" db:"first_name"`
	LastName  *string   `json:"last_name" db:"last_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
