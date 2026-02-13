package submission

import (
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/submission"
	"github.com/google/uuid"
)

type Service struct {
	repo submission.Repository
}

func NewService(repo submission.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Submit(input *submission.CreateSubmissionInput) (*submission.Submission, error) {
	sub := &submission.Submission{
		ID:           uuid.New().String(),
		AssignmentID: input.AssignmentID,
		StudentID:    input.StudentID,
		ContentText:  input.ContentText,
		FileURL:      input.FileURL,
		SubmittedAt:  time.Now(),
		Status:       "submitted",
	}
	if err := s.repo.Create(sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *Service) GetByID(id string) (*submission.Submission, error) {
	return s.repo.GetByID(id)
}

func (s *Service) GetByAssignmentAndStudent(assignmentID, studentID string) (*submission.Submission, error) {
	return s.repo.GetByAssignmentAndStudent(assignmentID, studentID)
}

func (s *Service) ListByAssignment(assignmentID string, skip, take int) ([]*submission.Submission, error) {
	return s.repo.ListByAssignment(assignmentID, skip, take)
}

func (s *Service) ListByStudent(studentID string, skip, take int) ([]*submission.Submission, error) {
	return s.repo.ListByStudent(studentID, skip, take)
}

func (s *Service) Update(input *submission.UpdateSubmissionInput) (*submission.Submission, error) {
	// Get existing submission
	existing, err := s.repo.GetByID(input.ID)
	if err != nil {
		return nil, err
	}

	// Update fields
	existing.ContentText = input.ContentText
	existing.FileURL = input.FileURL

	// Save to database
	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
