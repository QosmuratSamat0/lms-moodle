package teacher

import (
	"context"
	"database/sql"
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

func (r *PostgresRepository) Create(ctx context.Context, t *teacher.Teacher) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	query := `
		INSERT INTO teachers (id, user_id, employee_id, first_name, last_name, department, specialization, qualifications, bio, office_hours, phone, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		t.ID, t.UserID, t.EmployeeID, t.FirstName, t.LastName, t.Department, t.Specialization,
		t.Qualifications, t.Bio, t.OfficeHours, t.Phone, t.IsActive, t.CreatedAt, t.UpdatedAt,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*teacher.Teacher, error) {
	query := `
		SELECT t.id, t.user_id, t.employee_id, t.first_name, t.last_name, t.department, t.specialization, t.qualifications, t.bio, t.office_hours, t.phone, t.is_active, t.created_at, t.updated_at, u.email
		FROM teachers t
		LEFT JOIN users u ON t.user_id = u.id
		WHERE t.id = $1`

	var t teacher.Teacher
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.UserID, &t.EmployeeID, &t.FirstName, &t.LastName, &t.Department, &t.Specialization,
		&t.Qualifications, &t.Bio, &t.OfficeHours, &t.Phone, &t.IsActive, &t.CreatedAt, &t.UpdatedAt, &t.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &t, nil
}

func (r *PostgresRepository) GetByUserID(ctx context.Context, userID string) (*teacher.Teacher, error) {
	query := `
		SELECT t.id, t.user_id, t.employee_id, t.first_name, t.last_name, t.department, t.specialization, t.qualifications, t.bio, t.office_hours, t.phone, t.is_active, t.created_at, t.updated_at, u.email
		FROM teachers t
		LEFT JOIN users u ON t.user_id = u.id
		WHERE t.user_id = $1`

	var t teacher.Teacher
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&t.ID, &t.UserID, &t.EmployeeID, &t.FirstName, &t.LastName, &t.Department, &t.Specialization,
		&t.Qualifications, &t.Bio, &t.OfficeHours, &t.Phone, &t.IsActive, &t.CreatedAt, &t.UpdatedAt, &t.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &t, nil
}

func (r *PostgresRepository) GetByEmployeeID(ctx context.Context, employeeID string) (*teacher.Teacher, error) {
	query := `
		SELECT t.id, t.user_id, t.employee_id, t.first_name, t.last_name, t.department, t.specialization, t.qualifications, t.bio, t.office_hours, t.phone, t.is_active, t.created_at, t.updated_at, u.email
		FROM teachers t
		LEFT JOIN users u ON t.user_id = u.id
		WHERE t.employee_id = $1`

	var t teacher.Teacher
	err := r.db.QueryRow(ctx, query, employeeID).Scan(
		&t.ID, &t.UserID, &t.EmployeeID, &t.FirstName, &t.LastName, &t.Department, &t.Specialization,
		&t.Qualifications, &t.Bio, &t.OfficeHours, &t.Phone, &t.IsActive, &t.CreatedAt, &t.UpdatedAt, &t.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &t, nil
}

func (r *PostgresRepository) List(ctx context.Context, filter *teacher.TeacherFilter) ([]*teacher.Teacher, int64, error) {
	query := `
		SELECT t.id, t.user_id, t.employee_id, t.first_name, t.last_name, t.department, t.specialization, t.qualifications, t.bio, t.office_hours, t.phone, t.is_active, t.created_at, t.updated_at, u.email
		FROM teachers t
		LEFT JOIN users u ON t.user_id = u.id
		WHERE 1=1`

	args := []interface{}{}
	argCount := 1

	if filter.Department != "" {
		query += ` AND t.department = $` + string(rune(argCount))
		args = append(args, filter.Department)
		argCount++
	}

	if filter.Specialization != "" {
		query += ` AND t.specialization = $` + string(rune(argCount))
		args = append(args, filter.Specialization)
		argCount++
	}

	if filter.IsActive != nil {
		query += ` AND t.is_active = $` + string(rune(argCount))
		args = append(args, *filter.IsActive)
		argCount++
	}

	query += ` ORDER BY t.created_at DESC LIMIT $` + string(rune(argCount)) + ` OFFSET $` + string(rune(argCount+1))
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var teachers []*teacher.Teacher
	for rows.Next() {
		var t teacher.Teacher
		err := rows.Scan(
			&t.ID, &t.UserID, &t.EmployeeID, &t.FirstName, &t.LastName, &t.Department, &t.Specialization,
			&t.Qualifications, &t.Bio, &t.OfficeHours, &t.Phone, &t.IsActive, &t.CreatedAt, &t.UpdatedAt, &t.Email,
		)
		if err != nil {
			return nil, 0, err
		}
		teachers = append(teachers, &t)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM teachers t WHERE 1=1`
	countArgs := []interface{}{}
	countArgCount := 1

	if filter.Department != "" {
		countQuery += ` AND t.department = $` + string(rune(countArgCount))
		countArgs = append(countArgs, filter.Department)
		countArgCount++
	}

	if filter.Specialization != "" {
		countQuery += ` AND t.specialization = $` + string(rune(countArgCount))
		countArgs = append(countArgs, filter.Specialization)
		countArgCount++
	}

	if filter.IsActive != nil {
		countQuery += ` AND t.is_active = $` + string(rune(countArgCount))
		countArgs = append(countArgs, *filter.IsActive)
		countArgCount++
	}

	var total int64
	err = r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return teachers, total, nil
}

func (r *PostgresRepository) Update(ctx context.Context, t *teacher.Teacher) error {
	t.UpdatedAt = time.Now()

	query := `
		UPDATE teachers 
		SET department = $1, specialization = $2, qualifications = $3, bio = $4, office_hours = $5, phone = $6, is_active = $7, updated_at = $8
		WHERE id = $9
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query,
		t.Department, t.Specialization, t.Qualifications, t.Bio, t.OfficeHours, t.Phone, t.IsActive, t.UpdatedAt, t.ID,
	).Scan(&t.UpdatedAt)
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM teachers WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *PostgresRepository) GetTeacherCourses(ctx context.Context, teacherID string) ([]teacher.TeacherCourse, error) {
	query := `
		SELECT c.id, c.code, c.title, c.description, c.owner_teacher_id, c.max_points, c.is_active, c.created_at, c.updated_at
		FROM courses c
		WHERE c.owner_teacher_id = (SELECT user_id FROM teachers WHERE id = $1)`

	rows, err := r.db.Query(ctx, query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []teacher.TeacherCourse
	for rows.Next() {
		var c teacher.TeacherCourse
		err := rows.Scan(
			&c.ID, &c.Code, &c.Title, &c.Description, &c.TeacherID, &c.MaxPoints, &c.Active, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

func (r *PostgresRepository) GetTeacherGroups(ctx context.Context, teacherID string) ([]*teacher.TeacherGroup, error) {
	query := `
		SELECT DISTINCT g.id, g.course_id, g.name, g.description, g.max_students, g.created_at, g.updated_at
		FROM groups g
		JOIN courses c ON g.course_id = c.id
		WHERE c.owner_teacher_id = (SELECT user_id FROM teachers WHERE id = $1)
		ORDER BY g.created_at DESC`

	rows, err := r.db.Query(ctx, query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*teacher.TeacherGroup
	for rows.Next() {
		var g teacher.TeacherGroup
		err := rows.Scan(
			&g.ID, &g.CourseID, &g.Name, &g.Description, &g.MaxStudents, &g.CreatedAt, &g.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}
