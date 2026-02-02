// internal/domain/user/repository.go
package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	CreateProfile(ctx context.Context, tx pgx.Tx, userID uuid.UUID, role string, req *RegisterRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetProfile(ctx context.Context, userID uuid.UUID, role string) (*UserProfile, error)
	Update(ctx context.Context, user *User) error
	UpdateProfile(ctx context.Context, userID uuid.UUID, role string, req *UpdateProfileRequest) error
	List(ctx context.Context, limit, offset int) ([]User, int, error)
	Delete(ctx context.Context, id uuid.UUID) error
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.Begin(ctx)
}

func (r *repository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (email, password_hash, role, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

func (r *repository) CreateProfile(ctx context.Context, tx pgx.Tx, userID uuid.UUID, role string, req *RegisterRequest) error {
	switch role {
	case "student":
		query := `INSERT INTO students (user_id, first_name, last_name, group_name) VALUES ($1, $2, $3, $4)`
		_, err := tx.Exec(ctx, query, userID, req.FirstName, req.LastName, req.GroupName)
		return err
	case "teacher":
		query := `INSERT INTO teachers (user_id, first_name, last_name, department) VALUES ($1, $2, $3, $4)`
		_, err := tx.Exec(ctx, query, userID, req.FirstName, req.LastName, req.Department)
		return err
	case "manager":
		query := `INSERT INTO managers (user_id, first_name, last_name) VALUES ($1, $2, $3)`
		_, err := tx.Exec(ctx, query, userID, req.FirstName, req.LastName)
		return err
	default:
		return fmt.Errorf("invalid role: %s", role)
	}
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `SELECT id, email, password_hash, role, is_active, created_at, updated_at FROM users WHERE id = $1`

	var user User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, email, password_hash, role, is_active, created_at, updated_at FROM users WHERE email = $1`

	var user User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetProfile(ctx context.Context, userID uuid.UUID, role string) (*UserProfile, error) {
	user, err := r.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	profile := &UserProfile{User: *user}

	switch role {
	case "student":
		query := `SELECT user_id, first_name, last_name, group_name, created_at FROM students WHERE user_id = $1`
		var sp StudentProfile
		err := r.db.QueryRow(ctx, query, userID).Scan(&sp.UserID, &sp.FirstName, &sp.LastName, &sp.GroupName, &sp.CreatedAt)
		if err == nil {
			profile.Student = &sp
		}
	case "teacher":
		query := `SELECT user_id, first_name, last_name, department, created_at FROM teachers WHERE user_id = $1`
		var tp TeacherProfile
		err := r.db.QueryRow(ctx, query, userID).Scan(&tp.UserID, &tp.FirstName, &tp.LastName, &tp.Department, &tp.CreatedAt)
		if err == nil {
			profile.Teacher = &tp
		}
	case "manager":
		query := `SELECT user_id, first_name, last_name, created_at FROM managers WHERE user_id = $1`
		var mp ManagerProfile
		err := r.db.QueryRow(ctx, query, userID).Scan(&mp.UserID, &mp.FirstName, &mp.LastName, &mp.CreatedAt)
		if err == nil {
			profile.Manager = &mp
		}
	}

	return profile, nil
}

func (r *repository) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE users 
		SET email = $1, password_hash = $2, role = $3, is_active = $4, updated_at = now()
		WHERE id = $5`

	_, err := r.db.Exec(ctx, query,
		user.Email, user.PasswordHash, user.Role, user.IsActive, user.ID,
	)
	return err
}

func (r *repository) UpdateProfile(ctx context.Context, userID uuid.UUID, role string, req *UpdateProfileRequest) error {
	switch role {
	case "student":
		query := `UPDATE students SET first_name = COALESCE($1, first_name), last_name = COALESCE($2, last_name), group_name = COALESCE($3, group_name) WHERE user_id = $4`
		_, err := r.db.Exec(ctx, query, req.FirstName, req.LastName, req.GroupName, userID)
		return err
	case "teacher":
		query := `UPDATE teachers SET first_name = COALESCE($1, first_name), last_name = COALESCE($2, last_name), department = COALESCE($3, department) WHERE user_id = $4`
		_, err := r.db.Exec(ctx, query, req.FirstName, req.LastName, req.Department, userID)
		return err
	case "manager":
		query := `UPDATE managers SET first_name = COALESCE($1, first_name), last_name = COALESCE($2, last_name) WHERE user_id = $3`
		_, err := r.db.Exec(ctx, query, req.FirstName, req.LastName, userID)
		return err
	}
	return nil
}

func (r *repository) List(ctx context.Context, limit, offset int) ([]User, int, error) {
	countQuery := `SELECT COUNT(*) FROM users`
	var total int
	if err := r.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, email, password_hash, role, is_active, created_at, updated_at 
		FROM users 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	return users, total, nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
