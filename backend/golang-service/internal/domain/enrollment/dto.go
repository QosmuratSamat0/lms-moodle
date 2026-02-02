package enrollment

import "github.com/google/uuid"

// EnrollRequest represents a request to enroll in a course
type EnrollRequest struct {
	CourseID uuid.UUID `json:"course_id" validate:"required"`
}

// UpdateEnrollmentStatusRequest represents a request to update enrollment status
type UpdateEnrollmentStatusRequest struct {
	Status EnrollmentStatus `json:"status" validate:"required,oneof=active pending rejected dropped"`
}

// EnrollmentResponse represents an enrollment in API responses
type EnrollmentResponse struct {
	ID               string  `json:"id"`
	CourseID         string  `json:"course_id"`
	StudentID        string  `json:"student_id"`
	Status           string  `json:"status"`
	CourseTitle      string  `json:"course_title,omitempty"`
	StudentFirstName *string `json:"student_first_name,omitempty"`
	StudentLastName  *string `json:"student_last_name,omitempty"`
	StudentEmail     string  `json:"student_email,omitempty"`
	EnrolledAt       string  `json:"enrolled_at"`
}

// EnrollmentListResponse represents a paginated list of enrollments
type EnrollmentListResponse struct {
	Enrollments []EnrollmentResponse `json:"enrollments"`
	Total       int64                `json:"total"`
	Page        int                  `json:"page"`
	Limit       int                  `json:"limit"`
	TotalPages  int                  `json:"total_pages"`
}
