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
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO grades (id, submission_id, score, feedback, graded_by, graded_at, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		g.ID, g.SubmissionID, g.Score, g.Feedback, g.GradedBy, g.GradedAt, g.CreatedAt, g.UpdatedAt)
	return err
}

func (r *PostgresRepository) GetByID(id string) (*grade.Grade, error) {
	g := &grade.Grade{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, submission_id, score, feedback, graded_by, graded_at, created_at, updated_at
		 FROM grades WHERE id = $1`, id).
		Scan(&g.ID, &g.SubmissionID, &g.Score, &g.Feedback, &g.GradedBy, &g.GradedAt, &g.CreatedAt, &g.UpdatedAt)
	return g, err
}

func (r *PostgresRepository) GetBySubmission(submissionID string) (*grade.Grade, error) {
	g := &grade.Grade{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, submission_id, score, feedback, graded_by, graded_at, created_at, updated_at
		 FROM grades WHERE submission_id = $1`, submissionID).
		Scan(&g.ID, &g.SubmissionID, &g.Score, &g.Feedback, &g.GradedBy, &g.GradedAt, &g.CreatedAt, &g.UpdatedAt)
	return g, err
}

func (r *PostgresRepository) Update(g *grade.Grade) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE grades SET score=$1, feedback=$2, updated_at=$3 WHERE id=$4`,
		g.Score, g.Feedback, g.UpdatedAt, g.ID)
	return err
}

func (r *PostgresRepository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM grades WHERE id = $1", id)
	return err
}
