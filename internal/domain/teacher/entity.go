package teacher

import "time"

type Teacher struct {
	ID             string    `json:"id" db:"id"`
	UserID         string    `json:"user_id" db:"user_id"`
	EmployeeID     string    `json:"employee_id" db:"employee_id"`
	FirstName      string    `json:"first_name" db:"first_name"`
	LastName       string    `json:"last_name" db:"last_name"`
	Department     string    `json:"department" db:"department"`
	Specialization string    `json:"specialization" db:"specialization"`
	Qualifications string    `json:"qualifications" db:"qualifications"`
	Bio            string    `json:"bio" db:"bio"`
	OfficeHours    string    `json:"office_hours" db:"office_hours"`
	Phone          string    `json:"phone" db:"phone"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`

	// Populated via joins
	Email string `json:"email,omitempty" db:"email"`
}

type CreateTeacherInput struct {
	UserID         string `json:"user_id" binding:"required"`
	EmployeeID     string `json:"employee_id" binding:"required"`
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"required"`
	Department     string `json:"department"`
	Specialization string `json:"specialization"`
	Qualifications string `json:"qualifications"`
	Bio            string `json:"bio"`
	OfficeHours    string `json:"office_hours"`
	Phone          string `json:"phone"`
}

type UpdateTeacherInput struct {
	Department     *string `json:"department"`
	Specialization *string `json:"specialization"`
	Qualifications *string `json:"qualifications"`
	Bio            *string `json:"bio"`
	OfficeHours    *string `json:"office_hours"`
	Phone          *string `json:"phone"`
	IsActive       *bool   `json:"is_active"`
}

type TeacherFilter struct {
	Department     string `form:"department"`
	Specialization string `form:"specialization"`
	IsActive       *bool  `form:"is_active"`
	Limit          int    `form:"limit,default=20"`
	Offset         int    `form:"offset,default=0"`
}

type TeacherGroup struct {
	ID          string    `json:"id"`
	CourseID    string    `json:"course_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	MaxStudents int       `json:"max_students"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TeacherCourse struct {
	ID          string    `db:"id" json:"id"`
	Code        string    `db:"code" json:"code"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	TeacherID   string    `db:"teacher_id" json:"owner_teacher_id"`
	MaxPoints   int       `db:"max_points" json:"max_points"`
	Active      bool      `db:"active" json:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
