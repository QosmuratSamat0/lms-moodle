package student

import "context"

type Repository interface {
	Create(ctx context.Context, student *Student) error
	GetByID(ctx context.Context, id string) (*Student, error)
	GetByUserID(ctx context.Context, userID string) (*Student, error)
	GetByStudentCode(ctx context.Context, code string) (*Student, error)
	GetWithDetails(ctx context.Context, id string) (*StudentWithDetails, error)
	List(ctx context.Context, filter *StudentFilter) ([]*Student, int64, error)
	ListByMajor(ctx context.Context, major string, limit, offset int) ([]*Student, int64, error)
	Update(ctx context.Context, student *Student) error
	Delete(ctx context.Context, id string) error

	// Enrollment related
	GetEnrollments(ctx context.Context, studentID string) ([]*StudentEnrollment, error)

	// Group related
	GetGroups(ctx context.Context, studentID string) ([]*StudentGroup, error)
}

// StudentEnrollment combines enrollment with course info
type StudentEnrollment struct {
	EnrollmentID string `json:"enrollment_id"`
	CourseID     string `json:"course_id"`
	CourseName   string `json:"course_name"`
	CourseCode   string `json:"course_code"`
	TeacherName  string `json:"teacher_name"`
	Status       string `json:"status"`
	EnrolledAt   string `json:"enrolled_at"`
}

// StudentGroup combines group membership with group info
type StudentGroup struct {
	MembershipID string `json:"membership_id"`
	GroupID      string `json:"group_id"`
	GroupName    string `json:"group_name"`
	CourseID     string `json:"course_id"`
	CourseName   string `json:"course_name"`
	JoinedAt     string `json:"joined_at"`
}
