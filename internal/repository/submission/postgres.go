package submission

import (
	"context"

	"github.com/ap1-final-mini-moodle/internal/domain/submission"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) submission.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(s *submission.Submission) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO submissions (id, assignment_id, student_id, content, file_url, submitted_at, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		s.ID, s.AssignmentID, s.StudentID, s.Content, s.FileURL, s.SubmittedAt, s.CreatedAt, s.UpdatedAt)
	return err
}

func (r *PostgresRepository) GetByID(id string) (*submission.Submission, error) {
	s := &submission.Submission{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, assignment_id, student_id, content, file_url, submitted_at, graded_at, created_at, updated_at
		 FROM submissions WHERE id = $1`, id).
		Scan(&s.ID, &s.AssignmentID, &s.StudentID, &s.Content, &s.FileURL, &s.SubmittedAt, &s.GradedAt, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *PostgresRepository) GetByAssignmentAndStudent(assignmentID, studentID string) (*submission.Submission, error) {
	s := &submission.Submission{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, assignment_id, student_id, content, file_url, submitted_at, graded_at, created_at, updated_at
		 FROM submissions WHERE assignment_id=$1 AND student_id=$2`, assignmentID, studentID).
		Scan(&s.ID, &s.AssignmentID, &s.StudentID, &s.Content, &s.FileURL, &s.SubmittedAt, &s.GradedAt, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *PostgresRepository) ListByAssignment(assignmentID string, skip, take int) ([]*submission.Submission, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, assignment_id, student_id, content, file_url, submitted_at, graded_at, created_at, updated_at
		 FROM submissions WHERE assignment_id=$1 OFFSET $2 LIMIT $3`, assignmentID, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []*submission.Submission
	for rows.Next() {
		s := &submission.Submission{}
		if err := rows.Scan(&s.ID, &s.AssignmentID, &s.StudentID, &s.Content, &s.FileURL, &s.SubmittedAt, &s.GradedAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		submissions = append(submissions, s)
	}
	return submissions, rows.Err()
}

func (r *PostgresRepository) ListByStudent(studentID string, skip, take int) ([]*submission.Submission, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, assignment_id, student_id, content, file_url, submitted_at, graded_at, created_at, updated_at
		 FROM submissions WHERE student_id=$1 OFFSET $2 LIMIT $3`, studentID, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []*submission.Submission
	for rows.Next() {
		s := &submission.Submission{}
		if err := rows.Scan(&s.ID, &s.AssignmentID, &s.StudentID, &s.Content, &s.FileURL, &s.SubmittedAt, &s.GradedAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		submissions = append(submissions, s)
	}
	return submissions, rows.Err()
}

func (r *PostgresRepository) Update(s *submission.Submission) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE submissions SET content=$1, file_url=$2, graded_at=$3, updated_at=$4 WHERE id=$5`,
		s.Content, s.FileURL, s.GradedAt, s.UpdatedAt, s.ID)
	return err
}

func (r *PostgresRepository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM submissions WHERE id = $1", id)
	return err
}
