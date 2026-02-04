package attendance

import (
	"context"
	"errors"
	"time"

	"github.com/ap1-final-mini-moodle/internal/shared/errorx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Service defines the attendance service interface
type Service interface {
	// Session operations
	CreateSession(ctx context.Context, teacherID uuid.UUID, req *CreateSessionRequest) (*SessionResponse, error)
	GetSession(ctx context.Context, id uuid.UUID) (*SessionAttendanceResponse, error)
	UpdateSession(ctx context.Context, id uuid.UUID, userID uuid.UUID, role string, req *UpdateSessionRequest) error
	DeleteSession(ctx context.Context, id uuid.UUID, userID uuid.UUID, role string) error
	ListSessionsByCourse(ctx context.Context, courseID uuid.UUID, page, limit int) (*SessionListResponse, error)

	// Mark operations
	MarkAttendance(ctx context.Context, sessionID uuid.UUID, markedBy uuid.UUID, req *MarkAttendanceRequest) error
	BulkMarkAttendance(ctx context.Context, sessionID uuid.UUID, markedBy uuid.UUID, req *BulkMarkAttendanceRequest) error

	// Summary
	GetStudentCourseSummary(ctx context.Context, studentID, courseID uuid.UUID) (*AttendanceSummaryResponse, error)
}

type service struct {
	repo Repository
}

// NewService creates a new attendance service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateSession(ctx context.Context, teacherID uuid.UUID, req *CreateSessionRequest) (*SessionResponse, error) {
	// Parse dates
	sessionDate, err := time.Parse("2006-01-02", req.SessionDate)
	if err != nil {
		return nil, errorx.NewValidationError("session_date: invalid date format, use YYYY-MM-DD")
	}

	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return nil, errorx.NewValidationError("start_time: invalid time format, use HH:MM")
	}

	var endTime *time.Time
	if req.EndTime != nil {
		et, err := time.Parse("15:04", *req.EndTime)
		if err != nil {
			return nil, errorx.NewValidationError("end_time: invalid time format, use HH:MM")
		}
		endTime = &et
	}

	session := &Session{
		CourseID:    req.CourseID,
		CreatedByID: teacherID,
		Title:       req.Title,
		SessionDate: sessionDate,
		StartTime:   startTime,
		EndTime:     endTime,
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, errorx.Wrap(err, "create session")
	}

	// Get with details
	detailed, err := s.repo.GetSessionByID(ctx, session.ID)
	if err != nil {
		return nil, errorx.Wrap(err, "get created session")
	}

	return s.toSessionResponse(detailed), nil
}

func (s *service) GetSession(ctx context.Context, id uuid.UUID) (*SessionAttendanceResponse, error) {
	session, err := s.repo.GetSessionByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("attendance session")
		}
		return nil, errorx.Wrap(err, "get session")
	}

	marks, err := s.repo.ListMarksBySession(ctx, id)
	if err != nil {
		return nil, errorx.Wrap(err, "list marks")
	}

	return &SessionAttendanceResponse{
		Session: *s.toSessionResponse(session),
		Marks:   s.toMarkResponses(marks),
	}, nil
}

func (s *service) UpdateSession(ctx context.Context, id uuid.UUID, userID uuid.UUID, role string, req *UpdateSessionRequest) error {
	session, err := s.repo.GetSessionByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("attendance session")
		}
		return errorx.Wrap(err, "get session")
	}

	// Check permission: admin or creator
	if role != "admin" && session.CreatedByID != userID {
		return errorx.NewForbiddenError("you can only update sessions you created")
	}

	// Update fields
	if req.Title != nil {
		session.Title = req.Title
	}
	if req.SessionDate != nil {
		sd, err := time.Parse("2006-01-02", *req.SessionDate)
		if err != nil {
			return errorx.NewValidationError("session_date: invalid date format")
		}
		session.SessionDate = sd
	}
	if req.StartTime != nil {
		st, err := time.Parse("15:04", *req.StartTime)
		if err != nil {
			return errorx.NewValidationError("start_time: invalid time format")
		}
		session.StartTime = st
	}
	if req.EndTime != nil {
		et, err := time.Parse("15:04", *req.EndTime)
		if err != nil {
			return errorx.NewValidationError("end_time: invalid time format")
		}
		session.EndTime = &et
	}

	if err := s.repo.UpdateSession(ctx, &session.Session); err != nil {
		return errorx.Wrap(err, "update session")
	}

	return nil
}

func (s *service) DeleteSession(ctx context.Context, id uuid.UUID, userID uuid.UUID, role string) error {
	session, err := s.repo.GetSessionByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("attendance session")
		}
		return errorx.Wrap(err, "get session")
	}

	if role != "admin" && session.CreatedByID != userID {
		return errorx.NewForbiddenError("you can only delete sessions you created")
	}

	if err := s.repo.DeleteSession(ctx, id); err != nil {
		return errorx.Wrap(err, "delete session")
	}

	return nil
}

func (s *service) ListSessionsByCourse(ctx context.Context, courseID uuid.UUID, page, limit int) (*SessionListResponse, error) {
	offset := (page - 1) * limit
	sessions, total, err := s.repo.ListSessionsByCourse(ctx, courseID, limit, offset)
	if err != nil {
		return nil, errorx.Wrap(err, "list sessions")
	}

	var responses []SessionResponse
	for _, sess := range sessions {
		responses = append(responses, *s.toSessionResponse(&sess))
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &SessionListResponse{
		Sessions:   responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *service) MarkAttendance(ctx context.Context, sessionID uuid.UUID, markedBy uuid.UUID, req *MarkAttendanceRequest) error {
	// Check if session exists
	_, err := s.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorx.NewNotFoundError("attendance session")
		}
		return errorx.Wrap(err, "get session")
	}

	// Check if already marked
	existing, err := s.repo.GetMarkBySessionAndStudent(ctx, sessionID, req.StudentID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return errorx.Wrap(err, "check existing mark")
	}

	if existing != nil {
		// Update existing mark
		existing.Status = AttendanceStatus(req.Status)
		existing.Notes = req.Notes
		existing.MarkedBy = markedBy

		if err := s.repo.UpdateMark(ctx, existing); err != nil {
			return errorx.Wrap(err, "update mark")
		}
		return nil
	}

	// Create new mark
	mark := &Mark{
		SessionID: sessionID,
		StudentID: req.StudentID,
		Status:    AttendanceStatus(req.Status),
		Notes:     req.Notes,
		MarkedBy:  markedBy,
	}

	if err := s.repo.CreateMark(ctx, mark); err != nil {
		return errorx.Wrap(err, "create mark")
	}

	return nil
}

func (s *service) BulkMarkAttendance(ctx context.Context, sessionID uuid.UUID, markedBy uuid.UUID, req *BulkMarkAttendanceRequest) error {
	for _, markReq := range req.Marks {
		if err := s.MarkAttendance(ctx, sessionID, markedBy, &markReq); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) GetStudentCourseSummary(ctx context.Context, studentID, courseID uuid.UUID) (*AttendanceSummaryResponse, error) {
	summary, err := s.repo.GetStudentCourseSummary(ctx, studentID, courseID)
	if err != nil {
		return nil, errorx.Wrap(err, "get attendance summary")
	}

	return &AttendanceSummaryResponse{
		StudentID:      summary.StudentID.String(),
		CourseID:       summary.CourseID.String(),
		TotalSessions:  summary.TotalSessions,
		PresentCount:   summary.PresentCount,
		AbsentCount:    summary.AbsentCount,
		LateCount:      summary.LateCount,
		ExcusedCount:   summary.ExcusedCount,
		AttendanceRate: summary.AttendanceRate,
	}, nil
}

// Helper methods

func (s *service) toSessionResponse(sess *SessionWithDetails) *SessionResponse {
	resp := &SessionResponse{
		ID:          sess.ID.String(),
		CourseID:    sess.CourseID.String(),
		CourseTitle: sess.CourseTitle,
		Title:       sess.Title,
		SessionDate: sess.SessionDate.Format("2006-01-02"),
		StartTime:   sess.StartTime.Format("15:04"),
		TeacherName: sess.TeacherName,
		TotalMarks:  sess.TotalMarks,
		CreatedAt:   sess.CreatedAt.Format(time.RFC3339),
	}
	if sess.EndTime != nil {
		et := sess.EndTime.Format("15:04")
		resp.EndTime = &et
	}
	return resp
}

func (s *service) toMarkResponses(marks []MarkWithDetails) []MarkResponse {
	var responses []MarkResponse
	for _, m := range marks {
		responses = append(responses, MarkResponse{
			ID:               m.ID.String(),
			SessionID:        m.SessionID.String(),
			StudentID:        m.StudentID.String(),
			StudentFirstName: m.StudentFirstName,
			StudentLastName:  m.StudentLastName,
			StudentEmail:     m.StudentEmail,
			Status:           string(m.Status),
			Notes:            m.Notes,
			MarkedAt:         m.MarkedAt.Format(time.RFC3339),
		})
	}
	return responses
}

var _ Service = (*service)(nil)
