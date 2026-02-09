package categorymanager

import (
	"context"
	"database/sql"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/categorymanager"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) categorymanager.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, cm *categorymanager.CategoryManager) error {
	if cm.ID == "" {
		cm.ID = uuid.New().String()
	}
	cm.CreatedAt = time.Now()
	cm.UpdatedAt = time.Now()

	query := `
		INSERT INTO category_managers (id, user_id, category_id, permission_level, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		cm.ID, cm.UserID, cm.CategoryID, cm.PermissionLevel, cm.IsActive, cm.CreatedAt, cm.UpdatedAt,
	).Scan(&cm.ID, &cm.CreatedAt, &cm.UpdatedAt)
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*categorymanager.CategoryManager, error) {
	query := `
		SELECT cm.id, cm.user_id, cm.category_id, cm.permission_level, cm.is_active, cm.created_at, cm.updated_at,
		       u.email, COALESCE(u.first_name, ''), COALESCE(u.last_name, ''), COALESCE(cc.name, '')
		FROM category_managers cm
		LEFT JOIN users u ON cm.user_id = u.id
		LEFT JOIN course_categories cc ON cm.category_id = cc.id
		WHERE cm.id = $1`

	var cm categorymanager.CategoryManager
	var firstName, lastName, categoryName string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&cm.ID, &cm.UserID, &cm.CategoryID, &cm.PermissionLevel, &cm.IsActive, &cm.CreatedAt, &cm.UpdatedAt,
		&cm.Email, &firstName, &lastName, &categoryName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	cm.FirstName = firstName
	cm.LastName = lastName
	cm.CategoryName = categoryName

	return &cm, nil
}

func (r *PostgresRepository) GetByUserAndCategory(ctx context.Context, userID, categoryID string) (*categorymanager.CategoryManager, error) {
	query := `
		SELECT cm.id, cm.user_id, cm.category_id, cm.permission_level, cm.is_active, cm.created_at, cm.updated_at,
		       u.email, COALESCE(u.first_name, ''), COALESCE(u.last_name, ''), COALESCE(cc.name, '')
		FROM category_managers cm
		LEFT JOIN users u ON cm.user_id = u.id
		LEFT JOIN course_categories cc ON cm.category_id = cc.id
		WHERE cm.user_id = $1 AND cm.category_id = $2`

	var cm categorymanager.CategoryManager
	var firstName, lastName, categoryName string

	err := r.db.QueryRow(ctx, query, userID, categoryID).Scan(
		&cm.ID, &cm.UserID, &cm.CategoryID, &cm.PermissionLevel, &cm.IsActive, &cm.CreatedAt, &cm.UpdatedAt,
		&cm.Email, &firstName, &lastName, &categoryName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	cm.FirstName = firstName
	cm.LastName = lastName
	cm.CategoryName = categoryName

	return &cm, nil
}

func (r *PostgresRepository) GetByUserID(ctx context.Context, userID string) ([]*categorymanager.CategoryManager, error) {
	query := `
		SELECT cm.id, cm.user_id, cm.category_id, cm.permission_level, cm.is_active, cm.created_at, cm.updated_at,
		       u.email, COALESCE(u.first_name, ''), COALESCE(u.last_name, ''), COALESCE(cc.name, '')
		FROM category_managers cm
		LEFT JOIN users u ON cm.user_id = u.id
		LEFT JOIN course_categories cc ON cm.category_id = cc.id
		WHERE cm.user_id = $1
		ORDER BY cm.created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var managers []*categorymanager.CategoryManager
	for rows.Next() {
		var cm categorymanager.CategoryManager
		var firstName, lastName, categoryName string

		err := rows.Scan(
			&cm.ID, &cm.UserID, &cm.CategoryID, &cm.PermissionLevel, &cm.IsActive, &cm.CreatedAt, &cm.UpdatedAt,
			&cm.Email, &firstName, &lastName, &categoryName,
		)
		if err != nil {
			return nil, err
		}

		cm.FirstName = firstName
		cm.LastName = lastName
		cm.CategoryName = categoryName

		managers = append(managers, &cm)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return managers, nil
}

func (r *PostgresRepository) GetByCategoryID(ctx context.Context, categoryID string) ([]*categorymanager.CategoryManager, error) {
	query := `
		SELECT cm.id, cm.user_id, cm.category_id, cm.permission_level, cm.is_active, cm.created_at, cm.updated_at,
		       u.email, COALESCE(u.first_name, ''), COALESCE(u.last_name, ''), COALESCE(cc.name, '')
		FROM category_managers cm
		LEFT JOIN users u ON cm.user_id = u.id
		LEFT JOIN course_categories cc ON cm.category_id = cc.id
		WHERE cm.category_id = $1
		ORDER BY cm.created_at DESC`

	rows, err := r.db.Query(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var managers []*categorymanager.CategoryManager
	for rows.Next() {
		var cm categorymanager.CategoryManager
		var firstName, lastName, categoryName string

		err := rows.Scan(
			&cm.ID, &cm.UserID, &cm.CategoryID, &cm.PermissionLevel, &cm.IsActive, &cm.CreatedAt, &cm.UpdatedAt,
			&cm.Email, &firstName, &lastName, &categoryName,
		)
		if err != nil {
			return nil, err
		}

		cm.FirstName = firstName
		cm.LastName = lastName
		cm.CategoryName = categoryName

		managers = append(managers, &cm)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return managers, nil
}

func (r *PostgresRepository) List(ctx context.Context, filter *categorymanager.CategoryManagerFilter) ([]*categorymanager.CategoryManager, int64, error) {
	query := `
		SELECT cm.id, cm.user_id, cm.category_id, cm.permission_level, cm.is_active, cm.created_at, cm.updated_at,
		       u.email, COALESCE(u.first_name, ''), COALESCE(u.last_name, ''), COALESCE(cc.name, '')
		FROM category_managers cm
		LEFT JOIN users u ON cm.user_id = u.id
		LEFT JOIN course_categories cc ON cm.category_id = cc.id
		WHERE 1=1`

	args := []interface{}{}
	argCount := 1

	if filter.CategoryID != "" {
		query += ` AND cm.category_id = $` + string(rune(argCount))
		args = append(args, filter.CategoryID)
		argCount++
	}

	if filter.PermissionLevel != "" {
		query += ` AND cm.permission_level = $` + string(rune(argCount))
		args = append(args, filter.PermissionLevel)
		argCount++
	}

	if filter.IsActive != nil {
		query += ` AND cm.is_active = $` + string(rune(argCount))
		args = append(args, *filter.IsActive)
		argCount++
	}

	query += ` ORDER BY cm.created_at DESC LIMIT $` + string(rune(argCount)) + ` OFFSET $` + string(rune(argCount+1))
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var managers []*categorymanager.CategoryManager
	for rows.Next() {
		var cm categorymanager.CategoryManager
		var firstName, lastName, categoryName string

		err := rows.Scan(
			&cm.ID, &cm.UserID, &cm.CategoryID, &cm.PermissionLevel, &cm.IsActive, &cm.CreatedAt, &cm.UpdatedAt,
			&cm.Email, &firstName, &lastName, &categoryName,
		)
		if err != nil {
			return nil, 0, err
		}

		cm.FirstName = firstName
		cm.LastName = lastName
		cm.CategoryName = categoryName

		managers = append(managers, &cm)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM category_managers cm WHERE 1=1`
	countArgs := []interface{}{}
	countArgCount := 1

	if filter.CategoryID != "" {
		countQuery += ` AND cm.category_id = $` + string(rune(countArgCount))
		countArgs = append(countArgs, filter.CategoryID)
		countArgCount++
	}

	if filter.PermissionLevel != "" {
		countQuery += ` AND cm.permission_level = $` + string(rune(countArgCount))
		countArgs = append(countArgs, filter.PermissionLevel)
		countArgCount++
	}

	if filter.IsActive != nil {
		countQuery += ` AND cm.is_active = $` + string(rune(countArgCount))
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

func (r *PostgresRepository) Update(ctx context.Context, cm *categorymanager.CategoryManager) error {
	cm.UpdatedAt = time.Now()

	query := `
		UPDATE category_managers 
		SET permission_level = $1, is_active = $2, updated_at = $3
		WHERE id = $4
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query,
		cm.PermissionLevel, cm.IsActive, cm.UpdatedAt, cm.ID,
	).Scan(&cm.UpdatedAt)
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM category_managers WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *PostgresRepository) GetWithDetails(ctx context.Context, id string) (*categorymanager.CategoryManagerWithDetails, error) {
	cm, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	details := &categorymanager.CategoryManagerWithDetails{CategoryManager: *cm}

	// Get category details
	catQuery := `SELECT id, name, description FROM course_categories WHERE id = $1`
	var catID, catName, catDesc string
	err = r.db.QueryRow(ctx, catQuery, cm.CategoryID).Scan(&catID, &catName, &catDesc)
	if err == nil {
		details.Category = &categorymanager.CategoryInfo{
			ID:          catID,
			Name:        catName,
			Description: catDesc,
		}
	}

	// Get total courses in this category
	courseQuery := `SELECT COUNT(*) FROM courses WHERE category_id = $1`
	_ = r.db.QueryRow(ctx, courseQuery, cm.CategoryID).Scan(&details.TotalCourses)

	return details, nil
}
