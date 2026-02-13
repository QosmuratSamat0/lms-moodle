package submission

import "time"

type Submission struct {
	ID               string    `json:"id" db:"id"`
	AssignmentID     string    `json:"assignment_id" db:"assignment_id"`
	StudentID        string    `json:"student_id" db:"student_id"`
	ContentText      string    `json:"content_text" db:"content_text"`
	FileURL          *string   `json:"file_url" db:"file_url"`
	SubmittedAt      time.Time `json:"submitted_at" db:"submitted_at"`
	Status           string    `json:"status" db:"status"`
	StudentFirstName string    `json:"student_first_name,omitempty" db:"student_first_name"`
	StudentLastName  string    `json:"student_last_name,omitempty" db:"student_last_name"`
	StudentEmail     string    `json:"student_email,omitempty" db:"student_email"`
	GradeScore       *float64  `json:"grade_score,omitempty"`
	GradeFeedback    *string   `json:"grade_feedback,omitempty"`
	GradedAt         *time.Time `json:"graded_at,omitempty"`
	GraderFirstName  *string   `json:"grader_first_name,omitempty"`
	GraderLastName   *string   `json:"grader_last_name,omitempty"`
}

type CreateSubmissionInput struct {
	AssignmentID string
	StudentID    string
	ContentText  string
	FileURL      *string
}

type UpdateSubmissionInput struct {
	ID          string
	ContentText string
	FileURL     *string
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
