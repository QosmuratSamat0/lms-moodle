package errors

import "errors"

// Common domain errors
var (
	// User errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrForbidden          = errors.New("access forbidden")

	// Course errors
	ErrCourseNotFound      = errors.New("course not found")
	ErrCourseAlreadyExists = errors.New("course with this code already exists")
	ErrNotCourseOwner      = errors.New("you are not the owner of this course")
	ErrCourseInactive      = errors.New("course is not active")

	// Enrollment errors
	ErrAlreadyEnrolled    = errors.New("already enrolled in this course")
	ErrNotEnrolled        = errors.New("not enrolled in this course")
	ErrEnrollmentNotFound = errors.New("enrollment not found")

	// Student errors
	ErrStudentNotFound = errors.New("student not found")
	ErrAlreadyExists   = errors.New("record already exists")

	// Teacher errors
	ErrTeacherNotFound = errors.New("teacher not found")

	// Admin errors
	ErrAdminNotFound = errors.New("admin not found")

	// Manager errors
	ErrManagerNotFound = errors.New("manager not found")

	// Category Manager errors
	ErrCategoryManagerNotFound = errors.New("category manager not found")

	// Course Category errors
	ErrCategoryNotFound      = errors.New("course category not found")
	ErrCategoryAlreadyExists = errors.New("category with this name already exists")

	// Assignment errors
	ErrAssignmentNotFound   = errors.New("assignment not found")
	ErrAssignmentPastDue    = errors.New("assignment due date has passed")
	ErrAssignmentNotStarted = errors.New("assignment is not yet available")

	// Submission errors
	ErrSubmissionNotFound     = errors.New("submission not found")
	ErrAlreadySubmitted       = errors.New("assignment already submitted")
	ErrSubmissionPastDeadline = errors.New("submission deadline has passed")

	// Grade errors
	ErrGradeNotFound = errors.New("grade not found")
	ErrAlreadyGraded = errors.New("submission already graded")
	ErrInvalidScore  = errors.New("score cannot exceed maximum points")

	// Group errors
	ErrGroupNotFound  = errors.New("group not found")
	ErrGroupFull      = errors.New("group has reached maximum capacity")
	ErrAlreadyInGroup = errors.New("student is already in this group")
	ErrNotInGroup     = errors.New("student is not in this group")

	// Quiz errors
	ErrQuizNotFound            = errors.New("quiz not found")
	ErrQuizNotPublished        = errors.New("quiz is not published")
	ErrQuizNotStarted          = errors.New("quiz has not started yet")
	ErrQuizEnded               = errors.New("quiz has ended")
	ErrMaxAttemptsReached      = errors.New("maximum quiz attempts reached")
	ErrAttemptNotFound         = errors.New("quiz attempt not found")
	ErrAttemptAlreadySubmitted = errors.New("quiz attempt already submitted")
	ErrQuestionNotFound        = errors.New("question not found")
	ErrInvalidQuestionType     = errors.New("invalid question type")

	// Announcement errors
	ErrAnnouncementNotFound = errors.New("announcement not found")

	// Appeal errors
	ErrAppealNotFound        = errors.New("grade appeal not found")
	ErrAppealAlreadyExists   = errors.New("pending appeal already exists for this grade")
	ErrAppealAlreadyResolved = errors.New("appeal has already been resolved")
	ErrCannotAppealOwnGrade  = errors.New("teachers cannot appeal their own grades")

	// Attendance errors
	ErrAttendanceNotFound  = errors.New("attendance record not found")
	ErrSessionNotFound     = errors.New("attendance session not found")
	ErrDuplicateAttendance = errors.New("attendance already marked for this session")

	// Notification errors
	ErrNotificationNotFound = errors.New("notification not found")

	// Chat errors
	ErrChatNotFound  = errors.New("chat not found")
	ErrNotChatMember = errors.New("you are not a member of this chat")

	// Upload errors
	ErrUploadFailed    = errors.New("file upload failed")
	ErrInvalidFileType = errors.New("invalid file type")
	ErrFileTooLarge    = errors.New("file size exceeds limit")

	// Validation errors
	ErrInvalidInput     = errors.New("invalid input data")
	ErrMissingRequired  = errors.New("missing required field")
	ErrInvalidID        = errors.New("invalid ID format")
	ErrInvalidDateRange = errors.New("invalid date range")
)

// IsNotFoundError checks if the error is a "not found" type error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrUserNotFound) ||
		errors.Is(err, ErrCourseNotFound) ||
		errors.Is(err, ErrAssignmentNotFound) ||
		errors.Is(err, ErrSubmissionNotFound) ||
		errors.Is(err, ErrGradeNotFound) ||
		errors.Is(err, ErrGroupNotFound) ||
		errors.Is(err, ErrQuizNotFound) ||
		errors.Is(err, ErrAttemptNotFound) ||
		errors.Is(err, ErrQuestionNotFound) ||
		errors.Is(err, ErrAnnouncementNotFound) ||
		errors.Is(err, ErrAppealNotFound) ||
		errors.Is(err, ErrAttendanceNotFound) ||
		errors.Is(err, ErrSessionNotFound) ||
		errors.Is(err, ErrNotificationNotFound) ||
		errors.Is(err, ErrChatNotFound) ||
		errors.Is(err, ErrEnrollmentNotFound) ||
		errors.Is(err, ErrStudentNotFound) ||
		errors.Is(err, ErrTeacherNotFound) ||
		errors.Is(err, ErrAdminNotFound) ||
		errors.Is(err, ErrManagerNotFound) ||
		errors.Is(err, ErrCategoryManagerNotFound) ||
		errors.Is(err, ErrCategoryNotFound)
}

// IsConflictError checks if the error is a conflict/duplicate type error
func IsConflictError(err error) bool {
	return errors.Is(err, ErrUserAlreadyExists) ||
		errors.Is(err, ErrCourseAlreadyExists) ||
		errors.Is(err, ErrAlreadyEnrolled) ||
		errors.Is(err, ErrAlreadySubmitted) ||
		errors.Is(err, ErrAlreadyGraded) ||
		errors.Is(err, ErrAlreadyInGroup) ||
		errors.Is(err, ErrAppealAlreadyExists) ||
		errors.Is(err, ErrDuplicateAttendance) ||
		errors.Is(err, ErrAlreadyExists) ||
		errors.Is(err, ErrCategoryAlreadyExists)
}

// IsForbiddenError checks if the error is an authorization/permission error
func IsForbiddenError(err error) bool {
	return errors.Is(err, ErrUnauthorized) ||
		errors.Is(err, ErrForbidden) ||
		errors.Is(err, ErrNotCourseOwner) ||
		errors.Is(err, ErrNotEnrolled) ||
		errors.Is(err, ErrNotInGroup) ||
		errors.Is(err, ErrNotChatMember) ||
		errors.Is(err, ErrCannotAppealOwnGrade)
}

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	return errors.Is(err, ErrInvalidInput) ||
		errors.Is(err, ErrMissingRequired) ||
		errors.Is(err, ErrInvalidID) ||
		errors.Is(err, ErrInvalidDateRange) ||
		errors.Is(err, ErrInvalidScore) ||
		errors.Is(err, ErrInvalidQuestionType) ||
		errors.Is(err, ErrInvalidFileType) ||
		errors.Is(err, ErrFileTooLarge)
}
