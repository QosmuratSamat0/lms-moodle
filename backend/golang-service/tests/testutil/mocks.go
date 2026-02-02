package testutil

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ============================================================================
// Mock DBTX Interface (matches database.DBTX)
// ============================================================================

// MockDBTX is a mock implementation of the DBTX interface.
type MockDBTX struct {
	ExecFn     func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryFn    func(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRowFn func(ctx context.Context, sql string, args ...any) pgx.Row
}

func (m *MockDBTX) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	if m.ExecFn != nil {
		return m.ExecFn(ctx, sql, arguments...)
	}
	return pgconn.NewCommandTag(""), nil
}

func (m *MockDBTX) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if m.QueryFn != nil {
		return m.QueryFn(ctx, sql, args...)
	}
	return nil, nil
}

func (m *MockDBTX) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if m.QueryRowFn != nil {
		return m.QueryRowFn(ctx, sql, args...)
	}
	return &MockRow{}
}

// MockRow is a mock implementation of pgx.Row.
type MockRow struct {
	ScanFn func(dest ...any) error
}

func (m *MockRow) Scan(dest ...any) error {
	if m.ScanFn != nil {
		return m.ScanFn(dest...)
	}
	return nil
}

// ============================================================================
// Fixed UUIDs for deterministic tests
// ============================================================================

var (
	// User IDs
	TestUserID1    = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	TestUserID2    = uuid.MustParse("00000000-0000-0000-0000-000000000002")
	TestUserID3    = uuid.MustParse("00000000-0000-0000-0000-000000000003")
	TestAdminID    = uuid.MustParse("00000000-0000-0000-0000-000000000010")
	TestTeacherID  = uuid.MustParse("00000000-0000-0000-0000-000000000011")
	TestTeacherID2 = uuid.MustParse("00000000-0000-0000-0000-000000000012")
	TestStudentID  = uuid.MustParse("00000000-0000-0000-0000-000000000021")
	TestStudentID2 = uuid.MustParse("00000000-0000-0000-0000-000000000022")
	TestManagerID  = uuid.MustParse("00000000-0000-0000-0000-000000000031")

	// Course IDs
	TestCourseID1 = uuid.MustParse("00000000-0000-0000-0001-000000000001")
	TestCourseID2 = uuid.MustParse("00000000-0000-0000-0001-000000000002")
	TestCourseID3 = uuid.MustParse("00000000-0000-0000-0001-000000000003")

	// Assignment IDs
	TestAssignmentID1 = uuid.MustParse("00000000-0000-0000-0002-000000000001")
	TestAssignmentID2 = uuid.MustParse("00000000-0000-0000-0002-000000000002")

	// Submission IDs
	TestSubmissionID1 = uuid.MustParse("00000000-0000-0000-0003-000000000001")
	TestSubmissionID2 = uuid.MustParse("00000000-0000-0000-0003-000000000002")

	// Enrollment IDs
	TestEnrollmentID1 = uuid.MustParse("00000000-0000-0000-0004-000000000001")
	TestEnrollmentID2 = uuid.MustParse("00000000-0000-0000-0004-000000000002")

	// Grade IDs
	TestGradeID1 = uuid.MustParse("00000000-0000-0000-0005-000000000001")
	TestGradeID2 = uuid.MustParse("00000000-0000-0000-0005-000000000002")

	// Session IDs
	TestSessionID1 = uuid.MustParse("00000000-0000-0000-0006-000000000001")
	TestSessionID2 = uuid.MustParse("00000000-0000-0000-0006-000000000002")

	// Chat Room IDs
	TestChatRoomID1 = uuid.MustParse("00000000-0000-0000-0007-000000000001")
	TestChatRoomID2 = uuid.MustParse("00000000-0000-0000-0007-000000000002")
)

// NewTestUUID generates a deterministic UUID based on an index.
// Useful for generating multiple UUIDs in table-driven tests.
func NewTestUUID(index int) uuid.UUID {
	return uuid.MustParse(uuidFromIndex(index))
}

func uuidFromIndex(index int) string {
	return uuid.MustParse("00000000-0000-0000-9999-" + padLeft(index, 12)).String()
}

func padLeft(n, width int) string {
	s := ""
	for i := 0; i < width; i++ {
		s += "0"
	}
	ns := s + string(rune('0'+n%10))
	if n >= 10 {
		ns = s[:len(s)-2] + string(rune('0'+n/10)) + string(rune('0'+n%10))
	}
	return ns[len(ns)-width:]
}
