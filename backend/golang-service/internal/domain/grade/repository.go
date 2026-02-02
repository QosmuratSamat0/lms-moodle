package grade

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the grade repository interface
type Repository interface {
	Create(ctx context.Context, grade *Grade) error
	GetByID(ctx context.Context, id uuid.UUID) (*GradeWithDetails, error)
	GetBySubmissionID(ctx context.Context, submissionID uuid.UUID) (*Grade, error)
	Update(ctx context.Context, grade *Grade) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByStudent(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]GradeWithDetails, int64, error)
	ListByCourse(ctx context.Context, courseID uuid.UUID, limit, offset int) ([]GradeWithDetails, int64, error)
	ListByAssignment(ctx context.Context, assignmentID uuid.UUID, limit, offset int) ([]GradeWithDetails, int64, error)
	GetStudentCourseSummary(ctx context.Context, studentID, courseID uuid.UUID) (*StudentGradeSummary, error)
}

type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new grade repository
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, grade *Grade) error {
	query := `
		INSERT INTO grades (submission_id, graded_by_teacher_id, score, feedback)
		VALUES ($1, $2, $3, $4)
		RETURNING id, graded_at`

	return r.db.QueryRow(ctx, query,
		grade.SubmissionID,
		grade.GradedByTeacherID,
		grade.Score,
		grade.Feedback,
	).Scan(&grade.ID, &grade.GradedAt)
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*GradeWithDetails, error) {
	query := `
		SELECT 
			g.id, g.submission_id, g.graded_by_teacher_id, g.score, g.feedback, g.graded_at,
			s.student_id, s.assignment_id,
			a.title AS assignment_title, a.max_points,
			c.title AS course_title,
			st.first_name AS student_first_name, st.last_name AS student_last_name,
			u.email AS student_email,
			t.first_name AS teacher_first_name, t.last_name AS teacher_last_name
		FROM grades g
		JOIN submissions s ON g.submission_id = s.id
		JOIN assignments a ON s.assignment_id = a.id
		JOIN courses c ON a.course_id = c.id
		JOIN students st ON s.student_id = st.user_id
		JOIN users u ON st.user_id = u.id
		LEFT JOIN teachers t ON g.graded_by_teacher_id = t.user_id
		WHERE g.id = $1`

	var g GradeWithDetails
	err := r.db.QueryRow(ctx, query, id).Scan(
		&g.ID, &g.SubmissionID, &g.GradedByTeacherID, &g.Score, &g.Feedback, &g.GradedAt,
		&g.StudentID, &g.AssignmentID,
		&g.AssignmentTitle, &g.MaxPoints,
		&g.CourseTitle,
		&g.StudentFirstName, &g.StudentLastName, &g.StudentEmail,
		&g.TeacherFirstName, &g.TeacherLastName,
	)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *repository) GetBySubmissionID(ctx context.Context, submissionID uuid.UUID) (*Grade, error) {
	query := `
		SELECT id, submission_id, graded_by_teacher_id, score, feedback, graded_at
		FROM grades
		WHERE submission_id = $1`

	var g Grade
	err := r.db.QueryRow(ctx, query, submissionID).Scan(
		&g.ID, &g.SubmissionID, &g.GradedByTeacherID, &g.Score, &g.Feedback, &g.GradedAt,
	)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *repository) Update(ctx context.Context, grade *Grade) error {
	query := `
		UPDATE grades 
		SET score = $1, feedback = $2, graded_by_teacher_id = $3, graded_at = now()
		WHERE id = $4`

	_, err := r.db.Exec(ctx, query,
		grade.Score, grade.Feedback, grade.GradedByTeacherID, grade.ID,
	)
	return err
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM grades WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *repository) ListByStudent(ctx context.Context, studentID uuid.UUID, limit, offset int) ([]GradeWithDetails, int64, error) {
	countQuery := `
		SELECT COUNT(*) 
		FROM grades g
		JOIN submissions s ON g.submission_id = s.id
		WHERE s.student_id = $1`

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, studentID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			g.id, g.submission_id, g.graded_by_teacher_id, g.score, g.feedback, g.graded_at,
			s.student_id, s.assignment_id,
			a.title AS assignment_title, a.max_points,
			c.title AS course_title,
			st.first_name AS student_first_name, st.last_name AS student_last_name,
			u.email AS student_email,
			t.first_name AS teacher_first_name, t.last_name AS teacher_last_name
		FROM grades g
		JOIN submissions s ON g.submission_id = s.id
		JOIN assignments a ON s.assignment_id = a.id
		JOIN courses c ON a.course_id = c.id
		JOIN students st ON s.student_id = st.user_id
		JOIN users u ON st.user_id = u.id
		LEFT JOIN teachers t ON g.graded_by_teacher_id = t.user_id
		WHERE s.student_id = $1
		ORDER BY g.graded_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, studentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanGrades(rows, total)
}

func (r *repository) ListByCourse(ctx context.Context, courseID uuid.UUID, limit, offset int) ([]GradeWithDetails, int64, error) {
	countQuery := `
		SELECT COUNT(*) 
		FROM grades g
		JOIN submissions s ON g.submission_id = s.id
		JOIN assignments a ON s.assignment_id = a.id
		WHERE a.course_id = $1`

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, courseID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			g.id, g.submission_id, g.graded_by_teacher_id, g.score, g.feedback, g.graded_at,
			s.student_id, s.assignment_id,
			a.title AS assignment_title, a.max_points,
			c.title AS course_title,
			st.first_name AS student_first_name, st.last_name AS student_last_name,
			u.email AS student_email,
			t.first_name AS teacher_first_name, t.last_name AS teacher_last_name
		FROM grades g
		JOIN submissions s ON g.submission_id = s.id
		JOIN assignments a ON s.assignment_id = a.id
		JOIN courses c ON a.course_id = c.id
		JOIN students st ON s.student_id = st.user_id
		JOIN users u ON st.user_id = u.id
		LEFT JOIN teachers t ON g.graded_by_teacher_id = t.user_id
		WHERE a.course_id = $1
		ORDER BY g.graded_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, courseID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanGrades(rows, total)
}

func (r *repository) ListByAssignment(ctx context.Context, assignmentID uuid.UUID, limit, offset int) ([]GradeWithDetails, int64, error) {
	countQuery := `
		SELECT COUNT(*) 
		FROM grades g
		JOIN submissions s ON g.submission_id = s.id
		WHERE s.assignment_id = $1`

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, assignmentID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			g.id, g.submission_id, g.graded_by_teacher_id, g.score, g.feedback, g.graded_at,
			s.student_id, s.assignment_id,
			a.title AS assignment_title, a.max_points,
			c.title AS course_title,
			st.first_name AS student_first_name, st.last_name AS student_last_name,
			u.email AS student_email,
			t.first_name AS teacher_first_name, t.last_name AS teacher_last_name
		FROM grades g
		JOIN submissions s ON g.submission_id = s.id
		JOIN assignments a ON s.assignment_id = a.id
		JOIN courses c ON a.course_id = c.id
		JOIN students st ON s.student_id = st.user_id
		JOIN users u ON st.user_id = u.id
		LEFT JOIN teachers t ON g.graded_by_teacher_id = t.user_id
		WHERE s.assignment_id = $1
		ORDER BY g.graded_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, assignmentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanGrades(rows, total)
}

func (r *repository) GetStudentCourseSummary(ctx context.Context, studentID, courseID uuid.UUID) (*StudentGradeSummary, error) {
	query := `
		SELECT 
			c.id AS course_id,
			c.title AS course_title,
			COUNT(DISTINCT a.id) AS total_assignments,
			COUNT(DISTINCT g.id) AS graded_assignments,
			COALESCE(SUM(a.max_points), 0) AS total_points,
			COALESCE(SUM(g.score), 0) AS earned_points
		FROM courses c
		LEFT JOIN assignments a ON c.id = a.course_id
		LEFT JOIN submissions s ON a.id = s.assignment_id AND s.student_id = $1
		LEFT JOIN grades g ON s.id = g.submission_id
		WHERE c.id = $2
		GROUP BY c.id, c.title`

	var summary StudentGradeSummary
	var courseIDVal, courseTitle string
	var totalAssignments, gradedAssignments, totalPoints, earnedPoints int

	err := r.db.QueryRow(ctx, query, studentID, courseID).Scan(
		&courseIDVal, &courseTitle, &totalAssignments, &gradedAssignments, &totalPoints, &earnedPoints,
	)
	if err != nil {
		return nil, err
	}

	summary.StudentID = studentID.String()
	summary.CourseID = courseIDVal
	summary.CourseTitle = courseTitle
	summary.TotalAssignments = totalAssignments
	summary.GradedAssignments = gradedAssignments
	summary.TotalPoints = totalPoints
	summary.EarnedPoints = earnedPoints
	if totalPoints > 0 {
		summary.AveragePercentage = float64(earnedPoints) / float64(totalPoints) * 100
	}

	return &summary, nil
}

func (r *repository) scanGrades(rows pgx.Rows, total int64) ([]GradeWithDetails, int64, error) {
	var grades []GradeWithDetails
	for rows.Next() {
		var g GradeWithDetails
		if err := rows.Scan(
			&g.ID, &g.SubmissionID, &g.GradedByTeacherID, &g.Score, &g.Feedback, &g.GradedAt,
			&g.StudentID, &g.AssignmentID,
			&g.AssignmentTitle, &g.MaxPoints,
			&g.CourseTitle,
			&g.StudentFirstName, &g.StudentLastName, &g.StudentEmail,
			&g.TeacherFirstName, &g.TeacherLastName,
		); err != nil {
			return nil, 0, err
		}
		grades = append(grades, g)
	}
	return grades, total, nil
}
