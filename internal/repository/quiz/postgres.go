package quiz

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/quiz"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) quiz.Repository {
	return &PostgresRepository{db: db}
}

// Quiz CRUD
func (r *PostgresRepository) Create(ctx context.Context, q *quiz.Quiz) error {
	if q.ID == "" {
		q.ID = uuid.New().String()
	}
	q.CreatedAt = time.Now()
	q.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx,
		`INSERT INTO quizzes (id, course_id, created_by, title, description, time_limit, max_attempts, start_date, end_date, published, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		q.ID, q.CourseID, q.CreatedBy, q.Title, q.Description, q.TimeLimit, q.MaxAttempts, q.StartDate, q.EndDate, q.Published, q.CreatedAt, q.UpdatedAt)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*quiz.Quiz, error) {
	q := &quiz.Quiz{}
	err := r.db.QueryRow(ctx,
		`SELECT id, course_id, created_by, title, description, time_limit, max_attempts, start_date, end_date, published, created_at, updated_at
		 FROM quizzes WHERE id = $1`, id).
		Scan(&q.ID, &q.CourseID, &q.CreatedBy, &q.Title, &q.Description, &q.TimeLimit, &q.MaxAttempts, &q.StartDate, &q.EndDate, &q.Published, &q.CreatedAt, &q.UpdatedAt)
	return q, err
}

func (r *PostgresRepository) GetByCourseID(ctx context.Context, courseID string, limit, offset int) ([]*quiz.Quiz, int64, error) {
	var total int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM quizzes WHERE course_id = $1`, courseID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT id, course_id, created_by, title, description, time_limit, max_attempts, start_date, end_date, published, created_at, updated_at
		 FROM quizzes WHERE course_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`, courseID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var quizzes []*quiz.Quiz
	for rows.Next() {
		q := &quiz.Quiz{}
		if err := rows.Scan(&q.ID, &q.CourseID, &q.CreatedBy, &q.Title, &q.Description, &q.TimeLimit, &q.MaxAttempts, &q.StartDate, &q.EndDate, &q.Published, &q.CreatedAt, &q.UpdatedAt); err != nil {
			return nil, 0, err
		}
		quizzes = append(quizzes, q)
	}
	return quizzes, total, rows.Err()
}

func (r *PostgresRepository) Update(ctx context.Context, q *quiz.Quiz) error {
	q.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx,
		`UPDATE quizzes SET title=$1, description=$2, time_limit=$3, max_attempts=$4, start_date=$5, end_date=$6, updated_at=$7 WHERE id=$8`,
		q.Title, q.Description, q.TimeLimit, q.MaxAttempts, q.StartDate, q.EndDate, q.UpdatedAt, q.ID)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM quizzes WHERE id = $1`, id)
	return err
}

func (r *PostgresRepository) Publish(ctx context.Context, id string, published bool) error {
	_, err := r.db.Exec(ctx, `UPDATE quizzes SET published = $1, updated_at = $2 WHERE id = $3`, published, time.Now(), id)
	return err
}

// Questions
func (r *PostgresRepository) CreateQuestion(ctx context.Context, q *quiz.Question) error {
	if q.ID == "" {
		q.ID = uuid.New().String()
	}

	_, err := r.db.Exec(ctx,
		`INSERT INTO quiz_questions (id, quiz_id, type, text, options, correct_answer, points, order_index)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		q.ID, q.QuizID, q.Type, q.Text, q.Options, q.CorrectAnswer, q.Points, q.OrderIndex)
	return err
}

func (r *PostgresRepository) GetQuestionsByQuizID(ctx context.Context, quizID string) ([]*quiz.Question, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, quiz_id, type, text, options, correct_answer, points, order_index
		 FROM quiz_questions WHERE quiz_id = $1 ORDER BY order_index`, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*quiz.Question
	for rows.Next() {
		q := &quiz.Question{}
		if err := rows.Scan(&q.ID, &q.QuizID, &q.Type, &q.Text, &q.Options, &q.CorrectAnswer, &q.Points, &q.OrderIndex); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return questions, rows.Err()
}

func (r *PostgresRepository) GetQuestionByID(ctx context.Context, id string) (*quiz.Question, error) {
	q := &quiz.Question{}
	err := r.db.QueryRow(ctx,
		`SELECT id, quiz_id, type, text, options, correct_answer, points, order_index
		 FROM quiz_questions WHERE id = $1`, id).
		Scan(&q.ID, &q.QuizID, &q.Type, &q.Text, &q.Options, &q.CorrectAnswer, &q.Points, &q.OrderIndex)
	return q, err
}

func (r *PostgresRepository) UpdateQuestion(ctx context.Context, q *quiz.Question) error {
	_, err := r.db.Exec(ctx,
		`UPDATE quiz_questions SET text=$1, options=$2, correct_answer=$3, points=$4, order_index=$5 WHERE id=$6`,
		q.Text, q.Options, q.CorrectAnswer, q.Points, q.OrderIndex, q.ID)
	return err
}

func (r *PostgresRepository) DeleteQuestion(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM quiz_questions WHERE id = $1`, id)
	return err
}

// Attempts
func (r *PostgresRepository) CreateAttempt(ctx context.Context, a *quiz.QuizAttempt) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	a.StartedAt = time.Now()

	// Get max score from quiz
	var maxScore int
	err := r.db.QueryRow(ctx, `SELECT COALESCE(SUM(points), 0) FROM quiz_questions WHERE quiz_id = $1`, a.QuizID).Scan(&maxScore)
	if err != nil {
		return err
	}
	a.MaxScore = maxScore

	_, err = r.db.Exec(ctx,
		`INSERT INTO quiz_attempts (id, quiz_id, student_id, started_at, max_score)
		 VALUES ($1, $2, $3, $4, $5)`,
		a.ID, a.QuizID, a.StudentID, a.StartedAt, a.MaxScore)
	return err
}

func (r *PostgresRepository) GetAttemptByID(ctx context.Context, id string) (*quiz.QuizAttempt, error) {
	a := &quiz.QuizAttempt{}
	err := r.db.QueryRow(ctx,
		`SELECT id, quiz_id, student_id, started_at, submitted_at, score, max_score, percentage
		 FROM quiz_attempts WHERE id = $1`, id).
		Scan(&a.ID, &a.QuizID, &a.StudentID, &a.StartedAt, &a.SubmittedAt, &a.Score, &a.MaxScore, &a.Percentage)
	return a, err
}

func (r *PostgresRepository) GetAttemptsByQuizAndStudent(ctx context.Context, quizID, studentID string) ([]*quiz.QuizAttempt, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, quiz_id, student_id, started_at, submitted_at, score, max_score, percentage
		 FROM quiz_attempts WHERE quiz_id = $1 AND student_id = $2 ORDER BY started_at DESC`, quizID, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []*quiz.QuizAttempt
	for rows.Next() {
		a := &quiz.QuizAttempt{}
		if err := rows.Scan(&a.ID, &a.QuizID, &a.StudentID, &a.StartedAt, &a.SubmittedAt, &a.Score, &a.MaxScore, &a.Percentage); err != nil {
			return nil, err
		}
		attempts = append(attempts, a)
	}
	return attempts, rows.Err()
}

func (r *PostgresRepository) GetAttemptsByQuiz(ctx context.Context, quizID string) ([]*quiz.QuizAttempt, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, quiz_id, student_id, started_at, submitted_at, score, max_score, percentage
		 FROM quiz_attempts WHERE quiz_id = $1 ORDER BY submitted_at DESC NULLS LAST`, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []*quiz.QuizAttempt
	for rows.Next() {
		a := &quiz.QuizAttempt{}
		if err := rows.Scan(&a.ID, &a.QuizID, &a.StudentID, &a.StartedAt, &a.SubmittedAt, &a.Score, &a.MaxScore, &a.Percentage); err != nil {
			return nil, err
		}
		attempts = append(attempts, a)
	}
	return attempts, rows.Err()
}

func (r *PostgresRepository) CountAttempts(ctx context.Context, quizID, studentID string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM quiz_attempts WHERE quiz_id = $1 AND student_id = $2 AND submitted_at IS NOT NULL`,
		quizID, studentID).Scan(&count)
	return count, err
}

func (r *PostgresRepository) SubmitAttempt(ctx context.Context, attemptID string, score float64, maxScore int) error {
	now := time.Now()
	percentage := 0.0
	if maxScore > 0 {
		percentage = (score / float64(maxScore)) * 100
	}
	_, err := r.db.Exec(ctx,
		`UPDATE quiz_attempts SET submitted_at = $1, score = $2, percentage = $3 WHERE id = $4`,
		now, score, percentage, attemptID)
	return err
}

// Answers
func (r *PostgresRepository) SaveAnswer(ctx context.Context, a *quiz.QuizAnswer) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	a.CreatedAt = time.Now()

	_, err := r.db.Exec(ctx,
		`INSERT INTO quiz_answers (id, attempt_id, question_id, answer, is_correct, points_earned, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (attempt_id, question_id) DO UPDATE SET answer = $4, is_correct = $5, points_earned = $6`,
		a.ID, a.AttemptID, a.QuestionID, a.Answer, a.IsCorrect, a.Points, a.CreatedAt)
	return err
}

func (r *PostgresRepository) GetAnswersByAttempt(ctx context.Context, attemptID string) ([]*quiz.QuizAnswer, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, attempt_id, question_id, answer, is_correct, points_earned, created_at
		 FROM quiz_answers WHERE attempt_id = $1`, attemptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []*quiz.QuizAnswer
	for rows.Next() {
		a := &quiz.QuizAnswer{}
		if err := rows.Scan(&a.ID, &a.AttemptID, &a.QuestionID, &a.Answer, &a.IsCorrect, &a.Points, &a.CreatedAt); err != nil {
			return nil, err
		}
		answers = append(answers, a)
	}
	return answers, rows.Err()
}

// Helper to convert options to JSON
func OptionsToJSON(options []string) json.RawMessage {
	data, _ := json.Marshal(options)
	return data
}
