package quiz

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/quiz"
	appErrors "github.com/ap1-final-mini-moodle/internal/shared/errors"
)

type Service struct {
	repo quiz.Repository
}

func NewService(repo quiz.Repository) *Service {
	return &Service{repo: repo}
}

// Quiz CRUD

func (s *Service) CreateQuiz(ctx context.Context, createdBy string, input *quiz.CreateQuizInput) (*quiz.Quiz, error) {
	if input.CourseID == "" {
		return nil, appErrors.ErrMissingRequired
	}
	if input.Title == "" {
		return nil, appErrors.ErrMissingRequired
	}
	if len(input.Questions) == 0 {
		return nil, appErrors.ErrMissingRequired
	}

	// Validate date range
	if input.StartDate != nil && input.EndDate != nil {
		if input.EndDate.Before(*input.StartDate) {
			return nil, appErrors.ErrInvalidDateRange
		}
	}

	// Validate questions
	for _, q := range input.Questions {
		if err := validateQuestion(q); err != nil {
			return nil, err
		}
	}

	q := &quiz.Quiz{
		CourseID:    input.CourseID,
		CreatedBy:   createdBy,
		Title:       input.Title,
		Description: input.Description,
		TimeLimit:   input.TimeLimit,
		MaxAttempts: input.MaxAttempts,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		Published:   false,
	}

	if err := s.repo.Create(ctx, q); err != nil {
		return nil, err
	}

	// Create questions
	for i, qInput := range input.Questions {
		question := &quiz.Question{
			QuizID:        q.ID,
			Type:          qInput.Type,
			Text:          qInput.Text,
			Options:       optionsToJSON(qInput.Options),
			CorrectAnswer: qInput.CorrectAnswer,
			Points:        qInput.Points,
			OrderIndex:    i + 1,
		}
		if err := s.repo.CreateQuestion(ctx, question); err != nil {
			return nil, err
		}
	}

	return q, nil
}

func (s *Service) GetQuizByID(ctx context.Context, id string) (*quiz.Quiz, error) {
	if id == "" {
		return nil, appErrors.ErrInvalidID
	}

	q, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrQuizNotFound
	}

	return q, nil
}

func (s *Service) GetQuizzesByCourse(ctx context.Context, courseID string, limit, offset int) ([]*quiz.Quiz, int64, error) {
	if courseID == "" {
		return nil, 0, appErrors.ErrInvalidID
	}

	if limit <= 0 {
		limit = 20
	}

	return s.repo.GetByCourseID(ctx, courseID, limit, offset)
}

func (s *Service) GetQuestionsForQuiz(ctx context.Context, quizID string, includeAnswers bool) ([]*quiz.Question, error) {
	if quizID == "" {
		return nil, appErrors.ErrInvalidID
	}

	questions, err := s.repo.GetQuestionsByQuizID(ctx, quizID)
	if err != nil {
		return nil, err
	}

	// Hide correct answers if not requested
	if !includeAnswers {
		for _, q := range questions {
			q.CorrectAnswer = ""
		}
	}

	return questions, nil
}

func (s *Service) PublishQuiz(ctx context.Context, id string, published bool) error {
	if id == "" {
		return appErrors.ErrInvalidID
	}

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return appErrors.ErrQuizNotFound
	}

	return s.repo.Publish(ctx, id, published)
}

func (s *Service) DeleteQuiz(ctx context.Context, id string) error {
	if id == "" {
		return appErrors.ErrInvalidID
	}

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return appErrors.ErrQuizNotFound
	}

	return s.repo.Delete(ctx, id)
}

// Quiz Attempts

func (s *Service) StartAttempt(ctx context.Context, quizID, studentID string) (*quiz.QuizAttempt, error) {
	if quizID == "" || studentID == "" {
		return nil, appErrors.ErrMissingRequired
	}

	// Get quiz
	q, err := s.repo.GetByID(ctx, quizID)
	if err != nil {
		return nil, appErrors.ErrQuizNotFound
	}

	// Check if quiz is published
	if !q.Published {
		return nil, appErrors.ErrQuizNotPublished
	}

	// Check time restrictions
	now := time.Now()
	if q.StartDate != nil && now.Before(*q.StartDate) {
		return nil, appErrors.ErrQuizNotStarted
	}
	if q.EndDate != nil && now.After(*q.EndDate) {
		return nil, appErrors.ErrQuizEnded
	}

	// Check max attempts
	if q.MaxAttempts > 0 {
		attemptCount, _ := s.repo.CountAttempts(ctx, quizID, studentID)
		if attemptCount >= q.MaxAttempts {
			return nil, appErrors.ErrMaxAttemptsReached
		}
	}

	attempt := &quiz.QuizAttempt{
		QuizID:    quizID,
		StudentID: studentID,
	}

	if err := s.repo.CreateAttempt(ctx, attempt); err != nil {
		return nil, err
	}

	return attempt, nil
}

func (s *Service) SubmitAttempt(ctx context.Context, attemptID string, input *quiz.SubmitQuizInput) (*quiz.QuizResult, error) {
	if attemptID == "" {
		return nil, appErrors.ErrInvalidID
	}
	if len(input.Answers) == 0 {
		return nil, appErrors.ErrMissingRequired
	}

	// Get attempt
	attempt, err := s.repo.GetAttemptByID(ctx, attemptID)
	if err != nil {
		return nil, appErrors.ErrAttemptNotFound
	}

	// Check if already submitted
	if attempt.SubmittedAt != nil {
		return nil, appErrors.ErrAttemptAlreadySubmitted
	}

	// Get quiz
	q, err := s.repo.GetByID(ctx, attempt.QuizID)
	if err != nil {
		return nil, appErrors.ErrQuizNotFound
	}

	// Check time limit
	if q.TimeLimit > 0 {
		elapsed := time.Since(attempt.StartedAt)
		if elapsed.Minutes() > float64(q.TimeLimit)+1 { // +1 minute grace period
			return nil, appErrors.ErrQuizEnded
		}
	}

	// Get questions
	questions, err := s.repo.GetQuestionsByQuizID(ctx, attempt.QuizID)
	if err != nil {
		return nil, err
	}

	// Create question map for quick lookup
	questionMap := make(map[string]*quiz.Question)
	for _, question := range questions {
		questionMap[question.ID] = question
	}

	// Grade answers
	var totalScore float64
	var maxScore int
	var answerResults []quiz.AnswerResult

	for _, answer := range input.Answers {
		question, exists := questionMap[answer.QuestionID]
		if !exists {
			continue
		}

		isCorrect := gradeAnswer(question, answer.Answer)
		pointsEarned := 0
		if isCorrect {
			pointsEarned = question.Points
			totalScore += float64(pointsEarned)
		}
		maxScore += question.Points

		// Save answer
		quizAnswer := &quiz.QuizAnswer{
			AttemptID:  attemptID,
			QuestionID: answer.QuestionID,
			Answer:     answer.Answer,
			IsCorrect:  isCorrect,
			Points:     pointsEarned,
		}
		s.repo.SaveAnswer(ctx, quizAnswer)

		answerResults = append(answerResults, quiz.AnswerResult{
			QuestionID:    question.ID,
			QuestionText:  question.Text,
			YourAnswer:    answer.Answer,
			CorrectAnswer: question.CorrectAnswer,
			IsCorrect:     isCorrect,
			PointsEarned:  pointsEarned,
			MaxPoints:     question.Points,
		})
	}

	// Submit attempt
	if err := s.repo.SubmitAttempt(ctx, attemptID, totalScore, maxScore); err != nil {
		return nil, err
	}

	// Calculate percentage
	percentage := 0.0
	if maxScore > 0 {
		percentage = (totalScore / float64(maxScore)) * 100
	}

	result := &quiz.QuizResult{
		AttemptID:   attemptID,
		QuizID:      attempt.QuizID,
		QuizTitle:   q.Title,
		StudentID:   attempt.StudentID,
		Score:       totalScore,
		MaxScore:    maxScore,
		Percentage:  percentage,
		Passed:      percentage >= 60, // 60% passing grade
		SubmittedAt: time.Now(),
		Answers:     answerResults,
	}

	return result, nil
}

func (s *Service) GetAttemptResult(ctx context.Context, attemptID string) (*quiz.QuizResult, error) {
	if attemptID == "" {
		return nil, appErrors.ErrInvalidID
	}

	attempt, err := s.repo.GetAttemptByID(ctx, attemptID)
	if err != nil {
		return nil, appErrors.ErrAttemptNotFound
	}

	if attempt.SubmittedAt == nil {
		return nil, appErrors.ErrAttemptAlreadySubmitted
	}

	q, _ := s.repo.GetByID(ctx, attempt.QuizID)
	questions, _ := s.repo.GetQuestionsByQuizID(ctx, attempt.QuizID)
	answers, _ := s.repo.GetAnswersByAttempt(ctx, attemptID)

	// Build question map
	questionMap := make(map[string]*quiz.Question)
	for _, question := range questions {
		questionMap[question.ID] = question
	}

	var answerResults []quiz.AnswerResult
	for _, answer := range answers {
		question := questionMap[answer.QuestionID]
		if question != nil {
			answerResults = append(answerResults, quiz.AnswerResult{
				QuestionID:    answer.QuestionID,
				QuestionText:  question.Text,
				YourAnswer:    answer.Answer,
				CorrectAnswer: question.CorrectAnswer,
				IsCorrect:     answer.IsCorrect,
				PointsEarned:  answer.Points,
				MaxPoints:     question.Points,
			})
		}
	}

	result := &quiz.QuizResult{
		AttemptID:   attemptID,
		QuizID:      attempt.QuizID,
		QuizTitle:   q.Title,
		StudentID:   attempt.StudentID,
		Score:       *attempt.Score,
		MaxScore:    attempt.MaxScore,
		Percentage:  *attempt.Percentage,
		Passed:      *attempt.Percentage >= 60,
		SubmittedAt: *attempt.SubmittedAt,
		Answers:     answerResults,
	}

	return result, nil
}

func (s *Service) GetStudentAttempts(ctx context.Context, quizID, studentID string) ([]*quiz.QuizAttempt, error) {
	if quizID == "" || studentID == "" {
		return nil, appErrors.ErrInvalidID
	}

	return s.repo.GetAttemptsByQuizAndStudent(ctx, quizID, studentID)
}

// Helper functions

func validateQuestion(q quiz.QuestionInput) error {
	switch q.Type {
	case quiz.MultipleChoice:
		if len(q.Options) < 2 {
			return appErrors.ErrInvalidInput
		}
		// Check if correct answer is in options
		found := false
		for _, opt := range q.Options {
			if opt == q.CorrectAnswer {
				found = true
				break
			}
		}
		if !found {
			return appErrors.ErrInvalidInput
		}
	case quiz.TrueFalse:
		if q.CorrectAnswer != "true" && q.CorrectAnswer != "false" {
			return appErrors.ErrInvalidInput
		}
	case quiz.ShortAnswer:
		// No special validation
	default:
		return appErrors.ErrInvalidQuestionType
	}

	return nil
}

func gradeAnswer(question *quiz.Question, answer string) bool {
	switch question.Type {
	case quiz.MultipleChoice, quiz.TrueFalse:
		return strings.EqualFold(strings.TrimSpace(answer), strings.TrimSpace(question.CorrectAnswer))
	case quiz.ShortAnswer:
		// Case-insensitive comparison, trim whitespace
		return strings.EqualFold(strings.TrimSpace(answer), strings.TrimSpace(question.CorrectAnswer))
	default:
		return false
	}
}

func optionsToJSON(options []string) json.RawMessage {
	data, _ := json.Marshal(options)
	return data
}
