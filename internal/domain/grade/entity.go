package grade

import "time"

type Grade struct {
	ID           string    `db:"id"`
	SubmissionID string    `db:"submission_id"`
	Score        int       `db:"score"`
	Feedback     string    `db:"feedback"`
	GradedBy     string    `db:"graded_by"`
	GradedAt     time.Time `db:"graded_at"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type CreateGradeInput struct {
	SubmissionID string
	Score        int
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
