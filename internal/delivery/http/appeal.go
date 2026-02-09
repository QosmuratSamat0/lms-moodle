package http

import (
	"net/http"

	"github.com/ap1-final-mini-moodle/internal/domain/appeal"
	appealUC "github.com/ap1-final-mini-moodle/internal/usecase/appeal"
	"github.com/gin-gonic/gin"
)

type AppealHandler struct {
	service *appealUC.Service
}

type CreateAppealRequest struct {
	GradeID  string `json:"grade_id" binding:"required"`
	Reason   string `json:"reason" binding:"required,min=10"`
	Evidence string `json:"evidence"`
}

type ResolveAppealRequest struct {
	Status   appeal.AppealStatus `json:"status" binding:"required"`
	Response string              `json:"response" binding:"required"`
	NewScore *float64            `json:"new_score"`
}

func NewAppealHandler(service *appealUC.Service) *AppealHandler {
	return &AppealHandler{service: service}
}

// Create creates a new grade appeal
// @Summary Create appeal
// @Description Create a new appeal for a grade (Student only)
// @Tags appeals
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateAppealRequest true "Create Appeal Request"
// @Success 201 {object} appeal.GradeAppeal "Created appeal"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/appeals [post]
func (h *AppealHandler) Create(c *gin.Context) {
	userID, _ := c.Get("userID")
	studentID := userID.(string)

	var req CreateAppealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a, err := h.service.CreateAppeal(c.Request.Context(), studentID, &appeal.CreateAppealInput{
		GradeID:  req.GradeID,
		Reason:   req.Reason,
		Evidence: req.Evidence,
	})
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

// GetByID returns an appeal by ID
// @Summary Get appeal by ID
// @Description Returns appeal details by ID
// @Tags appeals
// @Security BearerAuth
// @Produce json
// @Param id path string true "Appeal ID"
// @Success 200 {object} appeal.GradeAppeal "Appeal details"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Appeal not found"
// @Router /api/v1/appeals/{id} [get]
func (h *AppealHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	a, err := h.service.GetAppealByID(c.Request.Context(), id)
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}

// GetStudentAppeals returns a list of appeals for the current student
// @Summary List student appeals
// @Description Returns a paginated list of appeals created by the current student
// @Tags appeals
// @Security BearerAuth
// @Produce json
// @Param status query string false "Filter by status"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{} "Appeals list"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/appeals [get]
func (h *AppealHandler) GetStudentAppeals(c *gin.Context) {
	userID, _ := c.Get("userID")
	studentID := userID.(string)

	var req struct {
		Status string `form:"status"`
		Limit  int    `form:"limit,default=20"`
		Offset int    `form:"offset,default=0"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var status *appeal.AppealStatus
	if req.Status != "" {
		s := appeal.AppealStatus(req.Status)
		status = &s
	}

	appeals, total, err := h.service.GetStudentAppeals(c.Request.Context(), studentID, status, req.Limit, req.Offset)
	if err != nil {
		st := getStatusCode(err)
		c.JSON(st, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  appeals,
		"total": total,
	})
}

// GetTeacherAppeals returns a list of appeals for the current teacher's courses
// @Summary List teacher appeals
// @Description Returns a paginated list of appeals for the current teacher's courses
// @Tags appeals
// @Security BearerAuth
// @Produce json
// @Param status query string false "Filter by status"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{} "Appeals list"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /api/v1/teacher/appeals [get]
func (h *AppealHandler) GetTeacherAppeals(c *gin.Context) {
	userID, _ := c.Get("userID")
	teacherID := userID.(string)

	var req struct {
		Status string `form:"status"`
		Limit  int    `form:"limit,default=20"`
		Offset int    `form:"offset,default=0"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var status *appeal.AppealStatus
	if req.Status != "" {
		s := appeal.AppealStatus(req.Status)
		status = &s
	}

	appeals, total, err := h.service.GetTeacherAppeals(c.Request.Context(), teacherID, status, req.Limit, req.Offset)
	if err != nil {
		st := getStatusCode(err)
		c.JSON(st, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  appeals,
		"total": total,
	})
}

// Resolve resolves a grade appeal
// @Summary Resolve appeal
// @Description Resolve a grade appeal by accepting or rejecting it (Teacher/Admin only)
// @Tags appeals
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Appeal ID"
// @Param request body ResolveAppealRequest true "Resolve Request"
// @Success 200 {object} appeal.GradeAppeal "Resolved appeal"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Appeal not found"
// @Router /api/v1/appeals/{id}/resolve [put]
func (h *AppealHandler) Resolve(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	teacherID := userID.(string)

	var req ResolveAppealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a, err := h.service.ResolveAppeal(c.Request.Context(), id, teacherID, &appeal.ResolveAppealInput{
		Status:   req.Status,
		Response: req.Response,
		NewScore: req.NewScore,
	})
	if err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}

// Delete deletes a grade appeal
// @Summary Delete appeal
// @Description Deletes an appeal by ID (Student/Admin only)
// @Tags appeals
// @Security BearerAuth
// @Produce json
// @Param id path string true "Appeal ID"
// @Success 204 "No content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Appeal not found"
// @Router /api/v1/appeals/{id} [delete]
func (h *AppealHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	studentID := userID.(string)

	if err := h.service.DeleteAppeal(c.Request.Context(), id, studentID); err != nil {
		status := getStatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
