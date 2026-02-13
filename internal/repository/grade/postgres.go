package grade

import (
	"context"

	"github.com/ap1-final-mini-moodle/internal/domain/grade"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) grade.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(g *grade.Grade) error {
	// Auto-ensure teacher record exists (migration 000010 changed teachers PK to separate id)
	if g.GradedBy != "" {
		_, _ = r.db.Exec(context.Background(),
			`INSERT INTO teachers (user_id) VALUES ($1) ON CONFLICT (user_id) DO NOTHING`,
			g.GradedBy)
	}
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO grades (id, submission_id, score, feedback, graded_by_teacher_id, graded_at)
		 VALUES ($1, $2, $3, $4, (SELECT id FROM teachers WHERE user_id = $5), $6)`,
		g.ID, g.SubmissionID, g.Score, g.Feedback, g.GradedBy, g.GradedAt)
	return err
}

func (r *PostgresRepository) GetByID(id string) (*grade.Grade, error) {
	g := &grade.Grade{}
	err := r.db.QueryRow(context.Background(),
		`SELECT g.id, g.submission_id, g.score, COALESCE(g.feedback, ''), COALESCE(t.user_id::text, ''), g.graded_at, g.graded_at, g.graded_at
		 FROM grades g
		 LEFT JOIN teachers t ON g.graded_by_teacher_id = t.id
		 WHERE g.id = $1`, id).
		Scan(&g.ID, &g.SubmissionID, &g.Score, &g.Feedback, &g.GradedBy, &g.GradedAt, &g.CreatedAt, &g.UpdatedAt)
	return g, err
}

func (r *PostgresRepository) GetBySubmission(submissionID string) (*grade.Grade, error) {
	g := &grade.Grade{}
	err := r.db.QueryRow(context.Background(),
		`SELECT g.id, g.submission_id, g.score, COALESCE(g.feedback, ''), COALESCE(t.user_id::text, ''), g.graded_at, g.graded_at, g.graded_at
		 FROM grades g
		 LEFT JOIN teachers t ON g.graded_by_teacher_id = t.id
		 WHERE g.submission_id = $1`, submissionID).
		Scan(&g.ID, &g.SubmissionID, &g.Score, &g.Feedback, &g.GradedBy, &g.GradedAt, &g.CreatedAt, &g.UpdatedAt)
	return g, err
}

func (r *PostgresRepository) Update(g *grade.Grade) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE grades SET score=$1, feedback=$2, graded_at=$3 WHERE id=$4`,
		g.Score, g.Feedback, g.GradedAt, g.ID)
	return err
}

func (r *PostgresRepository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM grades WHERE id = $1", id)
	return err
}
