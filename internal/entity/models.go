package entity

import "time"

type Course struct {
	ID          int64
	Title       string
	Description string
	TeacherID   int64
	TeacherName string
	CreatedAt   time.Time
}

type Assignment struct {
	ID          int64
	CourseID    int64
	Title       string
	Description string
	DueAt       time.Time
	MaxPoints   int
	CreatedAt   time.Time
}

type Submission struct {
	ID           int64
	AssignmentID int64
	StudentID    int64
	SubmittedAt  time.Time
	ContentRef   string
}

type Grade struct {
	SubmissionID int64
	GradedBy     int64
	Score        int
	Feedback     string
	GradedAt     time.Time
}
