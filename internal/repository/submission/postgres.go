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
		`INSERT INTO submissions (id, assignment_id, student_id, content_text, file_url, submitted_at, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		s.ID, s.AssignmentID, s.StudentID, s.ContentText, s.FileURL, s.SubmittedAt, s.Status)
	return err
}

func (r *PostgresRepository) GetByID(id string) (*submission.Submission, error) {
	s := &submission.Submission{}
	err := r.db.QueryRow(context.Background(),
		`SELECT s.id, s.assignment_id, s.student_id, COALESCE(s.content_text, ''), s.file_url, s.submitted_at,
		        CASE WHEN g.id IS NOT NULL THEN 'graded' ELSE COALESCE(s.status, 'submitted') END,
		        COALESCE(st.first_name, ''), COALESCE(st.last_name, ''), COALESCE(u.email, ''),
		        g.score, g.feedback, g.graded_at,
		        gt.first_name, gt.last_name
		 FROM submissions s
		 LEFT JOIN students st ON s.student_id = st.user_id
		 LEFT JOIN users u ON s.student_id = u.id
		 LEFT JOIN grades g ON g.submission_id = s.id
		 LEFT JOIN teachers gt ON g.graded_by_teacher_id = gt.id
		 WHERE s.id = $1`, id).
		Scan(&s.ID, &s.AssignmentID, &s.StudentID, &s.ContentText, &s.FileURL, &s.SubmittedAt, &s.Status,
			&s.StudentFirstName, &s.StudentLastName, &s.StudentEmail,
			&s.GradeScore, &s.GradeFeedback, &s.GradedAt,
			&s.GraderFirstName, &s.GraderLastName)
	return s, err
}

func (r *PostgresRepository) GetByAssignmentAndStudent(assignmentID, studentID string) (*submission.Submission, error) {
	s := &submission.Submission{}
	err := r.db.QueryRow(context.Background(),
		`SELECT s.id, s.assignment_id, s.student_id, COALESCE(s.content_text, ''), s.file_url, s.submitted_at,
		        CASE WHEN g.id IS NOT NULL THEN 'graded' ELSE COALESCE(s.status, 'submitted') END,
		        COALESCE(st.first_name, ''), COALESCE(st.last_name, ''), COALESCE(u.email, ''),
		        g.score, g.feedback, g.graded_at,
		        gt.first_name, gt.last_name
		 FROM submissions s
		 LEFT JOIN students st ON s.student_id = st.user_id
		 LEFT JOIN users u ON s.student_id = u.id
		 LEFT JOIN grades g ON g.submission_id = s.id
		 LEFT JOIN teachers gt ON g.graded_by_teacher_id = gt.id
		 WHERE s.assignment_id=$1 AND s.student_id=$2`, assignmentID, studentID).
		Scan(&s.ID, &s.AssignmentID, &s.StudentID, &s.ContentText, &s.FileURL, &s.SubmittedAt, &s.Status,
			&s.StudentFirstName, &s.StudentLastName, &s.StudentEmail,
			&s.GradeScore, &s.GradeFeedback, &s.GradedAt,
			&s.GraderFirstName, &s.GraderLastName)
	return s, err
}

func (r *PostgresRepository) ListByAssignment(assignmentID string, skip, take int) ([]*submission.Submission, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT s.id, s.assignment_id, s.student_id, COALESCE(s.content_text, ''), s.file_url, s.submitted_at,
		        CASE WHEN g.id IS NOT NULL THEN 'graded' ELSE COALESCE(s.status, 'submitted') END,
		        COALESCE(st.first_name, ''), COALESCE(st.last_name, ''), COALESCE(u.email, ''),
		        g.score, g.feedback, g.graded_at,
		        gt.first_name, gt.last_name
		 FROM submissions s
		 LEFT JOIN students st ON s.student_id = st.user_id
		 LEFT JOIN users u ON s.student_id = u.id
		 LEFT JOIN grades g ON g.submission_id = s.id
		 LEFT JOIN teachers gt ON g.graded_by_teacher_id = gt.id
		 WHERE s.assignment_id=$1 ORDER BY s.submitted_at DESC OFFSET $2 LIMIT $3`, assignmentID, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []*submission.Submission
	for rows.Next() {
		s := &submission.Submission{}
		if err := rows.Scan(&s.ID, &s.AssignmentID, &s.StudentID, &s.ContentText, &s.FileURL, &s.SubmittedAt, &s.Status,
			&s.StudentFirstName, &s.StudentLastName, &s.StudentEmail,
			&s.GradeScore, &s.GradeFeedback, &s.GradedAt,
			&s.GraderFirstName, &s.GraderLastName); err != nil {
			return nil, err
		}
		submissions = append(submissions, s)
	}
	return submissions, rows.Err()
}

func (r *PostgresRepository) ListByStudent(studentID string, skip, take int) ([]*submission.Submission, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT s.id, s.assignment_id, s.student_id, COALESCE(s.content_text, ''), s.file_url, s.submitted_at,
		        CASE WHEN g.id IS NOT NULL THEN 'graded' ELSE COALESCE(s.status, 'submitted') END,
		        COALESCE(st.first_name, ''), COALESCE(st.last_name, ''), COALESCE(u.email, ''),
		        g.score, g.feedback, g.graded_at,
		        gt.first_name, gt.last_name
		 FROM submissions s
		 LEFT JOIN students st ON s.student_id = st.user_id
		 LEFT JOIN users u ON s.student_id = u.id
		 LEFT JOIN grades g ON g.submission_id = s.id
		 LEFT JOIN teachers gt ON g.graded_by_teacher_id = gt.id
		 WHERE s.student_id=$1 ORDER BY s.submitted_at DESC OFFSET $2 LIMIT $3`, studentID, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []*submission.Submission
	for rows.Next() {
		s := &submission.Submission{}
		if err := rows.Scan(&s.ID, &s.AssignmentID, &s.StudentID, &s.ContentText, &s.FileURL, &s.SubmittedAt, &s.Status,
			&s.StudentFirstName, &s.StudentLastName, &s.StudentEmail,
			&s.GradeScore, &s.GradeFeedback, &s.GradedAt,
			&s.GraderFirstName, &s.GraderLastName); err != nil {
			return nil, err
		}
		submissions = append(submissions, s)
	}
	return submissions, rows.Err()
}

func (r *PostgresRepository) Update(s *submission.Submission) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE submissions SET content_text=$1, file_url=$2, status=$3 WHERE id=$4`,
		s.ContentText, s.FileURL, s.Status, s.ID)
	return err
}

func (r *PostgresRepository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM submissions WHERE id = $1", id)
	return err
}
