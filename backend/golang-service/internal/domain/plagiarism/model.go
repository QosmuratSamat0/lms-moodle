package plagiarism

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Report represents a plagiarism check report
type Report struct {
	ID           uuid.UUID       `json:"id" db:"id"`
	SubmissionID uuid.UUID       `json:"submission_id" db:"submission_id"`
	Score        float64         `json:"score" db:"score"`
	Details      json.RawMessage `json:"details" db:"details"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
}

// ReportDetails contains detailed plagiarism analysis
type ReportDetails struct {
	Matches      []Match  `json:"matches"`
	Sources      []Source `json:"sources"`
	AnalysisTime int64    `json:"analysis_time_ms"`
}

// Match represents a plagiarism match
type Match struct {
	SourceID    string  `json:"source_id"`
	StartIndex  int     `json:"start_index"`
	EndIndex    int     `json:"end_index"`
	MatchedText string  `json:"matched_text"`
	Similarity  float64 `json:"similarity"`
}

// Source represents a matched source document
type Source struct {
	ID           string  `json:"id"`
	Type         string  `json:"type"` // "submission", "external"
	SubmissionID *string `json:"submission_id,omitempty"`
	URL          *string `json:"url,omitempty"`
	Similarity   float64 `json:"similarity"`
}

// ReportWithSubmission includes submission details
type ReportWithSubmission struct {
	Report
	StudentID    uuid.UUID `json:"student_id" db:"student_id"`
	AssignmentID uuid.UUID `json:"assignment_id" db:"assignment_id"`
	ContentText  *string   `json:"content_text" db:"content_text"`
	StudentEmail string    `json:"student_email" db:"student_email"`
}

// CheckStatus represents the status of a plagiarism check
type CheckStatus string

const (
	CheckStatusPending    CheckStatus = "pending"
	CheckStatusProcessing CheckStatus = "processing"
	CheckStatusCompleted  CheckStatus = "completed"
	CheckStatusFailed     CheckStatus = "failed"
)

// IsSuspicious returns true if score exceeds threshold
func (r *Report) IsSuspicious(threshold float64) bool {
	return r.Score >= threshold
}
