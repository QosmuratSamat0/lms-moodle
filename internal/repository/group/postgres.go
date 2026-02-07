package group

import (
	"context"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/group"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) group.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, g *group.Group) error {
	if g.ID == "" {
		g.ID = uuid.New().String()
	}
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx,
		`INSERT INTO groups (id, course_id, name, description, max_students, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		g.ID, g.CourseID, g.Name, g.Description, g.MaxStudents, g.CreatedAt, g.UpdatedAt)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*group.Group, error) {
	g := &group.Group{}
	err := r.db.QueryRow(ctx,
		`SELECT id, course_id, name, description, max_students, created_at, updated_at
		 FROM groups WHERE id = $1`, id).
		Scan(&g.ID, &g.CourseID, &g.Name, &g.Description, &g.MaxStudents, &g.CreatedAt, &g.UpdatedAt)
	return g, err
}

func (r *PostgresRepository) GetByCourseID(ctx context.Context, courseID string) ([]*group.Group, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, course_id, name, description, max_students, created_at, updated_at
		 FROM groups WHERE course_id = $1 ORDER BY name`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*group.Group
	for rows.Next() {
		g := &group.Group{}
		if err := rows.Scan(&g.ID, &g.CourseID, &g.Name, &g.Description, &g.MaxStudents, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, rows.Err()
}

func (r *PostgresRepository) Update(ctx context.Context, g *group.Group) error {
	g.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx,
		`UPDATE groups SET name=$1, description=$2, max_students=$3, updated_at=$4 WHERE id=$5`,
		g.Name, g.Description, g.MaxStudents, g.UpdatedAt, g.ID)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM groups WHERE id = $1`, id)
	return err
}

func (r *PostgresRepository) AddMember(ctx context.Context, m *group.GroupMember) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	m.JoinedAt = time.Now()

	_, err := r.db.Exec(ctx,
		`INSERT INTO group_members (id, group_id, student_id, joined_at)
		 VALUES ($1, $2, $3, $4)`,
		m.ID, m.GroupID, m.StudentID, m.JoinedAt)
	return err
}

func (r *PostgresRepository) RemoveMember(ctx context.Context, groupID, studentID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM group_members WHERE group_id = $1 AND student_id = $2`, groupID, studentID)
	return err
}

func (r *PostgresRepository) GetMembers(ctx context.Context, groupID string) ([]*group.GroupMember, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, group_id, student_id, joined_at FROM group_members WHERE group_id = $1`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*group.GroupMember
	for rows.Next() {
		m := &group.GroupMember{}
		if err := rows.Scan(&m.ID, &m.GroupID, &m.StudentID, &m.JoinedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *PostgresRepository) GetStudentGroups(ctx context.Context, studentID string) ([]*group.Group, error) {
	rows, err := r.db.Query(ctx,
		`SELECT g.id, g.course_id, g.name, g.description, g.max_students, g.created_at, g.updated_at
		 FROM groups g
		 INNER JOIN group_members gm ON g.id = gm.group_id
		 WHERE gm.student_id = $1`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*group.Group
	for rows.Next() {
		g := &group.Group{}
		if err := rows.Scan(&g.ID, &g.CourseID, &g.Name, &g.Description, &g.MaxStudents, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, rows.Err()
}

func (r *PostgresRepository) IsMember(ctx context.Context, groupID, studentID string) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM group_members WHERE group_id = $1 AND student_id = $2`, groupID, studentID).
		Scan(&count)
	return count > 0, err
}

func (r *PostgresRepository) CountMembers(ctx context.Context, groupID string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM group_members WHERE group_id = $1`, groupID).Scan(&count)
	return count, err
}
