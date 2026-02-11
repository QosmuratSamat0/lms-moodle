package manager

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/manager"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) manager.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, m *manager.Manager) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	manageCategoriesJSON, _ := json.Marshal(m.ManagesCategories)
	manageTeachersJSON, _ := json.Marshal(m.ManagesTeachers)

	query := `
		INSERT INTO managers (id, user_id, employee_id, department, manages_categories, manages_teachers, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		m.ID, m.UserID, m.EmployeeID, m.Department, manageCategoriesJSON, manageTeachersJSON, m.IsActive, m.CreatedAt, m.UpdatedAt,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*manager.Manager, error) {
	query := `
		SELECT m.id, m.user_id, COALESCE(m.employee_id,''), COALESCE(m.department,''), m.manages_categories, m.manages_teachers, m.is_active, m.created_at, m.updated_at,
		       u.email, COALESCE(m.first_name, ''), COALESCE(m.last_name, '')
		FROM managers m
		LEFT JOIN users u ON m.user_id = u.id
		WHERE m.id = $1`

	var m manager.Manager
	var manageCategoriesJSON, manageTeachersJSON []byte
	var firstName, lastName string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.UserID, &m.EmployeeID, &m.Department, &manageCategoriesJSON, &manageTeachersJSON, &m.IsActive, &m.CreatedAt, &m.UpdatedAt,
		&m.Email, &firstName, &lastName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	if len(manageCategoriesJSON) > 0 {
		json.Unmarshal(manageCategoriesJSON, &m.ManagesCategories)
	}
	if len(manageTeachersJSON) > 0 {
		json.Unmarshal(manageTeachersJSON, &m.ManagesTeachers)
	}
	m.FirstName = firstName
	m.LastName = lastName

	return &m, nil
}

func (r *PostgresRepository) GetByUserID(ctx context.Context, userID string) (*manager.Manager, error) {
	query := `
		SELECT m.id, m.user_id, COALESCE(m.employee_id,''), COALESCE(m.department,''), m.manages_categories, m.manages_teachers, m.is_active, m.created_at, m.updated_at,
		       u.email, COALESCE(m.first_name, ''), COALESCE(m.last_name, '')
		FROM managers m
		LEFT JOIN users u ON m.user_id = u.id
		WHERE m.user_id = $1`

	var m manager.Manager
	var manageCategoriesJSON, manageTeachersJSON []byte
	var firstName, lastName string

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&m.ID, &m.UserID, &m.EmployeeID, &m.Department, &manageCategoriesJSON, &manageTeachersJSON, &m.IsActive, &m.CreatedAt, &m.UpdatedAt,
		&m.Email, &firstName, &lastName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	if len(manageCategoriesJSON) > 0 {
		json.Unmarshal(manageCategoriesJSON, &m.ManagesCategories)
	}
	if len(manageTeachersJSON) > 0 {
		json.Unmarshal(manageTeachersJSON, &m.ManagesTeachers)
	}
	m.FirstName = firstName
	m.LastName = lastName

	return &m, nil
}

func (r *PostgresRepository) GetByEmployeeID(ctx context.Context, employeeID string) (*manager.Manager, error) {
	query := `
		SELECT m.id, m.user_id, m.employee_id, m.department, m.manages_categories, m.manages_teachers, m.is_active, m.created_at, m.updated_at,
		       u.email, COALESCE(m.first_name, ''), COALESCE(m.last_name, '')
		FROM managers m
		LEFT JOIN users u ON m.user_id = u.id
		WHERE m.employee_id = $1`

	var m manager.Manager
	var manageCategoriesJSON, manageTeachersJSON []byte
	var firstName, lastName string

	err := r.db.QueryRow(ctx, query, employeeID).Scan(
		&m.ID, &m.UserID, &m.EmployeeID, &m.Department, &manageCategoriesJSON, &manageTeachersJSON, &m.IsActive, &m.CreatedAt, &m.UpdatedAt,
		&m.Email, &firstName, &lastName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	if len(manageCategoriesJSON) > 0 {
		json.Unmarshal(manageCategoriesJSON, &m.ManagesCategories)
	}
	if len(manageTeachersJSON) > 0 {
		json.Unmarshal(manageTeachersJSON, &m.ManagesTeachers)
	}
	m.FirstName = firstName
	m.LastName = lastName

	return &m, nil
}

func (r *PostgresRepository) GetByDepartment(ctx context.Context, department string) ([]*manager.Manager, error) {
	query := `
		SELECT m.id, m.user_id, m.employee_id, m.department, m.manages_categories, m.manages_teachers, m.is_active, m.created_at, m.updated_at,
		       u.email, COALESCE(m.first_name, ''), COALESCE(m.last_name, '')
		FROM managers m
		LEFT JOIN users u ON m.user_id = u.id
		WHERE m.department = $1
		ORDER BY m.created_at DESC`

	rows, err := r.db.Query(ctx, query, department)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var managers []*manager.Manager
	for rows.Next() {
		var m manager.Manager
		var manageCategoriesJSON, manageTeachersJSON []byte
		var firstName, lastName string

		err := rows.Scan(
			&m.ID, &m.UserID, &m.EmployeeID, &m.Department, &manageCategoriesJSON, &manageTeachersJSON, &m.IsActive, &m.CreatedAt, &m.UpdatedAt,
			&m.Email, &firstName, &lastName,
		)
		if err != nil {
			return nil, err
		}

		if len(manageCategoriesJSON) > 0 {
			json.Unmarshal(manageCategoriesJSON, &m.ManagesCategories)
		}
		if len(manageTeachersJSON) > 0 {
			json.Unmarshal(manageTeachersJSON, &m.ManagesTeachers)
		}
		m.FirstName = firstName
		m.LastName = lastName

		managers = append(managers, &m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return managers, nil
}

func (r *PostgresRepository) List(ctx context.Context, filter *manager.ManagerFilter) ([]*manager.Manager, int64, error) {
	query := `
		SELECT m.id, m.user_id, COALESCE(m.employee_id,''), COALESCE(m.department,''), m.manages_categories, m.manages_teachers, m.is_active, m.created_at, m.updated_at,
		       u.email, COALESCE(m.first_name, ''), COALESCE(m.last_name, '')
		FROM managers m
		LEFT JOIN users u ON m.user_id = u.id
		WHERE 1=1`

	args := []interface{}{}
	argCount := 1

	if filter.Department != "" {
		query += ` AND m.department = $` + strconv.Itoa(argCount)
		args = append(args, filter.Department)
		argCount++
	}

	if filter.IsActive != nil {
		query += ` AND m.is_active = $` + strconv.Itoa(argCount)
		args = append(args, *filter.IsActive)
		argCount++
	}

	query += ` ORDER BY m.created_at DESC LIMIT $` + strconv.Itoa(argCount) + ` OFFSET $` + strconv.Itoa(argCount+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var managers []*manager.Manager
	for rows.Next() {
		var m manager.Manager
		var manageCategoriesJSON, manageTeachersJSON []byte
		var firstName, lastName string

		err := rows.Scan(
			&m.ID, &m.UserID, &m.EmployeeID, &m.Department, &manageCategoriesJSON, &manageTeachersJSON, &m.IsActive, &m.CreatedAt, &m.UpdatedAt,
			&m.Email, &firstName, &lastName,
		)
		if err != nil {
			return nil, 0, err
		}

		if len(manageCategoriesJSON) > 0 {
			json.Unmarshal(manageCategoriesJSON, &m.ManagesCategories)
		}
		if len(manageTeachersJSON) > 0 {
			json.Unmarshal(manageTeachersJSON, &m.ManagesTeachers)
		}
		m.FirstName = firstName
		m.LastName = lastName

		managers = append(managers, &m)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM managers m WHERE 1=1`
	countArgs := []interface{}{}
	countArgCount := 1

	if filter.Department != "" {
		countQuery += ` AND m.department = $` + strconv.Itoa(countArgCount)
		countArgs = append(countArgs, filter.Department)
		countArgCount++
	}

	if filter.IsActive != nil {
		countQuery += ` AND m.is_active = $` + strconv.Itoa(countArgCount)
		countArgs = append(countArgs, *filter.IsActive)
		countArgCount++
	}

	var total int64
	err = r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return managers, total, nil
}

func (r *PostgresRepository) Update(ctx context.Context, m *manager.Manager) error {
	m.UpdatedAt = time.Now()

	manageCategoriesJSON, _ := json.Marshal(m.ManagesCategories)
	manageTeachersJSON, _ := json.Marshal(m.ManagesTeachers)

	query := `
		UPDATE managers 
		SET department = $1, manages_categories = $2, manages_teachers = $3, is_active = $4, updated_at = $5
		WHERE id = $6
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query,
		m.Department, manageCategoriesJSON, manageTeachersJSON, m.IsActive, m.UpdatedAt, m.ID,
	).Scan(&m.UpdatedAt)
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM managers WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *PostgresRepository) AddManagedCategory(ctx context.Context, managerID, categoryID string) error {
	query := `
		UPDATE managers
		SET manages_categories = CASE
			WHEN manages_categories @> $2::jsonb THEN manages_categories
			ELSE manages_categories || $2::jsonb
		END
		WHERE id = $1`

	categoryJSON := json.RawMessage(`["` + categoryID + `"]`)
	_, err := r.db.Exec(ctx, query, managerID, categoryJSON)
	return err
}

func (r *PostgresRepository) RemoveManagedCategory(ctx context.Context, managerID, categoryID string) error {
	query := `
		UPDATE managers
		SET manages_categories = manages_categories - $2
		WHERE id = $1`

	_, err := r.db.Exec(ctx, query, managerID, categoryID)
	return err
}

func (r *PostgresRepository) AddManagedTeacher(ctx context.Context, managerID, teacherID string) error {
	query := `
		UPDATE managers
		SET manages_teachers = CASE
			WHEN manages_teachers @> $2::jsonb THEN manages_teachers
			ELSE manages_teachers || $2::jsonb
		END
		WHERE id = $1`

	teacherJSON := json.RawMessage(`["` + teacherID + `"]`)
	_, err := r.db.Exec(ctx, query, managerID, teacherJSON)
	return err
}

func (r *PostgresRepository) RemoveManagedTeacher(ctx context.Context, managerID, teacherID string) error {
	query := `
		UPDATE managers
		SET manages_teachers = manages_teachers - $2
		WHERE id = $1`

	_, err := r.db.Exec(ctx, query, managerID, teacherID)
	return err
}

func (r *PostgresRepository) GetWithDetails(ctx context.Context, id string) (*manager.ManagerWithDetails, error) {
	m, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	details := &manager.ManagerWithDetails{Manager: *m}

	// Get categories
	if len(m.ManagesCategories) > 0 {
		catQuery := `SELECT id, name FROM course_categories WHERE id = ANY($1)`
		rows, err := r.db.Query(ctx, catQuery, pq.Array(m.ManagesCategories))
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var cat manager.CategorySummary
			if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
				return nil, err
			}
			details.Categories = append(details.Categories, cat)
		}
	}

	// Get teachers
	if len(m.ManagesTeachers) > 0 {
		teachQuery := `SELECT id, employee_id, CONCAT(first_name, ' ', last_name) as full_name FROM teachers WHERE id = ANY($1)`
		rows, err := r.db.Query(ctx, teachQuery, pq.Array(m.ManagesTeachers))
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var teacher manager.TeacherSummary
			if err := rows.Scan(&teacher.ID, &teacher.EmployeeID, &teacher.FullName); err != nil {
				return nil, err
			}
			details.Teachers = append(details.Teachers, teacher)
		}
	}

	return details, nil
}
