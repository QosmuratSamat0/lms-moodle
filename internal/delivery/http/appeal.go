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

func NewAppealHandler(service *appealUC.Service) *AppealHandler {
	return &AppealHandler{service: service}
}

func (h *AppealHandler) Create(c *gin.Context) {
	userID, _ := c.Get("userID")
	studentID := userID.(string)

	var req struct {
		GradeID  string `json:"grade_id" binding:"required"`
		Reason   string `json:"reason" binding:"required,min=10"`
		Evidence string `json:"evidence"`
	}
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

func (h *AppealHandler) Resolve(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	teacherID := userID.(string)

	var req struct {
		Status   appeal.AppealStatus `json:"status" binding:"required"`
		Response string              `json:"response" binding:"required"`
		NewScore *float64            `json:"new_score"`
	}
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
