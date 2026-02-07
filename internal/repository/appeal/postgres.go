package appeal

import (
	"context"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/appeal"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) appeal.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, a *appeal.GradeAppeal) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	a.Status = appeal.AppealPending
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx,
		`INSERT INTO grade_appeals (id, grade_id, student_id, reason, evidence, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		a.ID, a.GradeID, a.StudentID, a.Reason, a.Evidence, a.Status, a.CreatedAt, a.UpdatedAt)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*appeal.GradeAppeal, error) {
	a := &appeal.GradeAppeal{}
	err := r.db.QueryRow(ctx,
		`SELECT id, grade_id, student_id, reason, evidence, status, teacher_id, response, new_score, resolved_at, created_at, updated_at
		 FROM grade_appeals WHERE id = $1`, id).
		Scan(&a.ID, &a.GradeID, &a.StudentID, &a.Reason, &a.Evidence, &a.Status, &a.TeacherID, &a.Response, &a.NewScore, &a.ResolvedAt, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (r *PostgresRepository) GetByGradeID(ctx context.Context, gradeID string) (*appeal.GradeAppeal, error) {
	a := &appeal.GradeAppeal{}
	err := r.db.QueryRow(ctx,
		`SELECT id, grade_id, student_id, reason, evidence, status, teacher_id, response, new_score, resolved_at, created_at, updated_at
		 FROM grade_appeals WHERE grade_id = $1 ORDER BY created_at DESC LIMIT 1`, gradeID).
		Scan(&a.ID, &a.GradeID, &a.StudentID, &a.Reason, &a.Evidence, &a.Status, &a.TeacherID, &a.Response, &a.NewScore, &a.ResolvedAt, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (r *PostgresRepository) GetByStudentID(ctx context.Context, studentID string, status *appeal.AppealStatus, limit, offset int) ([]*appeal.AppealWithDetails, int64, error) {
	countQuery := `SELECT COUNT(*) FROM grade_appeals WHERE student_id = $1`
	listQuery := `
		SELECT ga.id, ga.grade_id, ga.student_id, ga.reason, ga.evidence, ga.status, ga.teacher_id, ga.response, ga.new_score, ga.resolved_at, ga.created_at, ga.updated_at,
		       u.first_name || ' ' || u.last_name as student_name, u.email as student_email,
		       c.title as course_name, a.title as assignment_title,
		       g.score as original_score, a.max_points
		FROM grade_appeals ga
		JOIN grades g ON ga.grade_id = g.id
		JOIN submissions s ON g.submission_id = s.id
		JOIN assignments a ON s.assignment_id = a.id
		JOIN courses c ON a.course_id = c.id
		JOIN users u ON ga.student_id = u.id
		WHERE ga.student_id = $1`

	args := []interface{}{studentID}
	argIdx := 2

	if status != nil {
		countQuery += ` AND status = $2`
		listQuery += ` AND ga.status = $2`
		args = append(args, *status)
		argIdx++
	}

	var total int64
	if status != nil {
		r.db.QueryRow(ctx, countQuery, studentID, *status).Scan(&total)
	} else {
		r.db.QueryRow(ctx, countQuery, studentID).Scan(&total)
	}

	listQuery += ` ORDER BY ga.created_at DESC LIMIT $` + string(rune('0'+argIdx)) + ` OFFSET $` + string(rune('0'+argIdx+1))
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var appeals []*appeal.AppealWithDetails
	for rows.Next() {
		a := &appeal.AppealWithDetails{}
		if err := rows.Scan(
			&a.ID, &a.GradeID, &a.StudentID, &a.Reason, &a.Evidence, &a.Status, &a.TeacherID, &a.Response, &a.NewScore, &a.ResolvedAt, &a.CreatedAt, &a.UpdatedAt,
			&a.StudentName, &a.StudentEmail, &a.CourseName, &a.AssignmentTitle, &a.OriginalScore, &a.MaxPoints,
		); err != nil {
			return nil, 0, err
		}
		appeals = append(appeals, a)
	}
	return appeals, total, rows.Err()
}

func (r *PostgresRepository) GetByTeacherCourses(ctx context.Context, teacherID string, status *appeal.AppealStatus, limit, offset int) ([]*appeal.AppealWithDetails, int64, error) {
	countQuery := `
		SELECT COUNT(*)
		FROM grade_appeals ga
		JOIN grades g ON ga.grade_id = g.id
		JOIN submissions s ON g.submission_id = s.id
		JOIN assignments a ON s.assignment_id = a.id
		JOIN courses c ON a.course_id = c.id
		WHERE c.teacher_id = $1`

	listQuery := `
		SELECT ga.id, ga.grade_id, ga.student_id, ga.reason, ga.evidence, ga.status, ga.teacher_id, ga.response, ga.new_score, ga.resolved_at, ga.created_at, ga.updated_at,
		       u.first_name || ' ' || u.last_name as student_name, u.email as student_email,
		       c.title as course_name, a.title as assignment_title,
		       g.score as original_score, a.max_points
		FROM grade_appeals ga
		JOIN grades g ON ga.grade_id = g.id
		JOIN submissions s ON g.submission_id = s.id
		JOIN assignments a ON s.assignment_id = a.id
		JOIN courses c ON a.course_id = c.id
		JOIN users u ON ga.student_id = u.id
		WHERE c.teacher_id = $1`

	args := []interface{}{teacherID}

	if status != nil {
		countQuery += ` AND ga.status = $2`
		listQuery += ` AND ga.status = $2`
		args = append(args, *status)
	}

	var total int64
	if status != nil {
		r.db.QueryRow(ctx, countQuery, teacherID, *status).Scan(&total)
	} else {
		r.db.QueryRow(ctx, countQuery, teacherID).Scan(&total)
	}

	listQuery += ` ORDER BY ga.created_at DESC LIMIT $2 OFFSET $3`
	if status != nil {
		listQuery = `
		SELECT ga.id, ga.grade_id, ga.student_id, ga.reason, ga.evidence, ga.status, ga.teacher_id, ga.response, ga.new_score, ga.resolved_at, ga.created_at, ga.updated_at,
		       u.first_name || ' ' || u.last_name as student_name, u.email as student_email,
		       c.title as course_name, a.title as assignment_title,
		       g.score as original_score, a.max_points
		FROM grade_appeals ga
		JOIN grades g ON ga.grade_id = g.id
		JOIN submissions s ON g.submission_id = s.id
		JOIN assignments a ON s.assignment_id = a.id
		JOIN courses c ON a.course_id = c.id
		JOIN users u ON ga.student_id = u.id
		WHERE c.teacher_id = $1 AND ga.status = $2
		ORDER BY ga.created_at DESC LIMIT $3 OFFSET $4`
		args = append(args, limit, offset)
	} else {
		args = append(args, limit, offset)
	}

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var appeals []*appeal.AppealWithDetails
	for rows.Next() {
		a := &appeal.AppealWithDetails{}
		if err := rows.Scan(
			&a.ID, &a.GradeID, &a.StudentID, &a.Reason, &a.Evidence, &a.Status, &a.TeacherID, &a.Response, &a.NewScore, &a.ResolvedAt, &a.CreatedAt, &a.UpdatedAt,
			&a.StudentName, &a.StudentEmail, &a.CourseName, &a.AssignmentTitle, &a.OriginalScore, &a.MaxPoints,
		); err != nil {
			return nil, 0, err
		}
		appeals = append(appeals, a)
	}
	return appeals, total, rows.Err()
}

func (r *PostgresRepository) Update(ctx context.Context, a *appeal.GradeAppeal) error {
	a.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx,
		`UPDATE grade_appeals SET status=$1, teacher_id=$2, response=$3, new_score=$4, resolved_at=$5, updated_at=$6 WHERE id=$7`,
		a.Status, a.TeacherID, a.Response, a.NewScore, a.ResolvedAt, a.UpdatedAt, a.ID)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM grade_appeals WHERE id = $1`, id)
	return err
}

func (r *PostgresRepository) HasPendingAppeal(ctx context.Context, gradeID string) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM grade_appeals WHERE grade_id = $1 AND status = 'pending'`, gradeID).Scan(&count)
	return count > 0, err
}
