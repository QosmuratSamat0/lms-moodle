package quiz

import (
	"time"
)

type QuestionType string

const (
	MultipleChoice QuestionType = "multiple_choice"
	TrueFalse      QuestionType = "true_false"
	ShortAnswer    QuestionType = "short_answer"
)

type Quiz struct {
	ID          string     `json:"id"`
	CourseID    string     `json:"course_id"`
	CreatedBy   string     `json:"created_by"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	TimeLimit   int        `json:"time_limit_minutes,omitempty"` // in minutes, 0 = no limit
	MaxAttempts int        `json:"max_attempts"`                 // 0 = unlimited
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Published   bool       `json:"published"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Question struct {
	ID            string       `json:"id"`
	QuizID        string       `json:"quiz_id"`
	Type          QuestionType `json:"type"`
	Text          string       `json:"text"`
	Options       []string     `json:"options,omitempty"` // Options for multiple choice
	CorrectAnswer string       `json:"correct_answer"`
	Points        int          `json:"points"`
	OrderIndex    int          `json:"order_index"`
}

type QuizAttempt struct {
	ID          string     `json:"id"`
	QuizID      string     `json:"quiz_id"`
	StudentID   string     `json:"student_id"`
	StartedAt   time.Time  `json:"started_at"`
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
	Score       *float64   `json:"score,omitempty"`
	MaxScore    int        `json:"max_score"`
	Percentage  *float64   `json:"percentage,omitempty"`
}

type QuizAnswer struct {
	ID         string    `json:"id"`
	AttemptID  string    `json:"attempt_id"`
	QuestionID string    `json:"question_id"`
	Answer     string    `json:"answer"`
	IsCorrect  bool      `json:"is_correct"`
	Points     int       `json:"points_earned"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateQuizInput struct {
	CourseID    string          `json:"course_id" binding:"required"`
	Title       string          `json:"title" binding:"required"`
	Description string          `json:"description"`
	TimeLimit   int             `json:"time_limit_minutes"`
	MaxAttempts int             `json:"max_attempts"`
	StartDate   *time.Time      `json:"start_date"`
	EndDate     *time.Time      `json:"end_date"`
	Questions   []QuestionInput `json:"questions" binding:"required,min=1"`
}

type QuestionInput struct {
	Type          QuestionType `json:"type" binding:"required"`
	Text          string       `json:"text" binding:"required"`
	Options       []string     `json:"options"`
	CorrectAnswer string       `json:"correct_answer" binding:"required"`
	Points        int          `json:"points" binding:"required,min=1"`
}

type SubmitQuizInput struct {
	Answers []AnswerInput `json:"answers" binding:"required"`
}

type AnswerInput struct {
	QuestionID string `json:"question_id" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
}

type QuizResult struct {
	AttemptID   string         `json:"attempt_id"`
	QuizID      string         `json:"quiz_id"`
	QuizTitle   string         `json:"quiz_title"`
	StudentID   string         `json:"student_id"`
	Score       float64        `json:"score"`
	MaxScore    int            `json:"max_score"`
	Percentage  float64        `json:"percentage"`
	Passed      bool           `json:"passed"`
	SubmittedAt time.Time      `json:"submitted_at"`
	Answers     []AnswerResult `json:"answers,omitempty"`
}

type AnswerResult struct {
	QuestionID    string `json:"question_id"`
	QuestionText  string `json:"question_text"`
	YourAnswer    string `json:"your_answer"`
	CorrectAnswer string `json:"correct_answer"`
	IsCorrect     bool   `json:"is_correct"`
	PointsEarned  int    `json:"points_earned"`
	MaxPoints     int    `json:"max_points"`
}
