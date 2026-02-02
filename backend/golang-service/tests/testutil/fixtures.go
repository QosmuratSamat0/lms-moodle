package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Fixtures provides test data seeding utilities.
type Fixtures struct {
	pool *pgxpool.Pool
}

// NewFixtures creates a new Fixtures helper.
func NewFixtures(pool *pgxpool.Pool) *Fixtures {
	return &Fixtures{pool: pool}
}

// CreateUser creates a user with the given parameters.
func (f *Fixtures) CreateUser(ctx context.Context, t *testing.T, id uuid.UUID, email, passwordHash, role string) {
	t.Helper()

	query := `
		INSERT INTO users (id, email, password_hash, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, true, now(), now())
		ON CONFLICT (id) DO NOTHING`

	_, err := f.pool.Exec(ctx, query, id, email, passwordHash, role)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
}

// CreateTeacher creates a teacher profile for an existing user.
func (f *Fixtures) CreateTeacher(ctx context.Context, t *testing.T, userID uuid.UUID, firstName, lastName, department string) {
	t.Helper()

	query := `
		INSERT INTO teachers (user_id, first_name, last_name, department, created_at)
		VALUES ($1, $2, $3, $4, now())
		ON CONFLICT (user_id) DO NOTHING`

	_, err := f.pool.Exec(ctx, query, userID, firstName, lastName, department)
	if err != nil {
		t.Fatalf("failed to create test teacher: %v", err)
	}
}

// CreateStudent creates a student profile for an existing user.
func (f *Fixtures) CreateStudent(ctx context.Context, t *testing.T, userID uuid.UUID, firstName, lastName, groupName string) {
	t.Helper()

	query := `
		INSERT INTO students (user_id, first_name, last_name, group_name, created_at)
		VALUES ($1, $2, $3, $4, now())
		ON CONFLICT (user_id) DO NOTHING`

	_, err := f.pool.Exec(ctx, query, userID, firstName, lastName, groupName)
	if err != nil {
		t.Fatalf("failed to create test student: %v", err)
	}
}

// CreateCourse creates a course with the given parameters and returns its ID.
func (f *Fixtures) CreateCourse(ctx context.Context, t *testing.T, id uuid.UUID, title string, description *string, ownerTeacherID *uuid.UUID, isActive bool) {
	t.Helper()

	query := `
		INSERT INTO courses (id, title, description, owner_teacher_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, now(), now())
		ON CONFLICT (id) DO NOTHING`

	_, err := f.pool.Exec(ctx, query, id, title, description, ownerTeacherID, isActive)
	if err != nil {
		t.Fatalf("failed to create test course: %v", err)
	}
}

// CreateEnrollment creates an enrollment.
func (f *Fixtures) CreateEnrollment(ctx context.Context, t *testing.T, id, courseID, studentID uuid.UUID, status string) {
	t.Helper()

	query := `
		INSERT INTO enrollments (id, course_id, student_id, status, enrolled_at)
		VALUES ($1, $2, $3, $4, now())
		ON CONFLICT (id) DO NOTHING`

	_, err := f.pool.Exec(ctx, query, id, courseID, studentID, status)
	if err != nil {
		t.Fatalf("failed to create test enrollment: %v", err)
	}
}

// CreateAssignment creates an assignment.
func (f *Fixtures) CreateAssignment(ctx context.Context, t *testing.T, id, courseID uuid.UUID, title string, description *string, dueAt *time.Time, maxPoints int, createdByTeacherID *uuid.UUID) {
	t.Helper()

	query := `
		INSERT INTO assignments (id, course_id, title, description, due_at, max_points, created_by_teacher_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, now())
		ON CONFLICT (id) DO NOTHING`

	_, err := f.pool.Exec(ctx, query, id, courseID, title, description, dueAt, maxPoints, createdByTeacherID)
	if err != nil {
		t.Fatalf("failed to create test assignment: %v", err)
	}
}

// CreateSubmission creates a submission.
func (f *Fixtures) CreateSubmission(ctx context.Context, t *testing.T, id, assignmentID, studentID uuid.UUID, contentText *string, status string) {
	t.Helper()

	query := `
		INSERT INTO submissions (id, assignment_id, student_id, content_text, status, submitted_at)
		VALUES ($1, $2, $3, $4, $5, now())
		ON CONFLICT (id) DO NOTHING`

	_, err := f.pool.Exec(ctx, query, id, assignmentID, studentID, contentText, status)
	if err != nil {
		t.Fatalf("failed to create test submission: %v", err)
	}
}

// CreateGrade creates a grade for a submission.
func (f *Fixtures) CreateGrade(ctx context.Context, t *testing.T, id, submissionID uuid.UUID, gradedByTeacherID *uuid.UUID, score int, feedback *string) {
	t.Helper()

	query := `
		INSERT INTO grades (id, submission_id, graded_by_teacher_id, score, feedback, graded_at)
		VALUES ($1, $2, $3, $4, $5, now())
		ON CONFLICT (id) DO NOTHING`

	_, err := f.pool.Exec(ctx, query, id, submissionID, gradedByTeacherID, score, feedback)
	if err != nil {
		t.Fatalf("failed to create test grade: %v", err)
	}
}

// SeedBasicData seeds a minimal set of test data.
func (f *Fixtures) SeedBasicData(ctx context.Context, t *testing.T) {
	t.Helper()

	// Create admin
	f.CreateUser(ctx, t, TestAdminID, "admin@test.com", "$2a$10$testhashadmin", "admin")

	// Create teacher with profile
	f.CreateUser(ctx, t, TestTeacherID, "teacher@test.com", "$2a$10$testhashteacher", "teacher")
	f.CreateTeacher(ctx, t, TestTeacherID, "Test", "Teacher", "Computer Science")

	// Create second teacher
	f.CreateUser(ctx, t, TestTeacherID2, "teacher2@test.com", "$2a$10$testhashteacher2", "teacher")
	f.CreateTeacher(ctx, t, TestTeacherID2, "Another", "Teacher", "Mathematics")

	// Create student with profile
	f.CreateUser(ctx, t, TestStudentID, "student@test.com", "$2a$10$testhashstudent", "student")
	f.CreateStudent(ctx, t, TestStudentID, "Test", "Student", "CS-101")

	// Create second student
	f.CreateUser(ctx, t, TestStudentID2, "student2@test.com", "$2a$10$testhashstudent2", "student")
	f.CreateStudent(ctx, t, TestStudentID2, "Another", "Student", "CS-102")

	// Create courses
	desc := "Test course description"
	f.CreateCourse(ctx, t, TestCourseID1, "Test Course 1", &desc, &TestTeacherID, true)
	f.CreateCourse(ctx, t, TestCourseID2, "Test Course 2", nil, &TestTeacherID, true)

	// Create enrollments
	f.CreateEnrollment(ctx, t, TestEnrollmentID1, TestCourseID1, TestStudentID, "active")

	// Create assignment
	dueDate := time.Now().Add(7 * 24 * time.Hour)
	f.CreateAssignment(ctx, t, TestAssignmentID1, TestCourseID1, "Test Assignment", &desc, &dueDate, 100, &TestTeacherID)
}

// CreateChatRoom creates a chat room with the given parameters.
func (f *Fixtures) CreateChatRoom(ctx context.Context, t *testing.T, id, createdBy uuid.UUID, name *string, roomType string) {
	t.Helper()

	query := `
		INSERT INTO chat_rooms (id, name, type, created_by, created_at)
		VALUES ($1, $2, $3, $4, now())
		ON CONFLICT (id) DO NOTHING`

	_, err := f.pool.Exec(ctx, query, id, name, roomType, createdBy)
	if err != nil {
		t.Fatalf("failed to create test chat room: %v", err)
	}
}

// CreateChatRoomMember adds a member to a chat room.
func (f *Fixtures) CreateChatRoomMember(ctx context.Context, t *testing.T, roomID, userID uuid.UUID, role string) {
	t.Helper()

	query := `
		INSERT INTO chat_room_members (room_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, now())
		ON CONFLICT (room_id, user_id) DO UPDATE SET left_at = NULL, role = $3`

	_, err := f.pool.Exec(ctx, query, roomID, userID, role)
	if err != nil {
		t.Fatalf("failed to create test chat room member: %v", err)
	}
}

// CreateChatMessage creates a chat message.
func (f *Fixtures) CreateChatMessage(ctx context.Context, t *testing.T, id, roomID, senderID uuid.UUID, content string) {
	t.Helper()

	query := `
		INSERT INTO chat_messages (id, room_id, sender_id, content, created_at)
		VALUES ($1, $2, $3, $4, now())
		ON CONFLICT (id) DO NOTHING`

	_, err := f.pool.Exec(ctx, query, id, roomID, senderID, content)
	if err != nil {
		t.Fatalf("failed to create test chat message: %v", err)
	}
}
