package teacher

import (
	"context"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/teacher"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) teacher.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, teacher *teacher.Teacher) error {
	if teacher.ID == "" {
		teacher.ID = uuid.New().String()
	}
	teacher.CreatedAt = time.Now()
	teacher.UpdatedAt = time.Now()

	_, err := r.db.QueryRow(ctx, ``)
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*teacher.Teacher, error) {

}

func (r *PostgresRepository) GetByUserID(ctx context.Context, userID string) (*teacher.Teacher, error) {

}

func (r *PostgresRepository) GetByEmployeeID(ctx context.Context, employeeID string) (*teacher.Teacher, error) {

}

func (r *PostgresRepository) List(ctx context.Context, skip, take int) ([]*teacher.Teacher, error) {

}

func (r *PostgresRepository) Update(ctx context.Context, teacher *teacher.Teacher) error {

}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {

}

func (r *PostgresRepository) GetTeacherGourses(ctx context.Context, teacherID string) ([]teacher.TeacherCourse, error) {

}

func (r *PostgresRepository) GetTeacherGroups(ctx context.Context, teacherID string) ([]*teacher.TeacherGroup, error) {

}
