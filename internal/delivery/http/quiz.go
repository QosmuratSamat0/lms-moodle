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

// Create creates a new quiz
// @Summary Create quiz
// @Description Create a new quiz with questions (Teacher/Admin only)
// @Tags quizzes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body struct{CourseID string `json:"course_id" binding:"required"`; Title string `json:"title" binding:"required"`; Description string `json:"description"`; TimeLimit int `json:"time_limit_minutes"`; MaxAttempts int `json:"max_attempts"`; StartDate *string `json:"start_date"`; EndDate *string `json:"end_date"`; Questions []questionReq `json:"questions" binding:"required,min=1"`} true "Create Quiz Request"
// @Success 201 {object} quiz.Quiz "Created quiz"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/quizzes [post]
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

// GetByID returns a quiz by ID
// @Summary Get quiz by ID
// @Description Returns quiz details by ID
// @Tags quizzes
// @Security BearerAuth
// @Produce json
// @Param id path string true "Quiz ID"
// @Success 200 {object} quiz.Quiz "Quiz details"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Quiz not found"
// @Router /api/v1/quizzes/{id} [get]
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

// ListByCourse returns quizzes for a specific course
// @Summary List course quizzes
// @Description Returns a paginated list of quizzes for a course
// @Tags quizzes
// @Security BearerAuth
// @Produce json
// @Param courseID path string true "Course ID"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{} "Quizzes list"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/courses/{courseID}/quizzes [get]
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

// Update updates quiz details (publishing)
// @Summary Update quiz
// @Description Update quiz details (e.g., publish status) (Teacher/Admin only)
// @Tags quizzes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Quiz ID"
// @Param request body struct{Published *bool `json:"published"`} true "Update Request"
// @Success 200 {object} map[string]string "Success message"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Quiz not found"
// @Router /api/v1/quizzes/{id} [put]
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

// Delete deletes a quiz
// @Summary Delete quiz
// @Description Deletes a quiz by ID (Admin only)
// @Tags quizzes
// @Security BearerAuth
// @Produce json
// @Param id path string true "Quiz ID"
// @Success 204 "No content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Quiz not found"
// @Router /api/v1/quizzes/{id} [delete]
func (h *QuizHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteQuiz(c.Request.Context(), id); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// StartAttempt starts a new quiz attempt for a student
// @Summary Start quiz attempt
// @Description Creates a new quiz attempt for the current student
// @Tags quizzes
// @Security BearerAuth
// @Produce json
// @Param id path string true "Quiz ID"
// @Success 201 {object} quiz.QuizAttempt "Created attempt"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Quiz not found"
// @Router /api/v1/quizzes/{id}/attempts [post]
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

// SubmitAttempt submits answers for a quiz attempt
// @Summary Submit quiz attempt
// @Description Submits answers and calculates the result for a quiz attempt
// @Tags quizzes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param attemptID path string true "Attempt ID"
// @Param request body struct{Answers []answerReq `json:"answers" binding:"required"`} true "Submit Request"
// @Success 200 {object} quiz.QuizAttempt "Submission result"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Attempt not found"
// @Router /api/v1/quiz-attempts/{attemptID}/submit [post]
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

// GetAttemptResult returns the result of a quiz attempt
// @Summary Get attempt result
// @Description Returns the result and score of a specific quiz attempt
// @Tags quizzes
// @Security BearerAuth
// @Produce json
// @Param attemptID path string true "Attempt ID"
// @Success 200 {object} quiz.QuizAttempt "Attempt details"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/quiz-attempts/{attemptID}/result [get]
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

// GetStudentAttempts returns all attempts of the current student for a quiz
// @Summary Get student attempts
// @Description Returns a list of all attempts by the current student for a specific quiz
// @Tags quizzes
// @Security BearerAuth
// @Produce json
// @Param id path string true "Quiz ID"
// @Success 200 {array} quiz.QuizAttempt "Attempts list"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/quizzes/{id}/attempts [get]
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

// GetQuizQuestions returns questions for a specific quiz
// @Summary Get quiz questions
// @Description Returns a list of questions for a quiz (Student sees questions without correct answers)
// @Tags quizzes
// @Security BearerAuth
// @Produce json
// @Param id path string true "Quiz ID"
// @Success 200 {array} quiz.Question "Questions list"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/quizzes/{id}/questions [get]
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
