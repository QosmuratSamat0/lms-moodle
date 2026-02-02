package plagiarism

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the plagiarism repository interface
type Repository interface {
	Create(ctx context.Context, report *Report) error
	GetByID(ctx context.Context, id uuid.UUID) (*ReportWithSubmission, error)
	GetBySubmissionID(ctx context.Context, submissionID uuid.UUID) (*Report, error)
	ListByAssignment(ctx context.Context, assignmentID uuid.UUID, page, limit int) ([]ReportWithSubmission, int, error)
	ListSuspicious(ctx context.Context, threshold float64, page, limit int) ([]ReportWithSubmission, int, error)
	Update(ctx context.Context, report *Report) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new plagiarism repository
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, report *Report) error {
	query := `
		INSERT INTO plagiarism_reports (submission_id, score, details)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	return r.db.QueryRow(ctx, query,
		report.SubmissionID, report.Score, report.Details,
	).Scan(&report.ID, &report.CreatedAt)
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*ReportWithSubmission, error) {
	query := `
		SELECT pr.id, pr.submission_id, pr.score, pr.details, pr.created_at,
		       s.student_id, s.assignment_id, s.content_text,
		       u.email as student_email
		FROM plagiarism_reports pr
		JOIN submissions s ON s.id = pr.submission_id
		JOIN users u ON u.id = s.student_id
		WHERE pr.id = $1`

	var report ReportWithSubmission
	err := r.db.QueryRow(ctx, query, id).Scan(
		&report.ID, &report.SubmissionID, &report.Score, &report.Details, &report.CreatedAt,
		&report.StudentID, &report.AssignmentID, &report.ContentText,
		&report.StudentEmail,
	)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (r *repository) GetBySubmissionID(ctx context.Context, submissionID uuid.UUID) (*Report, error) {
	query := `
		SELECT id, submission_id, score, details, created_at
		FROM plagiarism_reports
		WHERE submission_id = $1`

	var report Report
	err := r.db.QueryRow(ctx, query, submissionID).Scan(
		&report.ID, &report.SubmissionID, &report.Score, &report.Details, &report.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (r *repository) ListByAssignment(ctx context.Context, assignmentID uuid.UUID, page, limit int) ([]ReportWithSubmission, int, error) {
	countQuery := `
		SELECT COUNT(*) FROM plagiarism_reports pr
		JOIN submissions s ON s.id = pr.submission_id
		WHERE s.assignment_id = $1`

	var total int
	if err := r.db.QueryRow(ctx, countQuery, assignmentID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT pr.id, pr.submission_id, pr.score, pr.details, pr.created_at,
		       s.student_id, s.assignment_id, s.content_text,
		       u.email as student_email
		FROM plagiarism_reports pr
		JOIN submissions s ON s.id = pr.submission_id
		JOIN users u ON u.id = s.student_id
		WHERE s.assignment_id = $1
		ORDER BY pr.score DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, assignmentID, limit, (page-1)*limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reports []ReportWithSubmission
	for rows.Next() {
		var report ReportWithSubmission
		if err := rows.Scan(
			&report.ID, &report.SubmissionID, &report.Score, &report.Details, &report.CreatedAt,
			&report.StudentID, &report.AssignmentID, &report.ContentText,
			&report.StudentEmail,
		); err != nil {
			return nil, 0, err
		}
		reports = append(reports, report)
	}

	return reports, total, rows.Err()
}

func (r *repository) ListSuspicious(ctx context.Context, threshold float64, page, limit int) ([]ReportWithSubmission, int, error) {
	countQuery := `SELECT COUNT(*) FROM plagiarism_reports WHERE score >= $1`

	var total int
	if err := r.db.QueryRow(ctx, countQuery, threshold).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT pr.id, pr.submission_id, pr.score, pr.details, pr.created_at,
		       s.student_id, s.assignment_id, s.content_text,
		       u.email as student_email
		FROM plagiarism_reports pr
		JOIN submissions s ON s.id = pr.submission_id
		JOIN users u ON u.id = s.student_id
		WHERE pr.score >= $1
		ORDER BY pr.score DESC, pr.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, threshold, limit, (page-1)*limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reports []ReportWithSubmission
	for rows.Next() {
		var report ReportWithSubmission
		if err := rows.Scan(
			&report.ID, &report.SubmissionID, &report.Score, &report.Details, &report.CreatedAt,
			&report.StudentID, &report.AssignmentID, &report.ContentText,
			&report.StudentEmail,
		); err != nil {
			return nil, 0, err
		}
		reports = append(reports, report)
	}

	return reports, total, rows.Err()
}

func (r *repository) Update(ctx context.Context, report *Report) error {
	query := `
		UPDATE plagiarism_reports
		SET score = $2, details = $3
		WHERE id = $1`

	_, err := r.db.Exec(ctx, query, report.ID, report.Score, report.Details)
	return err
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM plagiarism_reports WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// Helper to convert details to JSON
func DetailsToJSON(details *ReportDetails) (json.RawMessage, error) {
	return json.Marshal(details)
}

// Helper to parse details from JSON
func JSONToDetails(data json.RawMessage) (*ReportDetails, error) {
	var details ReportDetails
	if err := json.Unmarshal(data, &details); err != nil {
		return nil, err
	}
	return &details, nil
}
