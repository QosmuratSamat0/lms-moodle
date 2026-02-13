package quiz

import "context"

type Repository interface {
	// Quiz CRUD
	Create(ctx context.Context, quiz *Quiz) error
	GetByID(ctx context.Context, id string) (*Quiz, error)
	GetByCourseID(ctx context.Context, courseID string, limit, offset int) ([]*Quiz, int64, error)
	Update(ctx context.Context, quiz *Quiz) error
	Delete(ctx context.Context, id string) error
	Publish(ctx context.Context, id string, published bool) error

	// Questions
	CreateQuestion(ctx context.Context, question *Question) error
	GetQuestionsByQuizID(ctx context.Context, quizID string) ([]*Question, error)
	GetQuestionByID(ctx context.Context, id string) (*Question, error)
	UpdateQuestion(ctx context.Context, question *Question) error
	DeleteQuestion(ctx context.Context, id string) error

	// Attempts
	CreateAttempt(ctx context.Context, attempt *QuizAttempt) error
	GetAttemptByID(ctx context.Context, id string) (*QuizAttempt, error)
	GetAttemptsByQuizAndStudent(ctx context.Context, quizID, studentID string) ([]*QuizAttempt, error)
	GetAttemptsByQuiz(ctx context.Context, quizID string) ([]*QuizAttempt, error)
	CountAttempts(ctx context.Context, quizID, studentID string) (int, error)
	SubmitAttempt(ctx context.Context, attemptID string, score float64, maxScore int) error

	// Answers
	SaveAnswer(ctx context.Context, answer *QuizAnswer) error
	GetAnswersByAttempt(ctx context.Context, attemptID string) ([]*QuizAnswer, error)
}
