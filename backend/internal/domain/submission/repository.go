package submission

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AssignmentInfo contains deadline info for submission validation
type AssignmentInfo struct {
	DueAt     *time.Time
	AllowLate bool
}

// Repository defines the submission repository interface
type Repository interface {
	Create(ctx context.Context, submission *Submission) error
	GetByID(ctx context.Context, id uuid.UUID) (*SubmissionWithDetails, error)
	GetByAssignmentAndStudent(ctx context.Context, assignmentID, studentID uuid.UUID) (*Submission, error)
	GetAssignmentInfo(ctx context.Context, assignmentID uuid.UUID) (*AssignmentInfo, error)
	Update(ctx context.Context, submission *Submission) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByAssignment(ctx context.Context, assignmentID uuid.UUID, groupID *uuid.UUID, limit, offset int) ([]SubmissionWithDetails, int64, error)
	ListByStudent(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]SubmissionWithDetails, int64, error)
}

type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new submission repository
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, submission *Submission) error {
	query := `
		INSERT INTO submissions (assignment_id, student_id, content_text, file_url, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, submitted_at`

	return r.db.QueryRow(ctx, query,
		submission.AssignmentID,
		submission.StudentID,
		submission.ContentText,
		submission.FileURL,
		submission.Status,
	).Scan(&submission.ID, &submission.SubmittedAt)
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*SubmissionWithDetails, error) {
	query := `
		SELECT
			s.id, s.assignment_id, s.student_id, s.content_text, s.file_url, s.submitted_at, s.status,
			a.title AS assignment_title, a.max_points, a.due_at, a.allow_late,
			c.title AS course_title,
			st.first_name AS student_first_name, st.last_name AS student_last_name,
			u.email AS student_email,
			g.score AS grade_score, g.feedback AS grade_feedback
		FROM submissions s
		JOIN assignments a ON s.assignment_id = a.id
		JOIN courses c ON a.course_id = c.id
		JOIN students st ON s.student_id = st.user_id
		JOIN users u ON st.user_id = u.id
		LEFT JOIN grades g ON s.id = g.submission_id
		WHERE s.id = $1`

	var sub SubmissionWithDetails
	err := r.db.QueryRow(ctx, query, id).Scan(
		&sub.ID, &sub.AssignmentID, &sub.StudentID, &sub.ContentText, &sub.FileURL, &sub.SubmittedAt, &sub.Status,
		&sub.AssignmentTitle, &sub.MaxPoints, &sub.DueAt, &sub.AllowLate,
		&sub.CourseTitle,
		&sub.StudentFirstName, &sub.StudentLastName, &sub.StudentEmail,
		&sub.GradeScore, &sub.GradeFeedback,
	)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *repository) GetAssignmentInfo(ctx context.Context, assignmentID uuid.UUID) (*AssignmentInfo, error) {
	query := `SELECT due_at, allow_late FROM assignments WHERE id = $1`
	var info AssignmentInfo
	err := r.db.QueryRow(ctx, query, assignmentID).Scan(&info.DueAt, &info.AllowLate)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (r *repository) GetByAssignmentAndStudent(ctx context.Context, assignmentID, studentID uuid.UUID) (*Submission, error) {
	query := `
		SELECT id, assignment_id, student_id, content_text, file_url, submitted_at, status
		FROM submissions
		WHERE assignment_id = $1 AND student_id = $2`

	var sub Submission
	err := r.db.QueryRow(ctx, query, assignmentID, studentID).Scan(
		&sub.ID, &sub.AssignmentID, &sub.StudentID, &sub.ContentText, &sub.FileURL, &sub.SubmittedAt, &sub.Status,
	)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *repository) Update(ctx context.Context, submission *Submission) error {
	query := `
		UPDATE submissions
		SET content_text = $1, file_url = $2, status = $3, submitted_at = now()
		WHERE id = $4`

	_, err := r.db.Exec(ctx, query,
		submission.ContentText, submission.FileURL, submission.Status, submission.ID,
	)
	return err
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM submissions WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *repository) ListByAssignment(ctx context.Context, assignmentID uuid.UUID, groupID *uuid.UUID, limit, offset int) ([]SubmissionWithDetails, int64, error) {
	var total int64
	var rows pgx.Rows
	var err error

	if groupID != nil {
		// Count with group filter
		countQuery := `
			SELECT COUNT(*)
			FROM submissions s
			JOIN students st ON s.student_id = st.user_id
			WHERE s.assignment_id = $1 AND st.group_id = $2`
		if err := r.db.QueryRow(ctx, countQuery, assignmentID, *groupID).Scan(&total); err != nil {
			return nil, 0, err
		}

		// Query with group filter
		query := `
			SELECT
				s.id, s.assignment_id, s.student_id, s.content_text, s.file_url, s.submitted_at, s.status,
				a.title AS assignment_title, a.max_points, a.due_at,
				c.title AS course_title,
				st.first_name AS student_first_name, st.last_name AS student_last_name,
				u.email AS student_email,
				g.score AS grade_score, g.feedback AS grade_feedback
			FROM submissions s
			JOIN assignments a ON s.assignment_id = a.id
			JOIN courses c ON a.course_id = c.id
			JOIN students st ON s.student_id = st.user_id
			JOIN users u ON st.user_id = u.id
			LEFT JOIN grades g ON s.id = g.submission_id
			WHERE s.assignment_id = $1 AND st.group_id = $2
			ORDER BY s.submitted_at DESC
			LIMIT $3 OFFSET $4`

		rows, err = r.db.Query(ctx, query, assignmentID, *groupID, limit, offset)
	} else {
		// Original query without group filter
		countQuery := `SELECT COUNT(*) FROM submissions WHERE assignment_id = $1`
		if err := r.db.QueryRow(ctx, countQuery, assignmentID).Scan(&total); err != nil {
			return nil, 0, err
		}

		query := `
			SELECT
				s.id, s.assignment_id, s.student_id, s.content_text, s.file_url, s.submitted_at, s.status,
				a.title AS assignment_title, a.max_points, a.due_at,
				c.title AS course_title,
				st.first_name AS student_first_name, st.last_name AS student_last_name,
				u.email AS student_email,
				g.score AS grade_score, g.feedback AS grade_feedback
			FROM submissions s
			JOIN assignments a ON s.assignment_id = a.id
			JOIN courses c ON a.course_id = c.id
			JOIN students st ON s.student_id = st.user_id
			JOIN users u ON st.user_id = u.id
			LEFT JOIN grades g ON s.id = g.submission_id
			WHERE s.assignment_id = $1
			ORDER BY s.submitted_at DESC
			LIMIT $2 OFFSET $3`

		rows, err = r.db.Query(ctx, query, assignmentID, limit, offset)
	}

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanSubmissions(rows, total)
}

func (r *repository) ListByStudent(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]SubmissionWithDetails, int64, error) {
	countQuery := `SELECT COUNT(*) FROM submissions WHERE student_id = $1`
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, studentID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT
			s.id, s.assignment_id, s.student_id, s.content_text, s.file_url, s.submitted_at, s.status,
			a.title AS assignment_title, a.max_points, a.due_at,
			c.title AS course_title,
			st.first_name AS student_first_name, st.last_name AS student_last_name,
			u.email AS student_email,
			g.score AS grade_score, g.feedback AS grade_feedback
		FROM submissions s
		JOIN assignments a ON s.assignment_id = a.id
		JOIN courses c ON a.course_id = c.id
		JOIN students st ON s.student_id = st.user_id
		JOIN users u ON st.user_id = u.id
		LEFT JOIN grades g ON s.id = g.submission_id
		WHERE s.student_id = $1
		ORDER BY s.submitted_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, studentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanSubmissions(rows, total)
}

func (r *repository) scanSubmissions(rows pgx.Rows, total int64) ([]SubmissionWithDetails, int64, error) {
	var submissions []SubmissionWithDetails
	for rows.Next() {
		var sub SubmissionWithDetails
		if err := rows.Scan(
			&sub.ID, &sub.AssignmentID, &sub.StudentID, &sub.ContentText, &sub.FileURL, &sub.SubmittedAt, &sub.Status,
			&sub.AssignmentTitle, &sub.MaxPoints, &sub.DueAt,
			&sub.CourseTitle,
			&sub.StudentFirstName, &sub.StudentLastName, &sub.StudentEmail,
			&sub.GradeScore, &sub.GradeFeedback,
		); err != nil {
			return nil, 0, err
		}
		submissions = append(submissions, sub)
	}
	return submissions, total, nil
}
