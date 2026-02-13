package grade

import "time"

type Grade struct {
	ID           string    `db:"id" json:"id"`
	SubmissionID string    `db:"submission_id" json:"submission_id"`
	Score        float64   `db:"score" json:"score"`
	Feedback     string    `db:"feedback" json:"feedback"`
	GradedBy     string    `db:"graded_by_teacher_id" json:"graded_by"`
	GradedAt     time.Time `db:"graded_at" json:"graded_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateGradeInput struct {
	SubmissionID string
	Score        float64
	Feedback     string
	GradedBy     string
}

type Repository interface {
	Create(grade *Grade) error
	GetByID(id string) (*Grade, error)
	GetBySubmission(submissionID string) (*Grade, error)
	Update(grade *Grade) error
	Delete(id string) error
}
