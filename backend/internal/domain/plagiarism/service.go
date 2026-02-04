package plagiarism

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/ap1-final-mini-moodle/internal/shared/errorx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const DefaultSuspiciousThreshold = 30.0 // 30% similarity

// Service defines the plagiarism service interface
type Service interface {
	CheckSubmission(ctx context.Context, submissionID uuid.UUID, text string) (*Report, error)
	GetReport(ctx context.Context, id uuid.UUID) (*ReportWithSubmission, error)
	GetReportBySubmission(ctx context.Context, submissionID uuid.UUID) (*Report, error)
	ListByAssignment(ctx context.Context, assignmentID uuid.UUID, page, limit int) (*ReportListResponse, error)
	ListSuspicious(ctx context.Context, threshold *float64, page, limit int) (*ReportListResponse, error)
	Recheck(ctx context.Context, submissionID uuid.UUID) (*Report, error)
}

type service struct {
	repo     Repository
	detector Detector
}

// NewService creates a new plagiarism service
func NewService(repo Repository, detector Detector) Service {
	return &service{
		repo:     repo,
		detector: detector,
	}
}

// ReportListResponse represents a paginated list of reports
type ReportListResponse struct {
	Reports    []ReportWithSubmission `json:"reports"`
	TotalCount int                    `json:"total_count"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
}

func (s *service) CheckSubmission(ctx context.Context, submissionID uuid.UUID, text string) (*Report, error) {
	if text == "" {
		return nil, errorx.NewValidationError("empty submission text")
	}

	// Check against corpus (placeholder - would query other submissions)
	startTime := time.Now()
	details := &ReportDetails{
		Matches:      []Match{},
		Sources:      []Source{},
		AnalysisTime: time.Since(startTime).Milliseconds(),
	}

	// For demo, use simple self-check score of 0
	score := 0.0

	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return nil, errorx.Wrap(err, "marshal details")
	}

	report := &Report{
		SubmissionID: submissionID,
		Score:        score,
		Details:      detailsJSON,
	}

	// Check if report already exists
	existing, err := s.repo.GetBySubmissionID(ctx, submissionID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, errorx.Wrap(err, "check existing report")
	}

	if existing != nil {
		// Update existing report
		report.ID = existing.ID
		if err := s.repo.Update(ctx, report); err != nil {
			return nil, errorx.Wrap(err, "update report")
		}
		return report, nil
	}

	// Create new report
	if err := s.repo.Create(ctx, report); err != nil {
		return nil, errorx.Wrap(err, "create report")
	}

	return report, nil
}

func (s *service) GetReport(ctx context.Context, id uuid.UUID) (*ReportWithSubmission, error) {
	report, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("plagiarism report")
		}
		return nil, errorx.Wrap(err, "get report")
	}

	return report, nil
}

func (s *service) GetReportBySubmission(ctx context.Context, submissionID uuid.UUID) (*Report, error) {
	report, err := s.repo.GetBySubmissionID(ctx, submissionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("plagiarism report")
		}
		return nil, errorx.Wrap(err, "get report by submission")
	}

	return report, nil
}

func (s *service) ListByAssignment(ctx context.Context, assignmentID uuid.UUID, page, limit int) (*ReportListResponse, error) {
	reports, total, err := s.repo.ListByAssignment(ctx, assignmentID, page, limit)
	if err != nil {
		return nil, errorx.Wrap(err, "list reports by assignment")
	}

	return &ReportListResponse{
		Reports:    reports,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
	}, nil
}

func (s *service) ListSuspicious(ctx context.Context, threshold *float64, page, limit int) (*ReportListResponse, error) {
	th := DefaultSuspiciousThreshold
	if threshold != nil {
		th = *threshold
	}

	reports, total, err := s.repo.ListSuspicious(ctx, th, page, limit)
	if err != nil {
		return nil, errorx.Wrap(err, "list suspicious reports")
	}

	return &ReportListResponse{
		Reports:    reports,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
	}, nil
}

func (s *service) Recheck(ctx context.Context, submissionID uuid.UUID) (*Report, error) {
	// In real implementation, would fetch submission text and recheck
	// For now, just get existing report
	existing, err := s.repo.GetBySubmissionID(ctx, submissionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorx.NewNotFoundError("no existing report to recheck")
		}
		return nil, errorx.Wrap(err, "get existing report")
	}

	return existing, nil
}
