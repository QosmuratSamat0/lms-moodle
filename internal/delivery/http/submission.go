package http

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/domain/submission"
	submissionUC "github.com/ap1-final-mini-moodle/internal/usecase/submission"
	"github.com/gin-gonic/gin"
)

type SubmissionHandler struct {
	service *submissionUC.Service
}

func NewSubmissionHandler(service *submissionUC.Service) *SubmissionHandler {
	return &SubmissionHandler{service: service}
}

// Submit submits an assignment
// @Summary Submit assignment
// @Description Creates a new submission for an assignment
// @Tags submissions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body struct{AssignmentID string `json:"assignment_id" binding:"required"`; StudentID string `json:"student_id" binding:"required"`; Content string `json:"content" binding:"required"`; FileURL *string `json:"file_url"`} true "Submission Request"
// @Success 201 {object} submission.Submission "Created submission"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/submissions [post]
func (h *SubmissionHandler) Submit(c *gin.Context) {
	var req struct {
		AssignmentID string  `json:"assignment_id" binding:"required"`
		StudentID    string  `json:"student_id" binding:"required"`
		Content      string  `json:"content" binding:"required"`
		FileURL      *string `json:"file_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s, err := h.service.Submit(&submission.CreateSubmissionInput{
		AssignmentID: req.AssignmentID,
		StudentID:    req.StudentID,
		Content:      req.Content,
		FileURL:      req.FileURL,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, s)
}

// GetByID returns a submission by ID
// @Summary Get submission by ID
// @Description Returns submission details by ID
// @Tags submissions
// @Security BearerAuth
// @Produce json
// @Param id path string true "Submission ID"
// @Success 200 {object} submission.Submission "Submission details"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Submission not found"
// @Router /api/v1/submissions/{id} [get]
func (h *SubmissionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	s, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "submission not found"})
		return
	}
	c.JSON(http.StatusOK, s)
}

// ListByAssignment returns submissions for a specific assignment
// @Summary List assignment submissions
// @Description Returns a paginated list of submissions for an assignment
// @Tags submissions
// @Security BearerAuth
// @Produce json
// @Param assignmentID path string true "Assignment ID"
// @Param skip query int false "Skip" default(0)
// @Param take query int false "Take" default(10)
// @Success 200 {array} submission.Submission "Submissions list"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/submissions/assignment/{assignmentID} [get]
func (h *SubmissionHandler) ListByAssignment(c *gin.Context) {
	assignmentID := c.Param("assignmentID")
	var req struct {
		Skip int `form:"skip,default=0"`
		Take int `form:"take,default=10"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	submissions, err := h.service.ListByAssignment(assignmentID, req.Skip, req.Take)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, submissions)
}

// ListByStudent returns submissions for a specific student
// @Summary List student submissions
// @Description Returns a paginated list of submissions from a student
// @Tags submissions
// @Security BearerAuth
// @Produce json
// @Param studentID path string true "Student ID"
// @Param skip query int false "Skip" default(0)
// @Param take query int false "Take" default(10)
// @Success 200 {array} submission.Submission "Submissions list"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/submissions/student/{studentID} [get]
func (h *SubmissionHandler) ListByStudent(c *gin.Context) {
	studentID := c.Param("studentID")
	var req struct {
		Skip int `form:"skip,default=0"`
		Take int `form:"take,default=10"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	submissions, err := h.service.ListByStudent(studentID, req.Skip, req.Take)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, submissions)
}

// Delete deletes a submission
// @Summary Delete submission
// @Description Deletes a submission by ID
// @Tags submissions
// @Security BearerAuth
// @Produce json
// @Param id path string true "Submission ID"
// @Success 204 "No content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Submission not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/submissions/{id} [delete]
func (h *SubmissionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
