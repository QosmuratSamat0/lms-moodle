package coursecategory

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/coursecategory"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) coursecategory.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, cc *coursecategory.CourseCategory) error {
	if cc.ID == "" {
		cc.ID = uuid.New().String()
	}
	cc.CreatedAt = time.Now()
	cc.UpdatedAt = time.Now()

	query := `
		INSERT INTO course_categories (id, name, description, icon, "order", is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		cc.ID, cc.Name, cc.Description, cc.Icon, cc.Order, cc.IsActive, cc.CreatedAt, cc.UpdatedAt,
	).Scan(&cc.ID, &cc.CreatedAt, &cc.UpdatedAt)
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*coursecategory.CourseCategory, error) {
	query := `
		SELECT cc.id, cc.name, cc.description, cc.icon, cc."order", cc.is_active, cc.created_at, cc.updated_at,
		       COUNT(c.id) as total_courses
		FROM course_categories cc
		LEFT JOIN courses c ON cc.id = c.category_id
		WHERE cc.id = $1
		GROUP BY cc.id`

	var cc coursecategory.CourseCategory

	err := r.db.QueryRow(ctx, query, id).Scan(
		&cc.ID, &cc.Name, &cc.Description, &cc.Icon, &cc.Order, &cc.IsActive, &cc.CreatedAt, &cc.UpdatedAt,
		&cc.TotalCourses,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &cc, nil
}

func (r *PostgresRepository) List(ctx context.Context, filter *coursecategory.CourseCategoryFilter) ([]*coursecategory.CourseCategory, int64, error) {
	query := `
		SELECT cc.id, cc.name, cc.description, cc.icon, cc."order", cc.is_active, cc.created_at, cc.updated_at,
		       COUNT(c.id) as total_courses
		FROM course_categories cc
		LEFT JOIN courses c ON cc.id = c.category_id
		WHERE 1=1`

	args := []interface{}{}
	argCount := 1

	if filter.IsActive != nil {
		query += ` AND cc.is_active = $` + strconv.Itoa(argCount)
		args = append(args, *filter.IsActive)
		argCount++
	}

	query += ` GROUP BY cc.id ORDER BY cc."order" ASC, cc.created_at DESC LIMIT $` + strconv.Itoa(argCount) + ` OFFSET $` + strconv.Itoa(argCount+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var categories []*coursecategory.CourseCategory
	for rows.Next() {
		var cc coursecategory.CourseCategory

		err := rows.Scan(
			&cc.ID, &cc.Name, &cc.Description, &cc.Icon, &cc.Order, &cc.IsActive, &cc.CreatedAt, &cc.UpdatedAt,
			&cc.TotalCourses,
		)
		if err != nil {
			return nil, 0, err
		}

		categories = append(categories, &cc)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM course_categories WHERE 1=1`
	countArgs := []interface{}{}
	countArgCount := 1

	if filter.IsActive != nil {
		countQuery += ` AND is_active = $` + strconv.Itoa(countArgCount)
		countArgs = append(countArgs, *filter.IsActive)
		countArgCount++
	}

	var total int64
	err = r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *PostgresRepository) Update(ctx context.Context, cc *coursecategory.CourseCategory) error {
	cc.UpdatedAt = time.Now()

	query := `
		UPDATE course_categories 
		SET name = $1, description = $2, icon = $3, "order" = $4, is_active = $5, updated_at = $6
		WHERE id = $7
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query,
		cc.Name, cc.Description, cc.Icon, cc.Order, cc.IsActive, cc.UpdatedAt, cc.ID,
	).Scan(&cc.UpdatedAt)
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM course_categories WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
