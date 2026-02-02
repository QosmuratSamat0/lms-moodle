package usecase

import (
	"strings"
	"time"

	"ap1-final-mini-moodle/internal/entity"
	"ap1-final-mini-moodle/internal/repository"
)

type SubmissionWithGrade struct {
	Submission entity.Submission
	Grade      *entity.Grade
}

type TeacherUsecase struct {
	store repository.Store
}

func NewTeacherUsecase(store repository.Store) *TeacherUsecase {
	return &TeacherUsecase{store: store}
}

func (u *TeacherUsecase) CreateCourse(teacherID int64, teacherName, title, description string) (entity.Course, error) {
	if teacherID <= 0 || strings.TrimSpace(title) == "" {
		return entity.Course{}, ErrInvalid
	}

	course := entity.Course{
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		TeacherID:   teacherID,
		TeacherName: strings.TrimSpace(teacherName),
		CreatedAt:   time.Now().UTC(),
	}

	return u.store.CreateCourse(course)
}

func (u *TeacherUsecase) CreateAssignment(teacherID, courseID int64, title, description string, dueAt time.Time, maxPoints int) (entity.Assignment, error) {
	if teacherID <= 0 || courseID <= 0 || strings.TrimSpace(title) == "" {
		return entity.Assignment{}, ErrInvalid
	}
	if maxPoints < 0 {
		return entity.Assignment{}, ErrInvalid
	}

	course, ok := u.store.GetCourse(courseID)
	if !ok {
		return entity.Assignment{}, ErrNotFound
	}
	if course.TeacherID != teacherID {
		return entity.Assignment{}, ErrForbidden
	}

	assignment := entity.Assignment{
		CourseID:    courseID,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		DueAt:       dueAt.UTC(),
		MaxPoints:   maxPoints,
		CreatedAt:   time.Now().UTC(),
	}

	return u.store.CreateAssignment(assignment)
}

func (u *TeacherUsecase) ListSubmissions(teacherID, assignmentID int64) ([]SubmissionWithGrade, error) {
	if teacherID <= 0 || assignmentID <= 0 {
		return nil, ErrInvalid
	}

	assignment, ok := u.store.GetAssignment(assignmentID)
	if !ok {
		return nil, ErrNotFound
	}
	course, ok := u.store.GetCourse(assignment.CourseID)
	if !ok {
		return nil, ErrNotFound
	}
	if course.TeacherID != teacherID {
		return nil, ErrForbidden
	}

	submissions := u.store.ListSubmissionsByAssignment(assignmentID)
	result := make([]SubmissionWithGrade, 0, len(submissions))
	for _, submission := range submissions {
		var gradePtr *entity.Grade
		if grade, ok := u.store.GetGrade(submission.ID); ok {
			gradeCopy := grade
			gradePtr = &gradeCopy
		}
		result = append(result, SubmissionWithGrade{
			Submission: submission,
			Grade:      gradePtr,
		})
	}

	return result, nil
}

func (u *TeacherUsecase) GradeSubmission(teacherID, submissionID int64, score int, feedback string) (entity.Grade, error) {
	if teacherID <= 0 || submissionID <= 0 {
		return entity.Grade{}, ErrInvalid
	}
	if score < 0 {
		return entity.Grade{}, ErrInvalid
	}

	submission, ok := u.store.GetSubmission(submissionID)
	if !ok {
		return entity.Grade{}, ErrNotFound
	}

	assignment, ok := u.store.GetAssignment(submission.AssignmentID)
	if !ok {
		return entity.Grade{}, ErrNotFound
	}

	course, ok := u.store.GetCourse(assignment.CourseID)
	if !ok {
		return entity.Grade{}, ErrNotFound
	}
	if course.TeacherID != teacherID {
		return entity.Grade{}, ErrForbidden
	}
	if assignment.MaxPoints > 0 && score > assignment.MaxPoints {
		return entity.Grade{}, ErrInvalid
	}

	grade := entity.Grade{
		SubmissionID: submissionID,
		GradedBy:     teacherID,
		Score:        score,
		Feedback:     strings.TrimSpace(feedback),
		GradedAt:     time.Now().UTC(),
	}

	return u.store.UpsertGrade(grade)
}
