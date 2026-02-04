package user

import "time"

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleTeacher  Role = "teacher"
	RoleStudent  Role = "student"
	RoleManager  Role = "manager"
)

type User struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Role      Role      `db:"role"`
	Active    bool      `db:"active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CreateUserInput struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Role      Role
}

type UpdateUserInput struct {
	FirstName *string
	LastName  *string
	Active    *bool
}

type Repository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	List(skip, take int) ([]*User, error)
	Delete(id string) error
}
