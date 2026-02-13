package group

import "time"

type Group struct {
	ID          string    `json:"id"`
	CourseID    string    `json:"course_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	MaxStudents int       `json:"max_students"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GroupMember struct {
	ID               string    `json:"id"`
	GroupID          string    `json:"group_id"`
	StudentID        string    `json:"student_id"`
	JoinedAt         time.Time `json:"joined_at"`
	StudentFirstName string    `json:"student_first_name,omitempty"`
	StudentLastName  string    `json:"student_last_name,omitempty"`
	StudentEmail     string    `json:"student_email,omitempty"`
}

type CreateGroupInput struct {
	CourseID    string `json:"course_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	MaxStudents int    `json:"max_students"`
}

type UpdateGroupInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	MaxStudents *int    `json:"max_students"`
}

type AddMemberInput struct {
	GroupID   string `json:"group_id" binding:"required"`
	StudentID string `json:"student_id" binding:"required"`
}
