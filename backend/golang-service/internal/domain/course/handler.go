package course

import (
	"net/http"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/errorx"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/middleware"
	"github.com/MaqsattoTeam/aLMS/golang-service/internal/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles course-related HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new course handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create handles course creation (teacher only)
// @Summary Create a new course
// @Tags courses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateCourseRequest true "Course details"
// @Success 201 {object} CourseResponse
// @Failure 400,403 {object} utils.Response
// @Router /courses [post]
func (h *Handler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	course, err := h.service.Create(c.Request.Context(), userID, &req)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.Created(c, course)
}

// GetByID retrieves a course by ID
// @Summary Get course by ID
// @Tags courses
// @Security BearerAuth
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} CourseResponse
// @Failure 404 {object} utils.Response
// @Router /courses/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid course ID")
		return
	}

	course, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, course)
}

// Update updates a course
// @Summary Update course
// @Tags courses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param request body UpdateCourseRequest true "Update details"
// @Success 200 {object} utils.Response
// @Failure 400,403,404 {object} utils.Response
// @Router /courses/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid course ID")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	role := middleware.GetUserRole(c)

	var req UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := h.service.Update(c.Request.Context(), id, userID, role, &req); err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, gin.H{"message": "course updated successfully"})
}

// Delete deletes a course
// @Summary Delete course
// @Tags courses
// @Security BearerAuth
// @Param id path string true "Course ID"
// @Success 204
// @Failure 403,404 {object} utils.Response
// @Router /courses/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid course ID")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	role := middleware.GetUserRole(c)

	if err := h.service.Delete(c.Request.Context(), id, userID, role); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// List lists all courses with optional filters
// @Summary List courses
// @Tags courses
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Param search query string false "Search in title/description"
// @Param is_active query bool false "Filter by active status"
// @Success 200 {object} CourseListResponse
// @Router /courses [get]
func (h *Handler) List(c *gin.Context) {
	pagination := utils.GetPaginationFromContext(c)

	filter := &CourseFilter{}

	if search := c.Query("search"); search != "" {
		filter.Search = &search
	}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	list, err := h.service.List(c.Request.Context(), filter, pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, list)
}

// ListMyCourses lists courses for the current teacher
// @Summary List my courses (teacher)
// @Tags courses
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} CourseListResponse
// @Router /teacher/courses [get]
func (h *Handler) ListMyCourses(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	pagination := utils.GetPaginationFromContext(c)

	list, err := h.service.ListByTeacher(c.Request.Context(), userID, pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, list)
}

// ListEnrolledCourses lists courses for the current student
// @Summary List enrolled courses (student)
// @Tags courses
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} CourseListResponse
// @Router /student/courses [get]
func (h *Handler) ListEnrolledCourses(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	pagination := utils.GetPaginationFromContext(c)

	list, err := h.service.ListByStudent(c.Request.Context(), userID, pagination.Page, pagination.Limit)
	if err != nil {
		errorx.HandleError(c, err)
		return
	}

	utils.OK(c, list)
}

// TransferOwnership transfers course ownership to another teacher
// @Summary Transfer course ownership
// @Tags courses
// @Security BearerAuth
// @Accept json
// @Param id path string true "Course ID"
// @Param request body TransferCourseRequest true "New owner"
// @Success 204
// @Failure 400,403,404 {object} utils.Response
// @Router /courses/{id}/transfer [post]
func (h *Handler) TransferOwnership(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequest(c, "invalid course ID")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		utils.Unauthorized(c, "user not authenticated")
		return
	}

	role := middleware.GetUserRole(c)

	var req TransferCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, utils.FormatValidationErrors(err))
		return
	}

	if err := h.service.TransferOwnership(c.Request.Context(), id, userID, req.NewOwnerTeacherID, role); err != nil {
		errorx.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
