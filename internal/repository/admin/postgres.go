package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/admin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) admin.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, a *admin.Admin) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	permissionsJSON, _ := json.Marshal(a.Permissions)

	query := `
		INSERT INTO admins (id, user_id, employee_id, department, access_level, permissions, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		a.ID, a.UserID, a.EmployeeID, a.Department, a.AccessLevel, permissionsJSON, a.IsActive, a.CreatedAt, a.UpdatedAt,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*admin.Admin, error) {
	query := `
		SELECT a.id, a.user_id, a.employee_id, a.department, a.access_level, a.permissions, a.is_active, a.created_at, a.updated_at, 
		       u.email, COALESCE(u.first_name, ''), COALESCE(u.last_name, '')
		FROM admins a
		LEFT JOIN users u ON a.user_id = u.id
		WHERE a.id = $1`

	var a admin.Admin
	var permissionsJSON []byte
	var firstName, lastName string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&a.ID, &a.UserID, &a.EmployeeID, &a.Department, &a.AccessLevel, &permissionsJSON, &a.IsActive, &a.CreatedAt, &a.UpdatedAt,
		&a.Email, &firstName, &lastName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	if len(permissionsJSON) > 0 {
		json.Unmarshal(permissionsJSON, &a.Permissions)
	}
	a.FirstName = firstName
	a.LastName = lastName

	return &a, nil
}

func (r *PostgresRepository) GetByUserID(ctx context.Context, userID string) (*admin.Admin, error) {
	query := `
		SELECT a.id, a.user_id, a.employee_id, a.department, a.access_level, a.permissions, a.is_active, a.created_at, a.updated_at,
		       u.email, COALESCE(u.first_name, ''), COALESCE(u.last_name, '')
		FROM admins a
		LEFT JOIN users u ON a.user_id = u.id
		WHERE a.user_id = $1`

	var a admin.Admin
	var permissionsJSON []byte
	var firstName, lastName string

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&a.ID, &a.UserID, &a.EmployeeID, &a.Department, &a.AccessLevel, &permissionsJSON, &a.IsActive, &a.CreatedAt, &a.UpdatedAt,
		&a.Email, &firstName, &lastName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	if len(permissionsJSON) > 0 {
		json.Unmarshal(permissionsJSON, &a.Permissions)
	}
	a.FirstName = firstName
	a.LastName = lastName

	return &a, nil
}

func (r *PostgresRepository) GetByEmployeeID(ctx context.Context, employeeID string) (*admin.Admin, error) {
	query := `
		SELECT a.id, a.user_id, a.employee_id, a.department, a.access_level, a.permissions, a.is_active, a.created_at, a.updated_at,
		       u.email, COALESCE(u.first_name, ''), COALESCE(u.last_name, '')
		FROM admins a
		LEFT JOIN users u ON a.user_id = u.id
		WHERE a.employee_id = $1`

	var a admin.Admin
	var permissionsJSON []byte
	var firstName, lastName string

	err := r.db.QueryRow(ctx, query, employeeID).Scan(
		&a.ID, &a.UserID, &a.EmployeeID, &a.Department, &a.AccessLevel, &permissionsJSON, &a.IsActive, &a.CreatedAt, &a.UpdatedAt,
		&a.Email, &firstName, &lastName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	if len(permissionsJSON) > 0 {
		json.Unmarshal(permissionsJSON, &a.Permissions)
	}
	a.FirstName = firstName
	a.LastName = lastName

	return &a, nil
}

func (r *PostgresRepository) List(ctx context.Context, filter *admin.AdminFilter) ([]*admin.Admin, int64, error) {
	query := `
		SELECT a.id, a.user_id, a.employee_id, a.department, a.access_level, a.permissions, a.is_active, a.created_at, a.updated_at,
		       u.email, COALESCE(u.first_name, ''), COALESCE(u.last_name, '')
		FROM admins a
		LEFT JOIN users u ON a.user_id = u.id
		WHERE 1=1`

	args := []interface{}{}
	argCount := 1

	if filter.Department != "" {
		query += ` AND a.department = $` + string(rune(argCount))
		args = append(args, filter.Department)
		argCount++
	}

	if filter.AccessLevel != "" {
		query += ` AND a.access_level = $` + string(rune(argCount))
		args = append(args, filter.AccessLevel)
		argCount++
	}

	if filter.IsActive != nil {
		query += ` AND a.is_active = $` + string(rune(argCount))
		args = append(args, *filter.IsActive)
		argCount++
	}

	query += ` ORDER BY a.created_at DESC LIMIT $` + string(rune(argCount)) + ` OFFSET $` + string(rune(argCount+1))
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var admins []*admin.Admin
	for rows.Next() {
		var a admin.Admin
		var permissionsJSON []byte
		var firstName, lastName string

		err := rows.Scan(
			&a.ID, &a.UserID, &a.EmployeeID, &a.Department, &a.AccessLevel, &permissionsJSON, &a.IsActive, &a.CreatedAt, &a.UpdatedAt,
			&a.Email, &firstName, &lastName,
		)
		if err != nil {
			return nil, 0, err
		}

		if len(permissionsJSON) > 0 {
			json.Unmarshal(permissionsJSON, &a.Permissions)
		}
		a.FirstName = firstName
		a.LastName = lastName

		admins = append(admins, &a)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM admins a WHERE 1=1`
	countArgs := []interface{}{}
	countArgCount := 1

	if filter.Department != "" {
		countQuery += ` AND a.department = $` + string(rune(countArgCount))
		countArgs = append(countArgs, filter.Department)
		countArgCount++
	}

	if filter.AccessLevel != "" {
		countQuery += ` AND a.access_level = $` + string(rune(countArgCount))
		countArgs = append(countArgs, filter.AccessLevel)
		countArgCount++
	}

	if filter.IsActive != nil {
		countQuery += ` AND a.is_active = $` + string(rune(countArgCount))
		countArgs = append(countArgs, *filter.IsActive)
		countArgCount++
	}

	var total int64
	err = r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return admins, total, nil
}

func (r *PostgresRepository) Update(ctx context.Context, a *admin.Admin) error {
	a.UpdatedAt = time.Now()

	permissionsJSON, _ := json.Marshal(a.Permissions)

	query := `
		UPDATE admins 
		SET department = $1, access_level = $2, permissions = $3, is_active = $4, updated_at = $5
		WHERE id = $6
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query,
		a.Department, a.AccessLevel, permissionsJSON, a.IsActive, a.UpdatedAt, a.ID,
	).Scan(&a.UpdatedAt)
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM admins WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
