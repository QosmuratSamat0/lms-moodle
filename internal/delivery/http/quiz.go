package http

import (
	"net/http"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/quiz"
	quizUC "github.com/ap1-final-mini-moodle/internal/usecase/quiz"
	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	service *quizUC.Service
}

func NewQuizHandler(service *quizUC.Service) *QuizHandler {
	return &QuizHandler{service: service}
}

func (h *QuizHandler) Create(c *gin.Context) {
	userID, _ := c.Get("userID")
	createdBy := ""
	if userID != nil {
		createdBy = userID.(string)
	}

	var req struct {
		CourseID    string        `json:"course_id" binding:"required"`
		Title       string        `json:"title" binding:"required"`
		Description string        `json:"description"`
		TimeLimit   int           `json:"time_limit_minutes"`
		MaxAttempts int           `json:"max_attempts"`
		StartDate   *time.Time    `json:"start_date"`
		EndDate     *time.Time    `json:"end_date"`
		Questions   []questionReq `json:"questions" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	questions := make([]quiz.QuestionInput, len(req.Questions))
	for i, q := range req.Questions {
		questions[i] = quiz.QuestionInput{
			Text:          q.Text,
			Type:          q.Type,
			Options:       q.Options,
			CorrectAnswer: q.CorrectAnswer,
			Points:        q.Points,
		}
	}

	qz, err := h.service.CreateQuiz(c.Request.Context(), createdBy, &quiz.CreateQuizInput{
		CourseID:    req.CourseID,
		Title:       req.Title,
		Description: req.Description,
		TimeLimit:   req.TimeLimit,
		MaxAttempts: req.MaxAttempts,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Questions:   questions,
	})
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, qz)
}

type questionReq struct {
	Text          string            `json:"text" binding:"required"`
	Type          quiz.QuestionType `json:"type" binding:"required"`
	Options       []string          `json:"options"`
	CorrectAnswer string            `json:"correct_answer" binding:"required"`
	Points        int               `json:"points" binding:"required,min=1"`
}

func (h *QuizHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	qz, err := h.service.GetQuizByID(c.Request.Context(), id)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, qz)
}

func (h *QuizHandler) ListByCourse(c *gin.Context) {
	courseID := c.Param("courseID")
	var req struct {
		Limit  int `form:"limit,default=20"`
		Offset int `form:"offset,default=0"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quizzes, total, err := h.service.GetQuizzesByCourse(c.Request.Context(), courseID, req.Limit, req.Offset)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  quizzes,
		"total": total,
	})
}

func (h *QuizHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Published *bool `json:"published"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Published != nil {
		if err := h.service.PublishQuiz(c.Request.Context(), id, *req.Published); err != nil {
			status := getStatusCode(err)
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *QuizHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteQuiz(c.Request.Context(), id); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *QuizHandler) StartAttempt(c *gin.Context) {
	quizID := c.Param("id")
	userID, _ := c.Get("userID")
	studentID := ""
	if userID != nil {
		studentID = userID.(string)
	}

	attempt, err := h.service.StartAttempt(c.Request.Context(), quizID, studentID)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, attempt)
}

func (h *QuizHandler) SubmitAttempt(c *gin.Context) {
	attemptID := c.Param("attemptID")

	var req struct {
		Answers []answerReq `json:"answers" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answers := make([]quiz.AnswerInput, len(req.Answers))
	for i, a := range req.Answers {
		answers[i] = quiz.AnswerInput{
			QuestionID: a.QuestionID,
			Answer:     a.Answer,
		}
	}

	result, err := h.service.SubmitAttempt(c.Request.Context(), attemptID, &quiz.SubmitQuizInput{
		Answers: answers,
	})
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

type answerReq struct {
	QuestionID string `json:"question_id" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
}

func (h *QuizHandler) GetAttemptResult(c *gin.Context) {
	attemptID := c.Param("attemptID")

	result, err := h.service.GetAttemptResult(c.Request.Context(), attemptID)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *QuizHandler) GetStudentAttempts(c *gin.Context) {
	quizID := c.Param("id")
	userID, _ := c.Get("userID")
	studentID := ""
	if userID != nil {
		studentID = userID.(string)
	}

	attempts, err := h.service.GetStudentAttempts(c.Request.Context(), quizID, studentID)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, attempts)
}

func (h *QuizHandler) GetQuizQuestions(c *gin.Context) {
	quizID := c.Param("id")
	questions, err := h.service.GetQuestionsForQuiz(c.Request.Context(), quizID, false)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, questions)
}
