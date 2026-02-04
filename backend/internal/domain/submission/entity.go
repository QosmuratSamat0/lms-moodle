package submission

import "time"

type Submission struct {
	ID           string     `db:"id"`
	AssignmentID string     `db:"assignment_id"`
	StudentID    string     `db:"student_id"`
	Content      string     `db:"content"`
	FileURL      *string    `db:"file_url"`
	SubmittedAt  time.Time  `db:"submitted_at"`
	GradedAt     *time.Time `db:"graded_at"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

type CreateSubmissionInput struct {
	AssignmentID string
	StudentID    string
	Content      string
	FileURL      *string
}

type Repository interface {
	Create(submission *Submission) error
	GetByID(id string) (*Submission, error)
	GetByAssignmentAndStudent(assignmentID, studentID string) (*Submission, error)
	ListByAssignment(assignmentID string, skip, take int) ([]*Submission, error)
	ListByStudent(studentID string, skip, take int) ([]*Submission, error)
	Update(submission *Submission) error
	Delete(id string) error
}
