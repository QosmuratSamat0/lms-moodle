package manager

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the manager repository interface
type Repository interface {
	GetByID(ctx context.Context, userID uuid.UUID) (*ManagerWithUser, error)
	List(ctx context.Context, filter *ManagerFilter, page, limit int) ([]ManagerWithUser, int, error)
	Update(ctx context.Context, userID uuid.UUID, req *UpdateManagerRequest) error
	GetSystemOverview(ctx context.Context) (*SystemOverview, error)
	ActivateUser(ctx context.Context, userID uuid.UUID) error
	DeactivateUser(ctx context.Context, userID uuid.UUID) error
	BulkActivateUsers(ctx context.Context, userIDs []uuid.UUID) error
	BulkDeactivateUsers(ctx context.Context, userIDs []uuid.UUID) error
}

type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new manager repository
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) GetByID(ctx context.Context, userID uuid.UUID) (*ManagerWithUser, error) {
	query := `
		SELECT m.user_id, m.first_name, m.last_name, m.created_at,
		       u.email, u.is_active
		FROM managers m
		JOIN users u ON u.id = m.user_id
		WHERE m.user_id = $1`

	var manager ManagerWithUser
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&manager.UserID, &manager.FirstName, &manager.LastName, &manager.CreatedAt,
		&manager.Email, &manager.IsActive,
	)
	if err != nil {
		return nil, err
	}

	return &manager, nil
}

func (r *repository) List(ctx context.Context, filter *ManagerFilter, page, limit int) ([]ManagerWithUser, int, error) {
	var conditions []string
	var args []any
	argIdx := 1

	if filter != nil {
		if filter.IsActive != nil {
			conditions = append(conditions, fmt.Sprintf("u.is_active = $%d", argIdx))
			args = append(args, *filter.IsActive)
			argIdx++
		}
		if filter.Search != nil && *filter.Search != "" {
			searchPattern := "%" + *filter.Search + "%"
			conditions = append(conditions, fmt.Sprintf(
				"(m.first_name ILIKE $%d OR m.last_name ILIKE $%d OR u.email ILIKE $%d)",
				argIdx, argIdx, argIdx,
			))
			args = append(args, searchPattern)
			argIdx++
		}
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM managers m
		JOIN users u ON u.id = m.user_id
		%s`, whereClause)

	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Main query
	query := fmt.Sprintf(`
		SELECT m.user_id, m.first_name, m.last_name, m.created_at,
		       u.email, u.is_active
		FROM managers m
		JOIN users u ON u.id = m.user_id
		%s
		ORDER BY m.created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)

	args = append(args, limit, (page-1)*limit)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var managers []ManagerWithUser
	for rows.Next() {
		var m ManagerWithUser
		if err := rows.Scan(
			&m.UserID, &m.FirstName, &m.LastName, &m.CreatedAt,
			&m.Email, &m.IsActive,
		); err != nil {
			return nil, 0, err
		}
		managers = append(managers, m)
	}

	return managers, total, rows.Err()
}

func (r *repository) Update(ctx context.Context, userID uuid.UUID, req *UpdateManagerRequest) error {
	query := `
		UPDATE managers
		SET first_name = COALESCE($2, first_name),
		    last_name = COALESCE($3, last_name)
		WHERE user_id = $1`

	_, err := r.db.Exec(ctx, query, userID, req.FirstName, req.LastName)
	return err
}

func (r *repository) GetSystemOverview(ctx context.Context) (*SystemOverview, error) {
	query := `
		SELECT 
			(SELECT COUNT(*) FROM users) as total_users,
			(SELECT COUNT(*) FROM users WHERE is_active = true) as active_users,
			(SELECT COUNT(*) FROM students) as total_students,
			(SELECT COUNT(*) FROM teachers) as total_teachers,
			(SELECT COUNT(*) FROM managers) as total_managers,
			(SELECT COUNT(*) FROM courses) as total_courses,
			(SELECT COUNT(*) FROM courses WHERE is_active = true) as active_courses,
			(SELECT COUNT(*) FROM enrollments WHERE status = 'active') as total_enrollment,
			(SELECT COUNT(*) FROM enrollments WHERE status = 'pending') as pending_enrollments,
			COALESCE((SELECT AVG(score)::float FROM grades), 0) as average_grade`

	var overview SystemOverview
	err := r.db.QueryRow(ctx, query).Scan(
		&overview.TotalUsers, &overview.ActiveUsers,
		&overview.TotalStudents, &overview.TotalTeachers, &overview.TotalManagers,
		&overview.TotalCourses, &overview.ActiveCourses,
		&overview.TotalEnrollment, &overview.PendingEnroll, &overview.AvgGrade,
	)
	if err != nil {
		return nil, err
	}

	return &overview, nil
}

func (r *repository) ActivateUser(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET is_active = true, updated_at = now() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *repository) DeactivateUser(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET is_active = false, updated_at = now() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *repository) BulkActivateUsers(ctx context.Context, userIDs []uuid.UUID) error {
	query := `UPDATE users SET is_active = true, updated_at = now() WHERE id = ANY($1)`
	_, err := r.db.Exec(ctx, query, userIDs)
	return err
}

func (r *repository) BulkDeactivateUsers(ctx context.Context, userIDs []uuid.UUID) error {
	query := `UPDATE users SET is_active = false, updated_at = now() WHERE id = ANY($1)`
	_, err := r.db.Exec(ctx, query, userIDs)
	return err
}
