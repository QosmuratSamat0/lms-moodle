package deliveryhttp

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"ap1-final-mini-moodle/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	teacher *usecase.TeacherUsecase
	admin   *usecase.AdminUsecase
}

func NewHandler(teacher *usecase.TeacherUsecase, admin *usecase.AdminUsecase) *Handler {
	return &Handler{
		teacher: teacher,
		admin:   admin,
	}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) CreateCourse(c *gin.Context) {
	user, ok := UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var payload struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	course, err := h.teacher.CreateCourse(user.ID, user.Name, payload.Title, payload.Description)
	if err != nil {
		writeUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusCreated, course)
}

func (h *Handler) CreateAssignment(c *gin.Context) {
	user, ok := UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	courseID, err := strconv.ParseInt(c.Param("courseId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}

	var payload struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DueAt       string `json:"due_at"`
		MaxPoints   int    `json:"max_points"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	dueAt, err := time.Parse(time.RFC3339, strings.TrimSpace(payload.DueAt))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due_at"})
		return
	}

	assignment, err := h.teacher.CreateAssignment(user.ID, courseID, payload.Title, payload.Description, dueAt, payload.MaxPoints)
	if err != nil {
		writeUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusCreated, assignment)
}

func (h *Handler) ListSubmissions(c *gin.Context) {
	user, ok := UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	assignmentID, err := strconv.ParseInt(c.Param("assignmentId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignment id"})
		return
	}

	submissions, err := h.teacher.ListSubmissions(user.ID, assignmentID)
	if err != nil {
		writeUsecaseError(c, err)
		return
	}

	response := make([]gin.H, 0, len(submissions))
	for _, submission := range submissions {
		item := gin.H{
			"id":            submission.Submission.ID,
			"assignment_id": submission.Submission.AssignmentID,
			"student_id":    submission.Submission.StudentID,
			"submitted_at":  submission.Submission.SubmittedAt,
			"content_ref":   submission.Submission.ContentRef,
			"grade":         nil,
		}
		if submission.Grade != nil {
			item["grade"] = gin.H{
				"score":    submission.Grade.Score,
				"feedback": submission.Grade.Feedback,
			}
		}
		response = append(response, item)
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GradeSubmission(c *gin.Context) {
	user, ok := UserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	submissionID, err := strconv.ParseInt(c.Param("submissionId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid submission id"})
		return
	}

	var payload struct {
		Score    int    `json:"score"`
		Feedback string `json:"feedback"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	grade, err := h.teacher.GradeSubmission(user.ID, submissionID, payload.Score, payload.Feedback)
	if err != nil {
		writeUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, grade)
}

func (h *Handler) AdminCourses(c *gin.Context) {
	courses := h.admin.ListCourses()
	response := make([]gin.H, 0, len(courses))
	for _, course := range courses {
		response = append(response, gin.H{
			"id":    course.ID,
			"title": course.Title,
			"teacher": gin.H{
				"id":   course.TeacherID,
				"name": course.TeacherName,
			},
			"student_count": course.StudentCount,
		})
	}

	c.JSON(http.StatusOK, response)
}

func writeUsecaseError(c *gin.Context, err error) {
	switch err {
	case usecase.ErrInvalid:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case usecase.ErrForbidden:
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case usecase.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}
