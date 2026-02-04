package user

import (
	"context"

	"github.com/ap1-final-mini-moodle/internal/domain/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) user.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(u *user.User) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO users (id, email, password_hash, role, is_active, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		u.ID, u.Email, u.Password, u.Role, u.Active, u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *PostgresRepository) GetByID(id string) (*user.User, error) {
	u := &user.User{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, email, password_hash, role, is_active, created_at, updated_at
		 FROM users WHERE id = $1`, id).
		Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.Active, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func (r *PostgresRepository) GetByEmail(email string) (*user.User, error) {
	u := &user.User{}
	err := r.db.QueryRow(context.Background(),
		`SELECT id, email, password_hash, role, is_active, created_at, updated_at
		 FROM users WHERE email = $1`, email).
		Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.Active, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func (r *PostgresRepository) Update(u *user.User) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE users SET is_active=$1, updated_at=$2 WHERE id=$3`,
		u.Active, u.UpdatedAt, u.ID)
	return err
}

func (r *PostgresRepository) List(skip, take int) ([]*user.User, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, email, password_hash, role, is_active, created_at, updated_at
		 FROM users OFFSET $1 LIMIT $2`, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		u := &user.User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.Active, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *PostgresRepository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	return err
}
