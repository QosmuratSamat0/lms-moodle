package appeal

import "time"

type AppealStatus string

const (
	AppealPending  AppealStatus = "pending"
	AppealApproved AppealStatus = "approved"
	AppealRejected AppealStatus = "rejected"
)

type GradeAppeal struct {
	ID         string       `json:"id"`
	GradeID    string       `json:"grade_id"`
	StudentID  string       `json:"student_id"`
	Reason     string       `json:"reason"`
	Evidence   string       `json:"evidence,omitempty"`
	Status     AppealStatus `json:"status"`
	TeacherID  *string      `json:"teacher_id,omitempty"`
	Response   *string      `json:"response,omitempty"`
	NewScore   *float64     `json:"new_score,omitempty"`
	ResolvedAt *time.Time   `json:"resolved_at,omitempty"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

type CreateAppealInput struct {
	GradeID  string `json:"grade_id" binding:"required"`
	Reason   string `json:"reason" binding:"required,min=10"`
	Evidence string `json:"evidence"`
}

type ResolveAppealInput struct {
	Status   AppealStatus `json:"status" binding:"required,oneof=approved rejected"`
	Response string       `json:"response" binding:"required"`
	NewScore *float64     `json:"new_score"` // only for approved
}

type AppealWithDetails struct {
	GradeAppeal
	StudentName     string  `json:"student_name"`
	StudentEmail    string  `json:"student_email"`
	CourseName      string  `json:"course_name"`
	AssignmentTitle string  `json:"assignment_title"`
	OriginalScore   float64 `json:"original_score"`
	MaxPoints       float64 `json:"max_points"`
}
