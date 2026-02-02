package repository

import "ap1-final-mini-moodle/internal/entity"

type CourseRepository interface {
	CreateCourse(course entity.Course) (entity.Course, error)
	GetCourse(id int64) (entity.Course, bool)
	ListCourses() []entity.Course
}

type AssignmentRepository interface {
	CreateAssignment(assignment entity.Assignment) (entity.Assignment, error)
	GetAssignment(id int64) (entity.Assignment, bool)
}

type SubmissionRepository interface {
	CreateSubmission(submission entity.Submission) (entity.Submission, error)
	GetSubmission(id int64) (entity.Submission, bool)
	ListSubmissionsByAssignment(assignmentID int64) []entity.Submission
}

type GradeRepository interface {
	UpsertGrade(grade entity.Grade) (entity.Grade, error)
	GetGrade(submissionID int64) (entity.Grade, bool)
}

type AnalyticsRepository interface {
	CountStudentsInCourse(courseID int64) int
}

type Store interface {
	CourseRepository
	AssignmentRepository
	SubmissionRepository
	GradeRepository
	AnalyticsRepository
}
