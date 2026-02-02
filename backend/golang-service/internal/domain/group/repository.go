package group

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the group repository interface
type Repository interface {
	Create(ctx context.Context, group *Group) error
	GetByID(ctx context.Context, id uuid.UUID) (*GroupWithStats, error)
	GetByCode(ctx context.Context, code string) (*GroupWithStats, error)
	Update(ctx context.Context, group *Group) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]GroupWithStats, int64, error)

	// Teacher-Course-Group assignments
	AssignTeacher(ctx context.Context, assignment *TeacherCourseGroup) error
	UnassignTeacher(ctx context.Context, id uuid.UUID) error
	ListTeacherAssignments(ctx context.Context, teacherID uuid.UUID, limit, offset int) ([]TeacherCourseGroupDetails, int64, error)
	ListGroupAssignments(ctx context.Context, groupID uuid.UUID, limit, offset int) ([]TeacherCourseGroupDetails, int64, error)
	GetTeacherAssignment(ctx context.Context, teacherID, courseID, groupID uuid.UUID) (*TeacherCourseGroup, error)
}

type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new group repository
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, group *Group) error {
	query := `
		INSERT INTO groups (code, name, description, year_of_admission)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	return r.db.QueryRow(ctx, query,
		group.Code,
		group.Name,
		group.Description,
		group.YearOfAdmission,
	).Scan(&group.ID, &group.CreatedAt)
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*GroupWithStats, error) {
	query := `
		SELECT 
			g.id, g.code, g.name, g.description, g.year_of_admission, g.created_at,
			COALESCE((SELECT COUNT(*) FROM students s WHERE s.group_id = g.id), 0) AS student_count
		FROM groups g
		WHERE g.id = $1`

	var g GroupWithStats
	err := r.db.QueryRow(ctx, query, id).Scan(
		&g.ID, &g.Code, &g.Name, &g.Description, &g.YearOfAdmission, &g.CreatedAt,
		&g.StudentCount,
	)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *repository) GetByCode(ctx context.Context, code string) (*GroupWithStats, error) {
	query := `
		SELECT 
			g.id, g.code, g.name, g.description, g.year_of_admission, g.created_at,
			COALESCE((SELECT COUNT(*) FROM students s WHERE s.group_id = g.id), 0) AS student_count
		FROM groups g
		WHERE g.code = $1`

	var g GroupWithStats
	err := r.db.QueryRow(ctx, query, code).Scan(
		&g.ID, &g.Code, &g.Name, &g.Description, &g.YearOfAdmission, &g.CreatedAt,
		&g.StudentCount,
	)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *repository) Update(ctx context.Context, group *Group) error {
	query := `
		UPDATE groups 
		SET code = $1, name = $2, description = $3, year_of_admission = $4
		WHERE id = $5`

	_, err := r.db.Exec(ctx, query,
		group.Code, group.Name, group.Description, group.YearOfAdmission, group.ID,
	)
	return err
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM groups WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *repository) List(ctx context.Context, limit, offset int) ([]GroupWithStats, int64, error) {
	countQuery := `SELECT COUNT(*) FROM groups`
	var total int64
	if err := r.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			g.id, g.code, g.name, g.description, g.year_of_admission, g.created_at,
			COALESCE((SELECT COUNT(*) FROM students s WHERE s.group_id = g.id), 0) AS student_count
		FROM groups g
		ORDER BY g.code ASC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var groups []GroupWithStats
	for rows.Next() {
		var g GroupWithStats
		if err := rows.Scan(
			&g.ID, &g.Code, &g.Name, &g.Description, &g.YearOfAdmission, &g.CreatedAt,
			&g.StudentCount,
		); err != nil {
			return nil, 0, err
		}
		groups = append(groups, g)
	}
	return groups, total, nil
}

func (r *repository) AssignTeacher(ctx context.Context, assignment *TeacherCourseGroup) error {
	query := `
		INSERT INTO teacher_course_groups (teacher_id, course_id, group_id)
		VALUES ($1, $2, $3)
		RETURNING id, assigned_at`

	return r.db.QueryRow(ctx, query,
		assignment.TeacherID,
		assignment.CourseID,
		assignment.GroupID,
	).Scan(&assignment.ID, &assignment.AssignedAt)
}

func (r *repository) UnassignTeacher(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM teacher_course_groups WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *repository) ListTeacherAssignments(ctx context.Context, teacherID uuid.UUID, limit, offset int) ([]TeacherCourseGroupDetails, int64, error) {
	countQuery := `SELECT COUNT(*) FROM teacher_course_groups WHERE teacher_id = $1`
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, teacherID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			tcg.id, tcg.teacher_id, tcg.course_id, tcg.group_id, tcg.assigned_at,
			t.first_name AS teacher_first_name, t.last_name AS teacher_last_name,
			c.title AS course_title,
			g.code AS group_code
		FROM teacher_course_groups tcg
		JOIN teachers t ON tcg.teacher_id = t.user_id
		JOIN courses c ON tcg.course_id = c.id
		JOIN groups g ON tcg.group_id = g.id
		WHERE tcg.teacher_id = $1
		ORDER BY c.title, g.code
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, teacherID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var assignments []TeacherCourseGroupDetails
	for rows.Next() {
		var a TeacherCourseGroupDetails
		if err := rows.Scan(
			&a.ID, &a.TeacherID, &a.CourseID, &a.GroupID, &a.AssignedAt,
			&a.TeacherFirstName, &a.TeacherLastName,
			&a.CourseTitle, &a.GroupCode,
		); err != nil {
			return nil, 0, err
		}
		assignments = append(assignments, a)
	}
	return assignments, total, nil
}

func (r *repository) ListGroupAssignments(ctx context.Context, groupID uuid.UUID, limit, offset int) ([]TeacherCourseGroupDetails, int64, error) {
	countQuery := `SELECT COUNT(*) FROM teacher_course_groups WHERE group_id = $1`
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, groupID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			tcg.id, tcg.teacher_id, tcg.course_id, tcg.group_id, tcg.assigned_at,
			t.first_name AS teacher_first_name, t.last_name AS teacher_last_name,
			c.title AS course_title,
			g.code AS group_code
		FROM teacher_course_groups tcg
		JOIN teachers t ON tcg.teacher_id = t.user_id
		JOIN courses c ON tcg.course_id = c.id
		JOIN groups g ON tcg.group_id = g.id
		WHERE tcg.group_id = $1
		ORDER BY c.title, t.last_name
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, groupID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var assignments []TeacherCourseGroupDetails
	for rows.Next() {
		var a TeacherCourseGroupDetails
		if err := rows.Scan(
			&a.ID, &a.TeacherID, &a.CourseID, &a.GroupID, &a.AssignedAt,
			&a.TeacherFirstName, &a.TeacherLastName,
			&a.CourseTitle, &a.GroupCode,
		); err != nil {
			return nil, 0, err
		}
		assignments = append(assignments, a)
	}
	return assignments, total, nil
}

func (r *repository) GetTeacherAssignment(ctx context.Context, teacherID, courseID, groupID uuid.UUID) (*TeacherCourseGroup, error) {
	query := `
		SELECT id, teacher_id, course_id, group_id, assigned_at
		FROM teacher_course_groups
		WHERE teacher_id = $1 AND course_id = $2 AND group_id = $3`

	var a TeacherCourseGroup
	err := r.db.QueryRow(ctx, query, teacherID, courseID, groupID).Scan(
		&a.ID, &a.TeacherID, &a.CourseID, &a.GroupID, &a.AssignedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
