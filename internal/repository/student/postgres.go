package student

import (
	"context"
	"strconv"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/student"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) student.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, s *student.Student) error {
	query := `
		INSERT INTO students (user_id, student_code, major, year, gpa, enrollment_status, group_id, admitted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	if s.AdmittedAt.IsZero() {
		s.AdmittedAt = time.Now()
	}
	if s.Status == "" {
		s.Status = "enrolled"
	}

	return r.db.QueryRow(ctx, query,
		s.UserID, s.StudentCode, s.Major, s.Year, s.GPA, s.Status, s.GroupID, s.AdmittedAt,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*student.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.student_code, s.major, s.year, s.gpa, s.enrollment_status,
			       COALESCE(s.group_id::text, ''), s.admitted_at, s.created_at, s.updated_at, u.first_name, u.last_name, u.email, COALESCE(g.name, '')
			FROM students s
			JOIN users u ON s.user_id = u.id
			LEFT JOIN groups g ON s.group_id = g.id
			WHERE s.id = $1`

	var s student.Student
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.UserID, &s.StudentCode, &s.Major, &s.Year, &s.GPA, &s.Status,
		&s.GroupID, &s.AdmittedAt, &s.CreatedAt, &s.UpdatedAt, &s.FirstName, &s.LastName, &s.Email, &s.GroupName,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *PostgresRepository) GetByUserID(ctx context.Context, userID string) (*student.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.student_code, s.major, s.year, s.gpa, s.enrollment_status,
			       COALESCE(s.group_id::text, ''), s.admitted_at, s.created_at, s.updated_at, u.first_name, u.last_name, u.email, COALESCE(g.name, '')
			FROM students s
			JOIN users u ON s.user_id = u.id
			LEFT JOIN groups g ON s.group_id = g.id
			WHERE s.user_id = $1`

	var s student.Student
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&s.ID, &s.UserID, &s.StudentCode, &s.Major, &s.Year, &s.GPA, &s.Status,
		&s.GroupID, &s.AdmittedAt, &s.CreatedAt, &s.UpdatedAt, &s.FirstName, &s.LastName, &s.Email, &s.GroupName,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *PostgresRepository) GetByStudentCode(ctx context.Context, code string) (*student.Student, error) {
	query := `
		SELECT s.id, s.user_id, s.student_code, s.major, s.year, s.gpa, s.enrollment_status,
			       COALESCE(s.group_id::text, ''), s.admitted_at, s.created_at, s.updated_at, u.first_name, u.last_name, u.email, COALESCE(g.name, '')
			FROM students s
			JOIN users u ON s.user_id = u.id
			LEFT JOIN groups g ON s.group_id = g.id
			WHERE s.student_code = $1`

	var s student.Student
	err := r.db.QueryRow(ctx, query, code).Scan(
		&s.ID, &s.UserID, &s.StudentCode, &s.Major, &s.Year, &s.GPA, &s.Status,
		&s.GroupID, &s.AdmittedAt, &s.CreatedAt, &s.UpdatedAt, &s.FirstName, &s.LastName, &s.Email, &s.GroupName,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *PostgresRepository) GetWithDetails(ctx context.Context, id string) (*student.StudentWithDetails, error) {
	s, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	details := &student.StudentWithDetails{
		Student: *s,
	}

	// Get enrollments
	enrollments, _ := r.GetEnrollments(ctx, id)
	details.TotalCourses = len(enrollments)

	// Get groups
	groups, _ := r.GetGroups(ctx, id)
	details.TotalGroups = len(groups)

	return details, nil
}

func (r *PostgresRepository) List(ctx context.Context, filter *student.StudentFilter) ([]*student.Student, int64, error) {
	query := `
		SELECT s.id, s.user_id, s.student_code, s.major, s.year, s.gpa, s.enrollment_status,
		       COALESCE(s.group_id::text, ''), s.admitted_at, s.created_at, s.updated_at, u.first_name, u.last_name, u.email, COALESCE(g.name, '')
		FROM students s
		JOIN users u ON s.user_id = u.id
		LEFT JOIN groups g ON s.group_id = g.id
		WHERE 1=1`

	countQuery := `SELECT COUNT(*) FROM students s WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if filter.Major != "" {
		argCount++
		query += ` AND s.major = $` + strconv.Itoa(argCount)
		countQuery += ` AND s.major = $` + strconv.Itoa(argCount)
		args = append(args, filter.Major)
	}
	if filter.Year > 0 {
		argCount++
		query += ` AND s.year = $` + strconv.Itoa(argCount)
		countQuery += ` AND s.year = $` + strconv.Itoa(argCount)
		args = append(args, filter.Year)
	}
	if filter.Status != "" {
		argCount++
		query += ` AND s.enrollment_status = $` + strconv.Itoa(argCount)
		countQuery += ` AND s.enrollment_status = $` + strconv.Itoa(argCount)
		args = append(args, filter.Status)
	}

	var total int64
	r.db.QueryRow(ctx, countQuery, args...).Scan(&total)

	query += ` ORDER BY s.created_at DESC`
	argCount++
	query += ` LIMIT $` + strconv.Itoa(argCount)
	args = append(args, filter.Limit)
	argCount++
	query += ` OFFSET $` + strconv.Itoa(argCount)
	args = append(args, filter.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var students []*student.Student
	for rows.Next() {
		var s student.Student
		err := rows.Scan(
			&s.ID, &s.UserID, &s.StudentCode, &s.Major, &s.Year, &s.GPA, &s.Status,
			&s.GroupID, &s.AdmittedAt, &s.CreatedAt, &s.UpdatedAt, &s.FirstName, &s.LastName, &s.Email, &s.GroupName,
		)
		if err != nil {
			return nil, 0, err
		}
		students = append(students, &s)
	}

	return students, total, nil
}

func (r *PostgresRepository) ListByMajor(ctx context.Context, major string, limit, offset int) ([]*student.Student, int64, error) {
	return r.List(ctx, &student.StudentFilter{
		Major:  major,
		Limit:  limit,
		Offset: offset,
	})
}

func (r *PostgresRepository) Update(ctx context.Context, s *student.Student) error {
	query := `
		UPDATE students
		SET major = $1, year = $2, gpa = $3, enrollment_status = $4, group_id = $5, updated_at = NOW()
		WHERE id = $6`

	_, err := r.db.Exec(ctx, query, s.Major, s.Year, s.GPA, s.Status, s.GroupID, s.ID)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM students WHERE id = $1`, id)
	return err
}

func (r *PostgresRepository) GetEnrollments(ctx context.Context, studentID string) ([]*student.StudentEnrollment, error) {
	query := `
		SELECT e.id, e.course_id, c.title, c.code, 
		       COALESCE(u.first_name || ' ' || u.last_name, 'N/A') as teacher_name,
		       e.status, e.enrolled_at
		FROM enrollments e
		JOIN courses c ON e.course_id = c.id
		JOIN students s ON e.student_id = s.user_id
		LEFT JOIN users u ON c.teacher_id = u.id
		WHERE s.id = $1
		ORDER BY e.enrolled_at DESC`

	rows, err := r.db.Query(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []*student.StudentEnrollment
	for rows.Next() {
		var e student.StudentEnrollment
		err := rows.Scan(&e.EnrollmentID, &e.CourseID, &e.CourseName, &e.CourseCode,
			&e.TeacherName, &e.Status, &e.EnrolledAt)
		if err != nil {
			continue
		}
		enrollments = append(enrollments, &e)
	}

	return enrollments, nil
}

func (r *PostgresRepository) GetGroups(ctx context.Context, studentID string) ([]*student.StudentGroup, error) {
	query := `
		SELECT gm.id, gm.group_id, g.name, g.course_id, c.title, gm.joined_at
		FROM group_members gm
		JOIN groups g ON gm.group_id = g.id
		JOIN courses c ON g.course_id = c.id
		JOIN students s ON gm.student_id = s.user_id
		WHERE s.id = $1
		ORDER BY gm.joined_at DESC`

	rows, err := r.db.Query(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*student.StudentGroup
	for rows.Next() {
		var g student.StudentGroup
		err := rows.Scan(&g.MembershipID, &g.GroupID, &g.GroupName, &g.CourseID,
			&g.CourseName, &g.JoinedAt)
		if err != nil {
			continue
		}
		groups = append(groups, &g)
	}

	return groups, nil
}
